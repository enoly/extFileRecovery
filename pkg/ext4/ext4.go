package ext4

import (
	"bytes"
	"fmt"
	"math"
	"os"

	structure "github.com/enoly/extFileRecovery/pkg/ext4/structure"
	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"
)

const (
	SUPERBLOCK_SEEK = 1024
	SUPERBLOCK_SIZE = 1024
	GDT_RECORD_SIZE = 32
	ROOT_DIR_INODE  = 2
)

const (
	INCOMPAT_META_BG = 0x10
	INCOMPAT_EXTENTS = 0x40
	INCOMPAT_64BIT   = 0x80
	INCOMPAT_FLEX_BG = 0x200
)

const (
	INODE_EXTENTS_FL = 0x80000
)

type Ext4 struct {
	Drive         string
	BlockSize     uint32
	Superblock    structure.Superblock
	GDT           []structure.GdtRecord
	GroupsPerFlex uint32
	drive         *os.File
}

func New(driveName string) (*Ext4, error) {
	driveFile, err := os.Open(driveName)
	if err != nil {
		return nil, err
	}

	ext4 := &Ext4{
		Drive: driveName,
		drive: driveFile,
	}

	superblock, err := ext4.readSuperblock()
	if err != nil {
		return nil, err
	}

	ext4.Superblock = *superblock
	ext4.GroupsPerFlex = uint32(math.Pow(2, float64(superblock.SLogGroupsPerFlex)))
	ext4.BlockSize = 1024 << superblock.LogBlockSize

	GDT, err := ext4.readGDT()
	if err != nil {
		return nil, err
	}

	ext4.GDT = *GDT

	return ext4, nil
}

func (e *Ext4) Close() {
	e.drive.Close()
}

func (e *Ext4) readFromDrive(buffer *[]byte, seek int64) error {
	_, err := e.drive.ReadAt(*buffer, seek)
	if err != nil {
		return err
	}

	return nil
}

func checkFlag(field uint32, flag uint32) bool {
	return field&flag == flag
}

func (e *Ext4) readSuperblock() (*structure.Superblock, error) {
	buffer := make([]byte, SUPERBLOCK_SIZE)
	if err := e.readFromDrive(&buffer, SUPERBLOCK_SEEK); err != nil {
		return nil, err
	}

	superblock := structure.NewSuperblock()
	if err := superblock.Read(kaitai.NewStream(bytes.NewReader(buffer)), superblock, superblock); err != nil {
		return nil, err
	}

	if checkFlag(superblock.FeatureIncompatable, INCOMPAT_META_BG) {
		return nil, fmt.Errorf("incompatable META_BG enabled")
	}

	if !checkFlag(superblock.FeatureIncompatable, INCOMPAT_EXTENTS) {
		return nil, fmt.Errorf("extents not enabled")
	}

	return superblock, nil
}

func (e *Ext4) readGDT() (*[]structure.GdtRecord, error) {
	blockGroupsCount := e.Superblock.BlocksCount / e.Superblock.BlocksPerGroup
	if e.Superblock.BlocksCount%e.Superblock.BlocksPerGroup != 0 {
		blockGroupsCount++
	}

	gdtSeek := int64((e.Superblock.FirstDataBlock + 1) * e.BlockSize)
	gdtSize := GDT_RECORD_SIZE * blockGroupsCount
	rawGDT := make([]byte, gdtSize)
	if err := e.readFromDrive(&rawGDT, gdtSeek); err != nil {
		return nil, err
	}

	GDT := []structure.GdtRecord{}
	for i := uint32(0); i < blockGroupsCount; i++ {
		record := structure.NewGdtRecord()
		if err := record.Read(kaitai.NewStream(bytes.NewReader(rawGDT[i*GDT_RECORD_SIZE:(i+1)*GDT_RECORD_SIZE])), record, record); err != nil {
			return nil, err
		}
		GDT = append(GDT, *record)
	}

	return &GDT, nil
}

func (e *Ext4) ReadBlock(blockPtr uint32) (*[]byte, error) {
	block := make([]byte, e.BlockSize)
	if err := e.readFromDrive(&block, int64(blockPtr*e.BlockSize)); err != nil {
		return nil, err
	}

	return &block, nil
}

