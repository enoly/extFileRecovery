package extworker

import (
	"fmt"

	"fyne.io/fyne/v2/data/binding"
	"github.com/enoly/extFileRecovery/pkg/ext4"
	"github.com/enoly/extFileRecovery/pkg/ext4/structure"
)

type Ext4Worker struct {
	ExtFs *ext4.Ext4
}

func NewExt4Worker(ext4fs *ext4.Ext4) *Ext4Worker {
	e := Ext4Worker{}
	e.ExtFs = ext4fs
	return &e
}

func (e *Ext4Worker) GetSuperblockInfo() *SuperblockInfo {
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

func (e *Ext4Worker) FindIndirectExtents(found chan uint64, counter *binding.Float) {
	for i := uint32(1); i < e.ExtFs.Superblock.BlocksCount; i++ {
		if e.ExtFs.CheckExtentMagic(i) {
			extent, err := e.ExtFs.GetExtentFromBlock(i)
			if err != nil || extent == nil {
				continue
			}

			found <- uint64(i)
		}

		if i%500 == 0 {
			(*counter).Set(float64(i) / float64(e.ExtFs.Superblock.BlocksCount))
		}
	}

	(*counter).Set(1)
	close(found)
}

func (e *Ext4Worker) GetExtentFromBlock(block uint64) (*structure.Extent, error) {
	extent, err := e.ExtFs.GetExtentFromBlock(uint32(block))
	if err != nil || extent == nil {
		return nil, err
	}

	return extent, nil
}

func (e *Ext4Worker) RestoreFileFromExtent(extent *structure.Extent, i int) error {
	return e.ExtFs.SaveFileFromExtent(extent, fmt.Sprint(i))
}
