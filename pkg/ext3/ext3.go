package ext3

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/enoly/extFileRecovery/pkg/ext3/model"
	structure "github.com/enoly/extFileRecovery/pkg/ext3/structure"
	journal "github.com/enoly/extFileRecovery/pkg/ext3/structure/journal"
	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"
)

const (
	FS_PERMISSION   = 0777
	SEEK_WHENCE     = 0
	SUPERBLOCK_SEEK = 1024
	SUPERBLOCK_SIZE = 1024
	GDT_RECORD_SIZE = 32
	ROOT_DIR_INODE  = 2
)

type Ext3 struct {
	Drive      string
	BlockSize  uint32
	Superblock structure.Superblock
	GDT        []structure.GdtRecord
	Journal    []model.JournalRecord
	drive      *os.File
}

func New(drive string) (*Ext3, error) {
	driveFile, err := os.Open(drive)
	if err != nil {
		return nil, err
	}

	ext3 := &Ext3{
		Drive: drive,
		drive: driveFile,
	}

	if err := ext3.readSuperblock(); err != nil {
		return nil, fmt.Errorf("unable to process ext3 superblock: %v", err)
	}

	ext3.BlockSize = 1024 << ext3.Superblock.LogBlockSize

	if err := ext3.readGDT(); err != nil {
		return nil, fmt.Errorf("unable to process ext3 GDT: %v", err)
	}

	if err := ext3.readJournal(); err != nil {
		return nil, fmt.Errorf("unable to process ext3 journal: %v", err)
	}

	return ext3, nil
}

func (e *Ext3) Close() {
	e.drive.Close()
}

func readBytes(driveFile *os.File, buffer *[]byte, seek int64) error {
	_, err := driveFile.ReadAt(*buffer, seek)
	if err != nil {
		return err
	}

	return nil
}

func (e *Ext3) readSuperblock() error {
	rawSb := make([]byte, SUPERBLOCK_SIZE)
	if err := readBytes(e.drive, &rawSb, SUPERBLOCK_SEEK); err != nil {
		return fmt.Errorf("unable to read superblock from drive: %v", err)
	}

	sb := structure.NewSuperblock()
	if err := sb.Read(kaitai.NewStream(bytes.NewReader(rawSb)), sb, sb); err != nil {
		return fmt.Errorf("unable to read superblock from bytes: %v", err)
	}

	e.Superblock = *sb
	return nil
}

func (e *Ext3) readGDT() error {
	blockGroupsCount := e.Superblock.BlocksCount / e.Superblock.BlocksPerGroup
	if e.Superblock.BlocksCount%e.Superblock.BlocksPerGroup != 0 {
		blockGroupsCount++
	}

	gdtSeek := int64((e.Superblock.FirstDataBlock + 1) * e.BlockSize)
	gdtSize := GDT_RECORD_SIZE * blockGroupsCount
	rawGDT := make([]byte, gdtSize)
	if err := readBytes(e.drive, &rawGDT, gdtSeek); err != nil {
		return fmt.Errorf("unable to read GDT: %v", err)
	}

	for i := uint32(0); i < blockGroupsCount; i++ {
		record := structure.NewGdtRecord()
		if err := record.Read(kaitai.NewStream(bytes.NewReader(rawGDT[i*GDT_RECORD_SIZE:(i+1)*GDT_RECORD_SIZE])), record, record); err != nil {
			return fmt.Errorf("unable to read GDT record: %v", err)
		}
		e.GDT = append(e.GDT, *record)
	}

	return nil
}

func (e *Ext3) ReadBlock(blockNum uint32) (*[]byte, error) {
	if blockNum >= e.Superblock.BlocksCount {
		return nil, fmt.Errorf("block %v not found", blockNum)
	}

	block := make([]byte, e.BlockSize)
	if blockNum == 0 {
		return &block, nil
	}

	blockSeek := int64(blockNum * e.BlockSize)

	if err := readBytes(e.drive, &block, blockSeek); err != nil {
		return nil, fmt.Errorf("unable to read FS block: %v", err)
	}

	return &block, nil
}

