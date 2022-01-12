package extworker

import (
	ext3 "github.com/enoly/extFileRecovery/pkg/ext3"
)

type Ext3Worker struct {
	extFs *ext3.Ext3
}

func NewExt3Worker(ext3fs *ext3.Ext3) *Ext3Worker {
	e := Ext3Worker{}
	e.extFs = ext3fs
	return &e
}

func (e *Ext3Worker) GetSuperblockInfo() *SuperblockInfo {
	info := SuperblockInfo{}
	if e.extFs != nil {
		info.BlockCount = uint64(e.extFs.Superblock.BlocksCount)
		info.InodeCount = uint64(e.extFs.Superblock.InodesCount)
		info.BlocksPerGroup = uint64(e.extFs.Superblock.BlocksPerGroup)
		info.InodesPerGroup = uint64(e.extFs.Superblock.InodesPerGroup)
		info.BlockSize = uint64(e.extFs.BlockSize)
		info.InodeSize = uint64(e.extFs.Superblock.InodeSize)
		info.JournalInode = uint64(e.extFs.Superblock.JournalInodeNum)
		info.UUID = e.extFs.Superblock.Uuid
		info.VolumeName = string(e.extFs.Superblock.VolumeName)
	}

	return &info
}
