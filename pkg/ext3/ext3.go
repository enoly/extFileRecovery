package ext3

import (
	"bytes"
	"fmt"
	"strconv"
	"syscall"

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
	drive      string
	blockSize  uint32
	superblock structure.Superblock
	gdt        []structure.GdtRecord
	journal    []model.JournalRecord
}

func New(drive string) (*Ext3, error) {
	ext3 := &Ext3{
		drive: drive,
	}

	if err := ext3.readSuperblock(); err != nil {
		return nil, fmt.Errorf("unable to process ext3 superblock: %v", err)
	}

	ext3.blockSize = 1024 << ext3.superblock.LogBlockSize

	if err := ext3.readGDT(); err != nil {
		return nil, fmt.Errorf("unable to process ext3 GDT: %v", err)
	}

	if err := ext3.readJournal(); err != nil {
		return nil, fmt.Errorf("unable to process ext3 journal: %v", err)
	}

	return ext3, nil
}

func readBytes(drive string, buffer *[]byte, seek int64) error {
	fd, err := syscall.Open(drive, syscall.O_RDONLY, FS_PERMISSION)
	if err != nil {
		return err
	}

	defer syscall.Close(fd)

	if _, err := syscall.Seek(fd, seek, SEEK_WHENCE); err != nil {
		return err
	}

	if _, err = syscall.Read(fd, *buffer); err != nil {
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

	e.superblock = *sb
	return nil
}

func (e *Ext3) readGDT() error {
	blockGroupsCount := e.superblock.BlocksCount / e.superblock.BlocksPerGroup
	if e.superblock.BlocksCount%e.superblock.BlocksPerGroup != 0 {
		blockGroupsCount++
	}

	gdtSeek := int64((e.superblock.FirstDataBlock + 1) * e.blockSize)
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
		e.gdt = append(e.gdt, *record)
	}

	return nil
}

func (e *Ext3) ReadBlock(blockNum uint32) (*[]byte, error) {
	if blockNum >= e.superblock.BlocksCount {
		return nil, fmt.Errorf("block %v not found", blockNum)
	}

	block := make([]byte, e.blockSize)
	if blockNum == 0 {
		return &block, nil
	}

	blockSeek := int64(blockNum * e.blockSize)

	if err := readBytes(e.drive, &block, blockSeek); err != nil {
		return nil, fmt.Errorf("unable to read FS block: %v", err)
	}

	return &block, nil
}

func (e *Ext3) GetInodeBlockPtr(inodeNum uint32) (uint32, error) {
	if inodeNum >= e.superblock.InodesCount {
		return 0, fmt.Errorf("inode %v not found", inodeNum)
	}

	inodeBlockGroup := uint32(inodeNum / e.superblock.InodesPerGroup)
	groupDescriptor := e.gdt[inodeBlockGroup]
	inodeTableStartBlock := groupDescriptor.InodeTableBlock
	inodeNumInTable := inodeNum % e.superblock.InodesPerGroup
	inodesPerBlock := e.blockSize / uint32(e.superblock.InodeSize)
	return uint32(inodeNumInTable/inodesPerBlock) + inodeTableStartBlock, nil
}