func (e *Ext3) GetInodeBlockGroup(inodeNum uint32) uint32 {
	return uint32(inodeNum / e.Superblock.InodesPerGroup)
}

func (e *Ext3) GetInodeBlockPtr(inodeNum uint32) (uint32, error) {
	if inodeNum >= e.Superblock.InodesCount {
		return 0, fmt.Errorf("inode %v not found", inodeNum)
	}

	inodeBlockGroup := e.GetInodeBlockGroup(inodeNum)
	groupDescriptor := e.GDT[inodeBlockGroup]
	inodeTableStartBlock := groupDescriptor.InodeTableBlock
	inodeNumInTable := inodeNum % e.Superblock.InodesPerGroup
	inodesPerBlock := e.BlockSize / uint32(e.Superblock.InodeSize)
	return uint32(inodeNumInTable/inodesPerBlock) + inodeTableStartBlock, nil
}

func (e *Ext3) GetInodeSeekInBlock(inodeNum uint32) uint32 {
	inodeNumInTable := inodeNum % e.Superblock.InodesPerGroup
	inodesPerBlock := e.BlockSize / uint32(e.Superblock.InodeSize)
	return (inodeNumInTable - 1) % inodesPerBlock
}

func (e *Ext3) ReadInode(inodeNum uint32) (*structure.Inode, error) {
	inodeTableBlock, err := e.GetInodeBlockPtr(inodeNum)
	if err != nil {
		return nil, err
	}

	rawInodeTablePart, err := e.ReadBlock(inodeTableBlock)
	if err != nil {
		return nil, fmt.Errorf("unable to read inode table: %v", err)
	}

	inodeSeek := e.GetInodeSeekInBlock(inodeNum)
	inodeSize := uint32(e.Superblock.InodeSize)
	rawInode := (*rawInodeTablePart)[inodeSeek*inodeSize : (inodeSeek+1)*inodeSize]
	inode := structure.NewInode()
	if err := inode.Read(kaitai.NewStream(bytes.NewReader(rawInode)), inode, inode); err != nil {
		return nil, fmt.Errorf("unable to process inode: %v", err)
	}
	return inode, nil
}

func (e *Ext3) ReadIndirectAddrBlock(blockNum uint32) (*structure.IndirectBlock, error) {
	rawBlock, err := e.ReadBlock(blockNum)
	if err != nil {
		return nil, fmt.Errorf("unable to read indirect block: %v", err)
	}

	indirectBlock := structure.NewIndirectBlock()
	if err := indirectBlock.Read(kaitai.NewStream(bytes.NewReader(*rawBlock)), indirectBlock, indirectBlock); err != nil {
		return nil, fmt.Errorf("unable to process indirect block: %v", err)
	}

	return indirectBlock, nil
}

func (e *Ext3) readPtrsFromIndirectBlock(indirectBlockPtr uint32) (*[]uint32, error) {
	ptrs := []uint32{}

	indirectBlock, err := e.ReadIndirectAddrBlock(indirectBlockPtr)
	if err != nil {
		return nil, err
	}
	for _, ptr := range indirectBlock.BlocksPtrs {
		ptrs = append(ptrs, ptr.Ptr)
	}

	return &ptrs, nil
}

