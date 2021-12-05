// This is a generated file! Please edit source .ksy file and use kaitai-struct-compiler to rebuild

package journal

import "github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"

type Superblock struct {
	Header                      []byte
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
	_root                       *Superblock
	_parent                     interface{}
}

func NewSuperblock() *Superblock {
	return &Superblock{}
}

func (this *Superblock) Read(io *kaitai.Stream, parent interface{}, root *Superblock) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp1, err := this._io.ReadBytes(int(12))
	if err != nil {
		return err
	}
	tmp1 = tmp1
	this.Header = tmp1
	tmp2, err := this._io.ReadU4be()
	if err != nil {
		return err
	}
	this.BlockSize = uint32(tmp2)
	tmp3, err := this._io.ReadU4be()
	if err != nil {
		return err
	}
	this.BlocksCount = uint32(tmp3)
	tmp4, err := this._io.ReadU4be()
	if err != nil {
		return err
	}
	this.FirstDataBlock = uint32(tmp4)
	tmp5, err := this._io.ReadU4be()
	if err != nil {
		return err
	}
	this.FirstTransactionNumber = uint32(tmp5)
	tmp6, err := this._io.ReadU4be()
	if err != nil {
		return err
	}
	this.FirstTransactionBlock = uint32(tmp6)
	tmp7, err := this._io.ReadU4be()
	if err != nil {
		return err
	}
	this.ErrorNumber = uint32(tmp7)
	tmp8, err := this._io.ReadU4be()
	if err != nil {
		return err
	}
	this.FeatureCompatable = uint32(tmp8)
	tmp9, err := this._io.ReadU4be()
	if err != nil {
		return err
	}
	this.FeatureIncompatable = uint32(tmp9)
	tmp10, err := this._io.ReadU4be()
	if err != nil {
		return err
	}
	this.FeatureReadOnly = uint32(tmp10)
	tmp11, err := this._io.ReadBytes(int(16))
	if err != nil {
		return err
	}
	tmp11 = tmp11
	this.Uuid = tmp11
	tmp12, err := this._io.ReadU4be()
	if err != nil {
		return err
	}
	this.FileSystemCount = uint32(tmp12)
	tmp13, err := this._io.ReadU4be()
	if err != nil {
		return err
	}
	this.SuperblockCopy = uint32(tmp13)
	tmp14, err := this._io.ReadU4be()
	if err != nil {
		return err
	}
	this.JournalBlocksPerTransaction = uint32(tmp14)
	tmp15, err := this._io.ReadU4be()
	if err != nil {
		return err
	}
	this.SystemBlocksPerTransaction = uint32(tmp15)
	tmp16, err := this._io.ReadBytes(int(176))
	if err != nil {
		return err
	}
	tmp16 = tmp16
	this.Unused = tmp16
	tmp17, err := this._io.ReadBytes(int(768))
	if err != nil {
		return err
	}
	tmp17 = tmp17
	this.FsUuids = tmp17
	return err
}
