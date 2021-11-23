package main

import (
	"bytes"
	"fmt"
	"syscall"

	ext3 "github.com/enoly/extFileRecovery/pkg/ext3"
	ext3Journal "github.com/enoly/extFileRecovery/pkg/ext3/journal"
	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"
)

func main() {
	disk := "/dev/sdb1"
	var err error

	fd, err := syscall.Open(disk, syscall.O_RDONLY, 0777)

	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	syscall.Seek(fd, 1024, 0)
	rawSuperblock := make([]byte, 1024)

	_, err = syscall.Read(fd, rawSuperblock)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	sb := ext3.NewSuperblock()
	err = sb.Read(kaitai.NewStream(bytes.NewReader(rawSuperblock)), sb, sb)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	var blockSize uint32 = 1024 << sb.LogBlockSize
	var GDTSeek uint64 = uint64((sb.FirstDataBlock + 1) * blockSize)

	var blockGroupNum uint32 = sb.BlocksCount / sb.BlocksPerGroup
	if sb.BlocksCount%sb.BlocksPerGroup != 0 {
		blockGroupNum++
	}

	var GDTSize uint32 = blockGroupNum * 32
	fmt.Printf("First data block is %v\n", sb.FirstDataBlock)
	fmt.Printf("GDT starts on %v byte from disk start\n", int64(GDTSeek))
	syscall.Seek(fd, int64(GDTSeek), 0)

	rawGDT := make([]byte, GDTSize)
	_, err = syscall.Read(fd, rawGDT)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	GDT := make([]ext3.GdtRecord, 0)
	for i := 0; i < int(blockGroupNum); i++ {
		GDTRecord := ext3.NewGdtRecord()
		GDTRecord.Read(kaitai.NewStream(bytes.NewReader(rawGDT[i*32:(i+1)*32])), GDTRecord, GDTRecord)
		GDT = append(GDT, *GDTRecord)
	}

	fmt.Printf("First block group descriptor %+v\n", GDT[0])
	fmt.Printf("Inode size %v\n", sb.InodeSize)

	var inodeTableSeek uint64 = uint64(GDT[0].InodeTableBlock) * uint64(blockSize)
	syscall.Seek(fd, int64(inodeTableSeek), 0)

	rawInodeTable := make([]byte, sb.InodesPerGroup*uint32(sb.InodeSize))
	_, err = syscall.Read(fd, rawInodeTable)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	inodeTable := make([]ext3.Inode, 0)
	for i := 0; i < int(sb.InodesPerGroup); i++ {
		inode := ext3.NewInode()
		inode.Read(kaitai.NewStream(bytes.NewReader(rawInodeTable[i*int(sb.InodeSize):(i+1)*int(sb.InodeSize)])), inode, inode)
		inodeTable = append(inodeTable, *inode)
	}

	rawRootDirectory := make([]byte, 0)
	fmt.Printf("Inode 2, size: %v, direct blocks:\n", inodeTable[1].Size)
	for _, blockNum := range inodeTable[1].DirectBlocks {
		fmt.Printf("%v ", blockNum.Ptr)

		if blockNum.Ptr != 0 {
			syscall.Seek(fd, int64(blockNum.Ptr)*int64(blockSize), 0)

			rawBlock := make([]byte, blockSize)
			_, err = syscall.Read(fd, rawBlock)
			if err != nil {
				fmt.Printf("Error: %v", err)
				return
			}

			rawRootDirectory = append(rawRootDirectory, rawBlock...)
		}
	}
	fmt.Println()

	rootDirectory := ext3.NewDirectory()
	rootDirectory.Read(kaitai.NewStream(bytes.NewReader(rawRootDirectory)), rootDirectory, rootDirectory)

	for _, entry := range rootDirectory.Entries {
		if entry.FileType != ext3.Directory_DirEntry_FileTypeEnum__Unknown && entry.InodePtr != 0 {
			fmt.Printf("%v\n", entry.Name)

			if len(entry.Padding) > 3 {
				deletedFiles := ext3.NewDeletedFiles()
				deletedFiles.Read(kaitai.NewStream(bytes.NewReader(entry.Padding)), deletedFiles, deletedFiles)

				for _, deletedFile := range deletedFiles.Entries {
					if deletedFile.FileType != ext3.DeletedFiles_DeletedEntry_FileTypeEnum__Unknown && deletedFile.InodePtr != 0 {
						fmt.Printf("deleted: %v\n", deletedFile.Name)
					}
				}
			}
		}
	}

	journalInode := inodeTable[sb.JournalInodeNum-1]
	journalBlocks := []uint32{}
	for _, directJournalBlock := range journalInode.DirectBlocks {
		if directJournalBlock.Ptr != 0 {
			journalBlocks = append(journalBlocks, directJournalBlock.Ptr)
		}
	}
	syscall.Seek(fd, int64(journalInode.FirstLevelIndirectBlock.Ptr*blockSize), 0)
	rawFirstIndirectJournalBlock := make([]byte, blockSize)
	_, err = syscall.Read(fd, rawFirstIndirectJournalBlock)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	firstIndirectJournalBlock := ext3.NewIndirectBlock()
	firstIndirectJournalBlock.Read(kaitai.NewStream(bytes.NewReader(rawFirstIndirectJournalBlock)), firstIndirectJournalBlock, firstIndirectJournalBlock)
	for _, directJournalBlock := range firstIndirectJournalBlock.BlocksPtrs {
		if directJournalBlock.Ptr != 0 {
			journalBlocks = append(journalBlocks, directJournalBlock.Ptr)
		}
	}

	for journalBlockNum, journalBlock := range journalBlocks {
		syscall.Seek(fd, int64(journalBlock*blockSize), 0)
		rawJournalBlock := make([]byte, blockSize)
		_, err = syscall.Read(fd, rawJournalBlock)
		if err != nil {
			fmt.Printf("Error: %v", err)
			return
		}

		journalBlockHeader := ext3Journal.NewHeader()
		journalBlockHeader.Read(kaitai.NewStream(bytes.NewReader(rawJournalBlock)), journalBlockHeader, journalBlockHeader)
		switch journalBlockHeader.BlockType {
		case ext3Journal.Header_BlockTypeEnum__SuperblockV2:
			fmt.Printf("[%v] Superblock\n", journalBlockNum)
		case ext3Journal.Header_BlockTypeEnum__DescriptorBlock:
			journalDescriptor := ext3Journal.NewDescriptor()
			journalDescriptor.Read(kaitai.NewStream(bytes.NewReader(rawJournalBlock)), journalDescriptor, journalDescriptor)
			fmt.Printf("[%v] Descriptor block\n", journalBlockNum)
			for i, descrBlock := range journalDescriptor.Descriptors {
				fmt.Printf("[%v] FS block %v\n", journalBlockNum+i+1, descrBlock.FsBlockNum)
			}
		case ext3Journal.Header_BlockTypeEnum__CommitBlock:
			fmt.Printf("[%v] Commit block\n", journalBlockNum)
		case ext3Journal.Header_BlockTypeEnum__RevokeBlock:
			fmt.Printf("[%v] Revoke block\n", journalBlockNum)
		}
	}

	err = syscall.Close(fd)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
}