func (e *Ext3) ReadFilePtrsFromInode(inode *structure.Inode) (*[]uint32, error) {
	fileBlocks := []uint32{}

	for _, ptr := range inode.DirectBlocks {
		fileBlocks = append(fileBlocks, ptr.Ptr)
	}
	if inode.FirstLevelIndirectBlock == nil || inode.FirstLevelIndirectBlock.Ptr == 0 {
		return &fileBlocks, nil
	}

	firstLevelIndirectBlockPtrs, err := e.readPtrsFromIndirectBlock(inode.FirstLevelIndirectBlock.Ptr)
	if err != nil {
		return nil, fmt.Errorf("unable to read blocks from first level indirect blocks: %v", err)
	}
	fileBlocks = append(fileBlocks, *firstLevelIndirectBlockPtrs...)
	if inode.SecondLevelIndirectBlock == nil || inode.SecondLevelIndirectBlock.Ptr == 0 {
		return &fileBlocks, nil
	}

	secondLevelIndirectBlockPtrs, err := e.readPtrsFromIndirectBlock(inode.SecondLevelIndirectBlock.Ptr)
	if err != nil {
		return nil, fmt.Errorf("unable to read second level indirect blocks: %v", err)
	}
	for _, indirectBlockPtr := range *secondLevelIndirectBlockPtrs {
		indirectBlockPtrs, err := e.readPtrsFromIndirectBlock(indirectBlockPtr)
		if err != nil {
			return nil, fmt.Errorf("unable to read indirect blocks from second level indirect block: %v", err)
		}

		fileBlocks = append(fileBlocks, *indirectBlockPtrs...)
	}
	if inode.ThirdLevelIndirectBlock == nil || inode.ThirdLevelIndirectBlock.Ptr == 0 {
		return &fileBlocks, nil
	}

	thirdLevelIndirectBlockPtrs, err := e.readPtrsFromIndirectBlock(inode.ThirdLevelIndirectBlock.Ptr)
	if err != nil {
		return nil, fmt.Errorf("unable to read third level indirect blocks: %v", err)
	}
	for _, secondLevelIndirectBlockPtr := range *thirdLevelIndirectBlockPtrs {
		secondLevelIndirectBlockPtrs, err := e.readPtrsFromIndirectBlock(secondLevelIndirectBlockPtr)
		if err != nil {
			return nil, fmt.Errorf("unable to read second level indirect blocks from third level indirect block: %v", err)
		}

		for _, indirectBlockPtr := range *secondLevelIndirectBlockPtrs {
			indirectBlockPtrs, err := e.readPtrsFromIndirectBlock(indirectBlockPtr)
			if err != nil {
				return nil, fmt.Errorf("unable to read indirect blocks from second level indirect block: %v", err)
			}

			fileBlocks = append(fileBlocks, *indirectBlockPtrs...)
		}
	}

	return &fileBlocks, nil
}

func (e *Ext3) ReadFileFromPtrs(filePtrs *[]uint32) (*[]byte, error) {
	filePtrsCopy := make([]uint32, len(*filePtrs))
	copy(filePtrsCopy, *filePtrs)

	for i := len(filePtrsCopy) - 1; i >= 0; i-- {
		if (*filePtrs)[i] == 0 {
			filePtrsCopy = filePtrsCopy[:i]
		} else {
			break
		}
	}

	file := []byte{}
	for _, ptr := range filePtrsCopy {
		block, err := e.ReadBlock(ptr)
		if err != nil {
			return nil, fmt.Errorf("unable to read file: %v", err)
		}

		file = append(file, *block...)
	}

	return &file, nil
}

func (e *Ext3) ReadFileFromInode(inode *structure.Inode) (*[]byte, error) {
	blocksPtrs, err := e.ReadFilePtrsFromInode(inode)
	if err != nil {
		return nil, fmt.Errorf("unable to read file blocks pointers: %v", err)
	}
	blocks, err := e.ReadFileFromPtrs(blocksPtrs)
	if err != nil {
		return nil, fmt.Errorf("unabel to read file blocks: %v", err)
	}

	return blocks, nil
}

func (e *Ext3) ReadFileFromInodePtr(inodeNum uint32) (*[]byte, error) {
	inode, err := e.ReadInode(inodeNum)
	if err != nil {
		return nil, fmt.Errorf("unable to read file inode: %v", err)
	}
	return e.ReadFileFromInode(inode)
}

func (e *Ext3) ReadDirFromInodePtr(inodeNum uint32) (*structure.Directory, error) {
	dirBlocks, err := e.ReadFileFromInodePtr(inodeNum)
	if err != nil {
		return nil, fmt.Errorf("unable to read directoryL %v", err)
	}

	dir := structure.NewDirectory()
	if err := dir.Read(kaitai.NewStream(bytes.NewReader(*dirBlocks)), dir, dir); err != nil {
		return nil, fmt.Errorf("unanble to process directory: %v", err)
	}

	return dir, nil
}

