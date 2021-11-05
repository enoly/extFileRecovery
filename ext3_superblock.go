// This is a generated file! Please edit source .ksy file and use kaitai-struct-compiler to rebuild
package main

import (
	"bytes"

	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"
)

type Ext3Superblock_StateEnum int

const (
	Ext3Superblock_StateEnum__ValidFs               Ext3Superblock_StateEnum = 1
	Ext3Superblock_StateEnum__ErrorFs               Ext3Superblock_StateEnum = 2
	Ext3Superblock_StateEnum__OrphansBeingRecovered Ext3Superblock_StateEnum = 4
)

type Ext3Superblock_ErrorsEnum int

const (
	Ext3Superblock_ErrorsEnum__ActContinue Ext3Superblock_ErrorsEnum = 1
	Ext3Superblock_ErrorsEnum__ActReadOnly Ext3Superblock_ErrorsEnum = 2
	Ext3Superblock_ErrorsEnum__ActPanic    Ext3Superblock_ErrorsEnum = 3
)

type Ext3Superblock_CreatorOsEnum int

const (
	Ext3Superblock_CreatorOsEnum__Linux   Ext3Superblock_CreatorOsEnum = 0
	Ext3Superblock_CreatorOsEnum__Hurd    Ext3Superblock_CreatorOsEnum = 1
	Ext3Superblock_CreatorOsEnum__Masix   Ext3Superblock_CreatorOsEnum = 2
	Ext3Superblock_CreatorOsEnum__FreeBsd Ext3Superblock_CreatorOsEnum = 3
	Ext3Superblock_CreatorOsEnum__Lites   Ext3Superblock_CreatorOsEnum = 4
)

type Ext3Superblock_MajorVersionEnum int

const (
	Ext3Superblock_MajorVersionEnum__Orignial Ext3Superblock_MajorVersionEnum = 0
	Ext3Superblock_MajorVersionEnum__Dynamic  Ext3Superblock_MajorVersionEnum = 1
)

type Ext3Superblock struct {
	InodesCount         uint32
	BlocksCount         uint32
	ReservedBlocksCount uint32
	FreeBlocksCount     uint32
	FreeInodesCount     uint32
	FirstDataBlock      uint32
	LogBlockSize        uint32
	LogFragSize         uint32
	BlocksPerGroup      uint32
	FragsPerGroup       uint32
	InodesPerGroup      uint32
	MountTime           uint32
	WrittenTime         uint32
	MountCount          uint16
	MaxMountCount       uint16
	Signature           []byte
	FsState             Ext3Superblock_StateEnum
	Errors              Ext3Superblock_ErrorsEnum
	MinorVersion        uint16
	LastCheck           uint32
	CheckInterval       uint32
	CreatorOs           Ext3Superblock_CreatorOsEnum
	MajorVersion        Ext3Superblock_MajorVersionEnum
	DefReservedUid      uint16
	DefReservedGid      uint16
	FirstInode          uint32
	InodeSize           uint16
	BlockGroupCopyLoc   uint16
	FeatureCompatable   uint32
	FeatureIncompatable uint32
	FeatureReadOnly     uint32
	Uuid                []byte
	VolumeName          []byte
	LastMounted         []byte
	AlgoBitmap          uint32
	PreallocBlocks      uint8
	PreallocDirBlocks   uint8
	Padding1            []byte
	JournalUuid         []byte
	JournalInodeNum     uint32
	JournalDevice       uint32
	OrphanInodes        uint32
	_io                 *kaitai.Stream
	_root               *Ext3Superblock
	_parent             interface{}
}

func NewExt3Superblock() *Ext3Superblock {
	return &Ext3Superblock{}
}

