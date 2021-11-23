// This is a generated file! Please edit source .ksy file and use kaitai-struct-compiler to rebuild

package journal

import (
	"bytes"

	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"
)

type Descriptor struct {
	Header      []byte
	Descriptors []*Descriptor_DescriptorRecord
	_io         *kaitai.Stream
	_root       *Descriptor
	_parent     interface{}
}

func NewDescriptor() *Descriptor {
	return &Descriptor{}
}

func (this *Descriptor) Read(io *kaitai.Stream, parent interface{}, root *Descriptor) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp1, err := this._io.ReadBytes(int(12))
	if err != nil {
		return err
	}
	tmp1 = tmp1
	this.Header = tmp1
	for i := 1; ; i++ {
		tmp2 := NewDescriptor_DescriptorRecord()
		err = tmp2.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		_it := tmp2
		this.Descriptors = append(this.Descriptors, _it)
		if _it.Flags.LastRecord == true {
			break
		}
	}
	return err
}

type Descriptor_DescriptorRecord struct {
	FsBlockNum uint32
	Flags      *Descriptor_DescriptorFlags
	Uuid       []byte
	_io        *kaitai.Stream
	_root      *Descriptor
	_parent    *Descriptor
	_raw_Flags []byte
}

func NewDescriptor_DescriptorRecord() *Descriptor_DescriptorRecord {
	return &Descriptor_DescriptorRecord{}
}

func (this *Descriptor_DescriptorRecord) Read(io *kaitai.Stream, parent *Descriptor, root *Descriptor) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp3, err := this._io.ReadU4be()
	if err != nil {
		return err
	}
	this.FsBlockNum = uint32(tmp3)
	tmp4, err := this._io.ReadBytes(int(4))
	if err != nil {
		return err
	}
	tmp4 = tmp4
	this._raw_Flags = tmp4
	_io__raw_Flags := kaitai.NewStream(bytes.NewReader(this._raw_Flags))
	tmp5 := NewDescriptor_DescriptorFlags()
	err = tmp5.Read(_io__raw_Flags, this, this._root)
	if err != nil {
		return err
	}
	this.Flags = tmp5
	tmp6, err := this._io.ReadBytes(int(16))
	if err != nil {
		return err
	}
	tmp6 = tmp6
	this.Uuid = tmp6
	return err
}

type Descriptor_DescriptorFlags struct {
	Reserved             uint64
	LastRecord           bool
	DeletedByTransaction bool
	SameUuid             bool
	SpecialHandling      bool
	_io                  *kaitai.Stream
	_root                *Descriptor
	_parent              *Descriptor_DescriptorRecord
}

func NewDescriptor_DescriptorFlags() *Descriptor_DescriptorFlags {
	return &Descriptor_DescriptorFlags{}
}

func (this *Descriptor_DescriptorFlags) Read(io *kaitai.Stream, parent *Descriptor_DescriptorRecord, root *Descriptor) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp7, err := this._io.ReadBitsIntBe(28)
	if err != nil {
		return err
	}
	this.Reserved = tmp7
	tmp8, err := this._io.ReadBitsIntBe(1)
	if err != nil {
		return err
	}
	this.LastRecord = tmp8 != 0
	tmp9, err := this._io.ReadBitsIntBe(1)
	if err != nil {
		return err
	}
	this.DeletedByTransaction = tmp9 != 0
	tmp10, err := this._io.ReadBitsIntBe(1)
	if err != nil {
		return err
	}
	this.SameUuid = tmp10 != 0
	tmp11, err := this._io.ReadBitsIntBe(1)
	if err != nil {
		return err
	}
	this.SpecialHandling = tmp11 != 0
	return err
}