func (e *Ext3) GetDeletedFilesFromDir(dir *structure.Directory) (*structure.DeletedFiles, error) {
	deleted := structure.NewDeletedFiles()

	for _, entry := range dir.Entries {
		padding := make([]byte, len(entry.Padding))
		copy(padding, entry.Padding)

		if (8+entry.NameLen)%4 != 0 {
			padding = padding[4-((8+entry.NameLen)%4):]
		}

		localDeleted := structure.NewDeletedFiles()
		if err := localDeleted.Read(kaitai.NewStream(bytes.NewReader(padding)), localDeleted, localDeleted); err != io.EOF && err != nil {
			continue
		}

		for _, deletedEntry := range localDeleted.Entries {
			if deletedEntry.NameLen > 0 {
				deleted.Entries = append(deleted.Entries, deletedEntry)
			}
		}
	}

	return deleted, nil
}

func (e *Ext3) readJournal() error {
	journalInodeNum := e.Superblock.JournalInodeNum
	journalInode, err := e.ReadInode(journalInodeNum)
	if err != nil {
		return fmt.Errorf("unable to read journal inode: %v", err)
	}
	journalPtrs, err := e.ReadFilePtrsFromInode(journalInode)
	if err != nil {
		return fmt.Errorf("unable to read journal ptrs: %v", err)
	}

	superblockRaw, err := e.ReadBlock((*journalPtrs)[0])
	if err != nil {
		return fmt.Errorf("unable to read journal first block: %v", err)
	}
	superblockHeader := journal.NewHeader()
	if err := superblockHeader.Read(kaitai.NewStream(bytes.NewReader((*superblockRaw)[:12])), superblockHeader, superblockHeader); err != nil {
		return fmt.Errorf("unable to process journal first header: %v", err)
	}
	if superblockHeader.BlockType != journal.Header_BlockTypeEnum__SuperblockV2 {
		return fmt.Errorf("unable to process journal superblock: not found")
	}
	superblock := journal.NewSuperblock()
	if err := superblock.Read(kaitai.NewStream(bytes.NewReader((*superblockRaw)[:1024])), superblock, superblock); err != nil {
		return fmt.Errorf("unable to process journal superblock: %v", err)
	}

	journalRecords := []model.JournalRecord{{
		Type:        model.JournalSuperblock,
		RecordNum:   0,
		RecordBlock: (*journalPtrs)[0],
		Description: "superblock",
	}}
	journalBlocksCount := int(journalInode.Size / superblock.BlockSize)

	journalBlock := 1
	journalBlocksPerFsBlock := e.BlockSize / superblock.BlockSize
	for journalBlock < journalBlocksCount {
		blockRaw, err := e.ReadBlock((*journalPtrs)[journalBlock])
		if err != nil {
			return fmt.Errorf("unable to read journal block: %v", err)
		}

		journalBlockStart := (journalBlock % int(journalBlocksPerFsBlock)) * int(superblock.BlockSize)
		header := journal.NewHeader()
		headerRaw := (*blockRaw)[journalBlockStart : journalBlockStart+12]
		if err := header.Read(kaitai.NewStream(bytes.NewReader(headerRaw)), header, header); err != nil {
			journalRecords = append(journalRecords, model.JournalRecord{
				Type:        model.JournalInvalid,
				RecordNum:   uint32(journalBlock),
				RecordBlock: (*journalPtrs)[uint32(journalBlock)*superblock.BlockSize/e.BlockSize],
				Description: "Invalid record",
			})
		}

		switch header.BlockType {
		case journal.Header_BlockTypeEnum__DescriptorBlock:
			{
				journalRecords = append(journalRecords, model.JournalRecord{
					Type:        model.JournalDescriptor,
					RecordNum:   uint32(journalBlock),
					RecordBlock: (*journalPtrs)[uint32(journalBlock)*superblock.BlockSize/e.BlockSize],
					Description: "Descriptor block",
				})

				descriptorRaw := (*blockRaw)[journalBlockStart : journalBlockStart+int(superblock.BlockSize)]
				descriptor := journal.NewDescriptor()
				if err := descriptor.Read(kaitai.NewStream(bytes.NewReader(descriptorRaw)), descriptor, descriptor); err != nil {
					return fmt.Errorf("unable to process journal descriptor: %v", err)
				}

				for i, record := range descriptor.Descriptors {
					journalRecords = append(journalRecords, model.JournalRecord{
						Type:        model.JournalMetadata,
						RecordNum:   uint32(journalBlock + i + 1),
						RecordBlock: (*journalPtrs)[uint32(journalBlock+i+1)*superblock.BlockSize/e.BlockSize],
						Description: fmt.Sprintf("%v", record.FsBlockNum),
					})
				}

				journalBlock += len(descriptor.Descriptors) + 1
				continue
			}
		case journal.Header_BlockTypeEnum__CommitBlock:
			{
				journalRecords = append(journalRecords, model.JournalRecord{
					Type:        model.JournalCommit,
					RecordNum:   uint32(journalBlock),
					RecordBlock: (*journalPtrs)[uint32(journalBlock)*superblock.BlockSize/e.BlockSize],
					Description: "Commit block",
				})
			}
		case journal.Header_BlockTypeEnum__RevokeBlock:
			{
				journalRecords = append(journalRecords, model.JournalRecord{
					Type:        model.JournalRevoke,
					RecordNum:   uint32(journalBlock),
					RecordBlock: (*journalPtrs)[uint32(journalBlock)*superblock.BlockSize/e.BlockSize],
					Description: "Revoke block",
				})
			}
		default:
			{
				journalRecords = append(journalRecords, model.JournalRecord{
					Type:        model.JournalInvalid,
					RecordNum:   uint32(journalBlock),
					RecordBlock: (*journalPtrs)[uint32(journalBlock)*superblock.BlockSize/e.BlockSize],
					Description: "Unknown block",
				})
			}
		}

		journalBlock++
	}

	e.Journal = journalRecords
	return nil
}

