// This is a generated file! Please edit source .ksy file and use kaitai-struct-compiler to rebuild
package main

import "github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"

type Ext3JournalSuperblock_BlockTypeEnum int

const (
	Ext3JournalSuperblock_BlockTypeEnum__DescriptorBlock Ext3JournalSuperblock_BlockTypeEnum = 1
	Ext3JournalSuperblock_BlockTypeEnum__CommitBlock     Ext3JournalSuperblock_BlockTypeEnum = 2
	Ext3JournalSuperblock_BlockTypeEnum__SuperblockV1    Ext3JournalSuperblock_BlockTypeEnum = 3
	Ext3JournalSuperblock_BlockTypeEnum__SuperblockV2    Ext3JournalSuperblock_BlockTypeEnum = 4
	Ext3JournalSuperblock_BlockTypeEnum__RevokeBlock     Ext3JournalSuperblock_BlockTypeEnum = 5
)

type Ext3JournalSuperblock struct {
	Signature                   uint32
	BlockType                   Ext3JournalSuperblock_BlockTypeEnum
	SerialNumber                uint32
	BlockSize                   uint32
	BlocksCount                 uint32
	FirstDataBlock              uint32
	FirstTransactionNumber      uint32
	FirstTransactionBlock       uint32
	ErrorNumber                 uint32
	FeatureCompatable           uint32
	FeatureIncompatable         uint32
	FeatureReadOnly             uint32
	Uuid                        []byte
	FileSystemCount             uint32
	SuperblockCopy              uint32
	JournalBlocksPerTransaction uint32
	SystemBlocksPerTransaction  uint32
	Unused                      []byte
	FsUuids                     []byte
	_io                         *kaitai.Stream
	_root                       *Ext3JournalSuperblock
	_parent                     interface{}
}

func NewExt3JournalSuperblock() *Ext3JournalSuperblock {
	return &Ext3JournalSuperblock{}
}

func (this *Ext3JournalSuperblock) Read(io *kaitai.Stream, parent interface{}, root *Ext3JournalSuperblock) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp1, err := this._io.ReadU4be()
	if err != nil {
		return err
	}
	this.Signature = uint32(tmp1)
	tmp2, err := this._io.ReadU4be()
	if err != nil {
		return err
	}
	this.BlockType = Ext3JournalSuperblock_BlockTypeEnum(tmp2)
	tmp3, err := this._io.ReadU4be()
	if err != nil {
		return err
	}
	this.SerialNumber = uint32(tmp3)
	tmp4, err := this._io.ReadU4be()
	if err != nil {
		return err
	}
	this.BlockSize = uint32(tmp4)
	tmp5, err := this._io.ReadU4be()
	if err != nil {
		return err
	}
	this.BlocksCount = uint32(tmp5)
	tmp6, err := this._io.ReadU4be()
	if err != nil {
		return err
	}
	this.FirstDataBlock = uint32(tmp6)
	tmp7, err := this._io.ReadU4be()
	if err != nil {
		return err
	}
	this.FirstTransactionNumber = uint32(tmp7)
	tmp8, err := this._io.ReadU4be()
	if err != nil {
		return err
	}
	this.FirstTransactionBlock = uint32(tmp8)
	tmp9, err := this._io.ReadU4be()
	if err != nil {
		return err
	}
	this.ErrorNumber = uint32(tmp9)
	tmp10, err := this._io.ReadU4be()
	if err != nil {
		return err
	}
	this.FeatureCompatable = uint32(tmp10)
	tmp11, err := this._io.ReadU4be()
	if err != nil {
		return err
	}
	this.FeatureIncompatable = uint32(tmp11)
	tmp12, err := this._io.ReadU4be()
	if err != nil {
		return err
	}
	this.FeatureReadOnly = uint32(tmp12)
	tmp13, err := this._io.ReadBytes(int(16))
	if err != nil {
		return err
	}
	tmp13 = tmp13
	this.Uuid = tmp13
	tmp14, err := this._io.ReadU4be()
	if err != nil {
		return err
	}
	this.FileSystemCount = uint32(tmp14)
	tmp15, err := this._io.ReadU4be()
	if err != nil {
		return err
	}
	this.SuperblockCopy = uint32(tmp15)
	tmp16, err := this._io.ReadU4be()
	if err != nil {
		return err
	}
	this.JournalBlocksPerTransaction = uint32(tmp16)
	tmp17, err := this._io.ReadU4be()
	if err != nil {
		return err
	}
	this.SystemBlocksPerTransaction = uint32(tmp17)
	tmp18, err := this._io.ReadBytes(int(176))
	if err != nil {
		return err
	}
	tmp18 = tmp18
	this.Unused = tmp18
	tmp19, err := this._io.ReadBytes(int(768))
	if err != nil {
		return err
	}
	tmp19 = tmp19
	this.FsUuids = tmp19
	return err
}