func (this *Ext3Superblock) Read(io *kaitai.Stream, parent interface{}, root *Ext3Superblock) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp1, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.InodesCount = uint32(tmp1)
	tmp2, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.BlocksCount = uint32(tmp2)
	tmp3, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.ReservedBlocksCount = uint32(tmp3)
	tmp4, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.FreeBlocksCount = uint32(tmp4)
	tmp5, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.FreeInodesCount = uint32(tmp5)
	tmp6, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.FirstDataBlock = uint32(tmp6)
	tmp7, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.LogBlockSize = uint32(tmp7)
	tmp8, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.LogFragSize = uint32(tmp8)
	tmp9, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.BlocksPerGroup = uint32(tmp9)
	tmp10, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.FragsPerGroup = uint32(tmp10)
	tmp11, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.InodesPerGroup = uint32(tmp11)
	tmp12, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.MountTime = uint32(tmp12)
	tmp13, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.WrittenTime = uint32(tmp13)
	tmp14, err := this._io.ReadU2le()
	if err != nil {
		return err
	}
	this.MountCount = uint16(tmp14)
	tmp15, err := this._io.ReadU2le()
	if err != nil {
		return err
	}
	this.MaxMountCount = uint16(tmp15)
	tmp16, err := this._io.ReadBytes(int(2))
	if err != nil {
		return err
	}
	tmp16 = tmp16
	this.Signature = tmp16
	if !(bytes.Equal(this.Signature, []uint8{83, 239})) {
		return kaitai.NewValidationNotEqualError([]uint8{83, 239}, this.Signature, this._io, "/seq/15")
	}
	tmp17, err := this._io.ReadU2le()
	if err != nil {
		return err
	}
	this.FsState = Ext3Superblock_StateEnum(tmp17)
	tmp18, err := this._io.ReadU2le()
	if err != nil {
		return err
	}
	this.Errors = Ext3Superblock_ErrorsEnum(tmp18)
	tmp19, err := this._io.ReadU2le()
	if err != nil {
		return err
	}
	this.MinorVersion = uint16(tmp19)
	tmp20, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.LastCheck = uint32(tmp20)
	tmp21, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.CheckInterval = uint32(tmp21)
	tmp22, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.CreatorOs = Ext3Superblock_CreatorOsEnum(tmp22)
	tmp23, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.MajorVersion = Ext3Superblock_MajorVersionEnum(tmp23)
	tmp24, err := this._io.ReadU2le()
	if err != nil {
		return err
	}
	this.DefReservedUid = uint16(tmp24)
	tmp25, err := this._io.ReadU2le()
	if err != nil {
		return err
	}
	this.DefReservedGid = uint16(tmp25)
	tmp26, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.FirstInode = uint32(tmp26)
	tmp27, err := this._io.ReadU2le()
	if err != nil {
		return err
	}
	this.InodeSize = uint16(tmp27)
	tmp28, err := this._io.ReadU2le()
	if err != nil {
		return err
	}
	this.BlockGroupCopyLoc = uint16(tmp28)
	tmp29, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.FeatureCompatable = uint32(tmp29)
	tmp30, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.FeatureIncompatable = uint32(tmp30)
	tmp31, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.FeatureReadOnly = uint32(tmp31)
	tmp32, err := this._io.ReadBytes(int(16))
	if err != nil {
		return err
	}
	tmp32 = tmp32
	this.Uuid = tmp32
	tmp33, err := this._io.ReadBytes(int(16))
	if err != nil {
		return err
	}
	tmp33 = tmp33
	this.VolumeName = tmp33
	tmp34, err := this._io.ReadBytes(int(64))
	if err != nil {
		return err
	}
	tmp34 = tmp34
	this.LastMounted = tmp34
	tmp35, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.AlgoBitmap = uint32(tmp35)
	tmp36, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.PreallocBlocks = tmp36
	tmp37, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.PreallocDirBlocks = tmp37
	tmp38, err := this._io.ReadBytes(int(2))
	if err != nil {
		return err
	}
	tmp38 = tmp38
	this.Padding1 = tmp38
	tmp39, err := this._io.ReadBytes(int(16))
	if err != nil {
		return err
	}
	tmp39 = tmp39
	this.JournalUuid = tmp39
	tmp40, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.JournalInodeNum = uint32(tmp40)
	tmp41, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.JournalDevice = uint32(tmp41)
	tmp42, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.OrphanInodes = uint32(tmp42)
	return err
}