func (e *Ext3) GetJournal() *[]model.JournalRecord {
	return &e.Journal
}

func (e *Ext3) GetFreshInodeFromJournal(inodeNum uint32) (*structure.Inode, error) {
	var freshInode structure.Inode
	inodeBlock, err := e.GetInodeBlockPtr(inodeNum)
	if err != nil {
		return nil, err
	}
	inodeSeek := e.GetInodeSeekInBlock(inodeNum)
	inodeSize := uint32(e.Superblock.InodeSize)

	for _, record := range e.Journal {
		if record.Type == model.JournalMetadata {
			recordBlockPtr, err := strconv.ParseUint(record.Description, 10, 32)
			if err != nil || uint32(recordBlockPtr) != inodeBlock {
				continue
			}

			journalBlock, err := e.ReadBlock(record.RecordBlock)
			if err != nil {
				continue
			}

			tmpInode := structure.NewInode()
			rawInode := (*journalBlock)[inodeSeek*inodeSize : (inodeSeek+1)*inodeSize]
			if err := tmpInode.Read(kaitai.NewStream(bytes.NewReader(rawInode)), tmpInode, tmpInode); err != nil {
				continue
			}

			ptrs, err := e.ReadFilePtrsFromInode(tmpInode)
			if err != nil {
				continue
			}

			flag := false
			for _, ptr := range *ptrs {
				if ptr != 0 {
					flag = true
				}

				if ptr >= e.Superblock.BlocksCount {
					flag = false
					break
				}
			}

			if !flag {
				continue
			}

			freshInode = *tmpInode
		}
	}

	return &freshInode, nil
}

func (e *Ext3) GetBlockGroupBitmap(blockGroupNum uint32) (*[]bool, error) {
	bitmapBlockPtr := e.GDT[blockGroupNum].BlockBitmapBlock
	bitmapBlock, err := e.ReadBlock(bitmapBlockPtr)
	if err != nil {
		return nil, err
	}

	bitmap := structure.NewBitmap()
	if err := bitmap.Read(kaitai.NewStream(bytes.NewReader(*bitmapBlock)), bitmap, bitmap); err != nil {
		return nil, err
	}

	return &bitmap.Flags, nil
}
