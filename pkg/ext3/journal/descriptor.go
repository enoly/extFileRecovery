// This is a generated file! Please edit source .ksy file and use kaitai-struct-compiler to rebuild

package journal

import "github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"

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
		if (_it.Flags & 8) == 8 {
			break
		}
	}
	return err
}

type Descriptor_DescriptorRecord struct {
	FsBlockNum uint32
	Flags      uint32
	Uuid       []byte
	_io        *kaitai.Stream
	_root      *Descriptor
	_parent    *Descriptor
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
	tmp4, err := this._io.ReadU4be()
	if err != nil {
		return err
	}
	this.Flags = uint32(tmp4)
	var tmp5 int8
	if (this.Flags & 2) == 2 {
		tmp5 = 0
	} else {
		tmp5 = 16
	}
	tmp6, err := this._io.ReadBytes(int(tmp5))
	if err != nil {
		return err
	}
	tmp6 = tmp6
	this.Uuid = tmp6
	return err
}