func (e *Ext3) GetInodeSeekInBlock(inodeNum uint32) uint32 {
	inodeNumInTable := inodeNum % e.superblock.InodesPerGroup
	inodesPerBlock := e.blockSize / uint32(e.superblock.InodeSize)
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
	inodeSize := uint32(e.superblock.InodeSize)
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

func (e *Ext3) readFilePtrsFromInode(inode *structure.Inode) (*[]uint32, error) {
	fileBlocks := []uint32{}

	for _, ptr := range inode.DirectBlocks {
		fileBlocks = append(fileBlocks, ptr.Ptr)
	}
	if inode.FirstLevelIndirectBlock.Ptr == 0 {
		return &fileBlocks, nil
	}

	firstLevelIndirectBlockPtrs, err := e.readPtrsFromIndirectBlock(inode.FirstLevelIndirectBlock.Ptr)
	if err != nil {
		return nil, fmt.Errorf("unable to read blocks from first level indirect blocks: %v", err)
	}
	fileBlocks = append(fileBlocks, *firstLevelIndirectBlockPtrs...)
	if inode.SecondLevelIndirectBlock.Ptr == 0 {
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
	if inode.ThirdLevelIndirectBlock.Ptr == 0 {
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

func (e *Ext3) readFileFromPtrs(filePtrs *[]uint32) (*[]byte, error) {
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
	blocksPtrs, err := e.readFilePtrsFromInode(inode)
	if err != nil {
		return nil, fmt.Errorf("unable to read file blocks pointers: %v", err)
	}
	blocks, err := e.readFileFromPtrs(blocksPtrs)
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
		localDeleted := structure.NewDeletedFiles()
		if err := localDeleted.Read(kaitai.NewStream(bytes.NewReader(entry.Padding)), localDeleted, localDeleted); err != nil {
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
	journalInodeNum := e.superblock.JournalInodeNum
	journalInode, err := e.ReadInode(journalInodeNum)
	if err != nil {
		return fmt.Errorf("unable to read journal inode: %v", err)
	}
	journalPtrs, err := e.readFilePtrsFromInode(journalInode)
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
	journalBlocksPerFsBlock := e.blockSize / superblock.BlockSize
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
				RecordBlock: (*journalPtrs)[uint32(journalBlock)*superblock.BlockSize/e.blockSize],
				Description: "Invalid record",
			})
		}

		switch header.BlockType {
		case journal.Header_BlockTypeEnum__DescriptorBlock:
			{
				journalRecords = append(journalRecords, model.JournalRecord{
					Type:        model.JournalDescriptor,
					RecordNum:   uint32(journalBlock),
					RecordBlock: (*journalPtrs)[uint32(journalBlock)*superblock.BlockSize/e.blockSize],
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
						RecordBlock: (*journalPtrs)[uint32(journalBlock+i+1)*superblock.BlockSize/e.blockSize],
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
					RecordBlock: (*journalPtrs)[uint32(journalBlock)*superblock.BlockSize/e.blockSize],
					Description: "Commit block",
				})
			}
		case journal.Header_BlockTypeEnum__RevokeBlock:
			{
				journalRecords = append(journalRecords, model.JournalRecord{
					Type:        model.JournalRevoke,
					RecordNum:   uint32(journalBlock),
					RecordBlock: (*journalPtrs)[uint32(journalBlock)*superblock.BlockSize/e.blockSize],
					Description: "Revoke block",
				})
			}
		default:
			{
				journalRecords = append(journalRecords, model.JournalRecord{
					Type:        model.JournalInvalid,
					RecordNum:   uint32(journalBlock),
					RecordBlock: (*journalPtrs)[uint32(journalBlock)*superblock.BlockSize/e.blockSize],
					Description: "Unknown block",
				})
			}
		}

		journalBlock++
	}

	e.journal = journalRecords
	return nil
}

func (e *Ext3) GetJournal() *[]model.JournalRecord {
	return &e.journal
}

func (e *Ext3) GetFreshInodeFromJournal(inodeNum uint32) (*structure.Inode, error) {
	var freshInode structure.Inode
	inodeBlock, err := e.GetInodeBlockPtr(inodeNum)
	if err != nil {
		return nil, err
	}
	inodeSeek := e.GetInodeSeekInBlock(inodeNum)
	inodeSize := uint32(e.superblock.InodeSize)

	for _, record := range e.journal {
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

			ptrs, err := e.readFilePtrsFromInode(tmpInode)
			if err != nil {
				continue
			}

			flag := false
			for _, ptr := range *ptrs {
				if ptr != 0 {
					flag = true
					break
				}
			}

			if flag {
				freshInode = *tmpInode
			}
		}
	}

	return &freshInode, nil
}