func (e *Ext4) ReadBlockGroupBitmap(blockGroup uint32) (*structure.Bitmap, error) {
	blockBitmapPtr := e.GDT[blockGroup].BlockBitmapBlock
	rawBitmap, err := e.ReadBlock(blockBitmapPtr)
	if err != nil {
		return nil, err
	}

	blockBitmap := structure.NewBitmap()
	if err := blockBitmap.Read(kaitai.NewStream(bytes.NewReader(*rawBitmap)), blockBitmap, blockBitmap); err != nil {
		return nil, err
	}

	return blockBitmap, nil
}

func (e *Ext4) GetInodeBlockGroup(inodeNum uint32) uint32 {
	return uint32(inodeNum / e.Superblock.InodesPerGroup)
}

func (e *Ext4) GetInodeBlockPtr(inodeNum uint32) (uint32, error) {
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

func (e *Ext4) GetInodeSeekInBlock(inodeNum uint32) uint32 {
	inodeNumInTable := inodeNum % e.Superblock.InodesPerGroup
	inodesPerBlock := e.BlockSize / uint32(e.Superblock.InodeSize)
	return (inodeNumInTable - 1) % inodesPerBlock
}

func (e *Ext4) ReadInode(inodeNum uint32) (*structure.Inode, error) {
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

func (e *Ext4) IsBlockUsed(blockPtr uint32) (bool, error) {
	blockGroup := blockPtr / e.Superblock.BlocksPerGroup
	blockBitmap, err := e.ReadBlockGroupBitmap(blockGroup)
	if err != nil {
		return false, err
	}

	return blockBitmap.Flags[blockPtr%e.Superblock.BlocksPerGroup], nil
}

func (e *Ext4) GetExtentFromInode(inode *structure.Inode) (*structure.Extent, error) {
	if !checkFlag(inode.Flags, INODE_EXTENTS_FL) {
		return nil, fmt.Errorf("inode not uses extents")
	}

	extent := structure.NewExtent()
	if err := extent.Read(kaitai.NewStream(bytes.NewReader(inode.IBlock)), extent, extent); err != nil {
		return nil, err
	}

	return extent, nil
}

func (e *Ext4) GetExtentFromBlock(blockPtr uint32) (*structure.Extent, error) {
	block, err := e.ReadBlock(blockPtr)
	if err != nil {
		return nil, err
	}

	extent := structure.NewExtent()
	if err := extent.Read(kaitai.NewStream(bytes.NewReader(*block)), extent, extent); err != nil {
		return nil, err
	}

	return extent, nil
}

func (e *Ext4) ReadFileFromExtent(extent *structure.Extent) (*[]byte, error) {
	if extent == nil || extent.Depth != 0 {
		return nil, fmt.Errorf("extent is inner")
	}

	file := []byte{}
	for _, leaf := range extent.LeafNodes {
		firstBlock := leaf.FirstBlockLower
		length := leaf.CoveredBlocks
		raw := make([]byte, e.BlockSize*uint32(length))

		if err := e.readFromDrive(&raw, int64(firstBlock*e.BlockSize)); err != nil {
			return nil, err
		}

		file = append(file, raw...)
	}

	return &file, nil
}

func (e *Ext4) SaveFileFromExtent(extent *structure.Extent, fileName string) error {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	if extent == nil || extent.Depth != 0 {
		return fmt.Errorf("extent is inner")
	}
	for _, leaf := range extent.LeafNodes {
		firstBlock := leaf.FirstBlockLower
		length := leaf.CoveredBlocks
		raw := make([]byte, e.BlockSize*uint32(length))

		if err := e.readFromDrive(&raw, int64(firstBlock*e.BlockSize)); err != nil {
			return err
		}

		_, err := f.Write(raw)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Ext4) CheckExtentMagic(blockPtr uint32) bool {
	magic := []byte{10, 243}
	blockFirstTwoBytes := make([]byte, 2)
	if err := e.readFromDrive(&blockFirstTwoBytes, int64(blockPtr*e.BlockSize)); err != nil {
		return false
	}

	if bytes.Equal(magic, blockFirstTwoBytes) {
		return true
	}

	return false
}
