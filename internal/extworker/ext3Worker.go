package extworker

import (
	"fmt"
	"os"
	"strings"

	ext3 "github.com/enoly/extFileRecovery/pkg/ext3"
	structure "github.com/enoly/extFileRecovery/pkg/ext3/structure"
)

type Ext3Worker struct {
	ExtFs *ext3.Ext3
}

func NewExt3Worker(ext3fs *ext3.Ext3) *Ext3Worker {
	e := Ext3Worker{}
	e.ExtFs = ext3fs
	return &e
}

func (e *Ext3Worker) GetSuperblockInfo() *SuperblockInfo {
	info := SuperblockInfo{}
	if e.ExtFs != nil {
		info.BlockCount = uint64(e.ExtFs.Superblock.BlocksCount)
		info.InodeCount = uint64(e.ExtFs.Superblock.InodesCount)
		info.BlocksPerGroup = uint64(e.ExtFs.Superblock.BlocksPerGroup)
		info.InodesPerGroup = uint64(e.ExtFs.Superblock.InodesPerGroup)
		info.BlockSize = uint64(e.ExtFs.BlockSize)
		info.InodeSize = uint64(e.ExtFs.Superblock.InodeSize)
		info.JournalInode = uint64(e.ExtFs.Superblock.JournalInodeNum)
		info.UUID = e.ExtFs.Superblock.Uuid
		info.VolumeName = string(e.ExtFs.Superblock.VolumeName)
	}

	return &info
}

func (e *Ext3Worker) ReadDirectory(parent *structure.Directory, path string) (*structure.Directory, error) {
	splittedPath := strings.Split(path, "/")

	if parent == nil && splittedPath[0] == "" {
		root, err := e.ExtFs.ReadDirFromInodePtr(2)
		if err != nil {
			return nil, fmt.Errorf("unable to read root directory: %v", err)
		}

		return e.ReadDirectory(root, path[1:])
	}

	if len(splittedPath) == 1 && splittedPath[0] == "" {
		return parent, nil
	}

	for _, entry := range parent.Entries {
		if entry.FileType == structure.Directory_DirEntry_FileTypeEnum__Dir && entry.Name == splittedPath[0] {
			dir, err := e.ExtFs.ReadDirFromInodePtr(entry.InodePtr)
			if err != nil {
				return nil, fmt.Errorf("unable to read directory %v: %v", entry.Name, err)
			}

			seek := 0
			if len(splittedPath) > 1 {
				seek = 1
			}

			return e.ReadDirectory(dir, path[len(entry.Name)+seek:])
		}
	}

	return nil, fmt.Errorf("unable to find directory")
}

func (e *Ext3Worker) FindInJournal(path string) (*map[string]*structure.Inode, error) {
	foundInodes := make(map[string]*structure.Inode)

	dir, err := e.ReadDirectory(nil, path)
	if err != nil {
		return &foundInodes, err
	}

	deletedFromDir, err := e.ExtFs.GetDeletedFilesFromDir(dir)
	if err != nil {
		return &foundInodes, err
	}

	if len(deletedFromDir.Entries) == 0 {
		return &foundInodes, nil
	}

	for _, entry := range deletedFromDir.Entries {
		inode, err := e.ExtFs.GetFreshInodeFromJournal(entry.InodePtr)
		if err != nil || inode == nil {
			continue
		}

		foundInodes[entry.Name] = inode
	}

	return &foundInodes, nil
}

func (e *Ext3Worker) RestoreFromInode(name string, inode *structure.Inode) error {
	file, err := e.ExtFs.ReadFileFromInode(inode)
	if err != nil {
		return err
	}

	if err := os.WriteFile(name, *file, 0777); err != nil {
		return err
	}

	return nil
}

func (e *Ext3Worker) readAddrBlock(blockPtr uint32) (*structure.IndirectBlock, error) {
	addrBlock, err := e.ExtFs.ReadIndirectAddrBlock(blockPtr)
	if err != nil {
		return nil, err
	}

	tmp := addrBlock.BlocksPtrs[0].Ptr
	for j := 1; j < len(addrBlock.BlocksPtrs); j++ {
		if addrBlock.BlocksPtrs[j].Ptr != 0 {
			if j < 4 && addrBlock.BlocksPtrs[j].Ptr <= tmp {
				return nil, nil
			}

			if addrBlock.BlocksPtrs[j].Ptr > e.ExtFs.Superblock.BlocksCount {
				return nil, nil
			}

			tmp = addrBlock.BlocksPtrs[j].Ptr
		}
	}

	if addrBlock.BlocksPtrs[0].Ptr == 0 {
		return nil, nil
	}

	return addrBlock, nil
}

func (e *Ext3Worker) RestoreFromIndirectBlocks(path string) (*map[string][][]uint32, error) {
	foundIndirect := make(map[string][][]uint32)

	dir, err := e.ReadDirectory(nil, path)
	if err != nil {
		return &foundIndirect, err
	}

	deletedFromDir, err := e.ExtFs.GetDeletedFilesFromDir(dir)
	if err != nil {
		return &foundIndirect, err
	}

	if len(deletedFromDir.Entries) == 0 {
		return &foundIndirect, nil
	}

	for _, entry := range deletedFromDir.Entries {
		fileBlockGroup := e.ExtFs.GetInodeBlockGroup(entry.InodePtr)
		blockBitmap, err := e.ExtFs.GetBlockGroupBitmap(fileBlockGroup)
		if err != nil {
			return &foundIndirect, err
		}

		inodeTableSize := e.ExtFs.Superblock.InodesPerGroup * uint32(e.ExtFs.Superblock.InodeSize) / e.ExtFs.BlockSize
		blocksStart := e.ExtFs.GDT[fileBlockGroup].InodeTableBlock + inodeTableSize
		for i := blocksStart; i < (fileBlockGroup+1)*e.ExtFs.Superblock.BlocksPerGroup && i < e.ExtFs.Superblock.BlocksCount; i++ {
			if (*blockBitmap)[i-fileBlockGroup*e.ExtFs.Superblock.BlocksPerGroup] {
				addrBlock, err := e.readAddrBlock(uint32(i))
				if err != nil || addrBlock == nil {
					continue
				}

				directBlocks := []uint32{}
				for j := addrBlock.BlocksPtrs[0].Ptr - 12; j < addrBlock.BlocksPtrs[0].Ptr; j++ {
					directBlocks = append(directBlocks, j)
				}

				for _, ptr := range addrBlock.BlocksPtrs {
					if ptr.Ptr == 0 {
						break
					}
					directBlocks = append(directBlocks, ptr.Ptr)
				}
				foundIndirect[entry.Name] = append(foundIndirect[entry.Name], directBlocks)
			}
		}
	}

	return &foundIndirect, nil
}

func (e *Ext3Worker) RestoreFromPtrs(name string, ptrs *[]uint32) error {
	file, err := e.ExtFs.ReadFileFromPtrs(ptrs)
	if err != nil {
		return err
	}

	if err := os.WriteFile(name, *file, 0777); err != nil {
		return err
	}

	return nil
}
