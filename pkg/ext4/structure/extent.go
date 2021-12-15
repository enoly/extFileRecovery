// This is a generated file! Please edit source .ksy file and use kaitai-struct-compiler to rebuild

package structure

import (
	"bytes"

	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"
)

type Extent struct {
	Signature       []byte
	ValidEntriesNum uint16
	MaxEntriesNum   uint16
	Depth           uint16
	Generation      uint32
	InternalNodes   []*Extent_InternalNode
	LeafNodes       []*Extent_LeafNode
	_io             *kaitai.Stream
	_root           *Extent
	_parent         interface{}
}

func NewExtent() *Extent {
	return &Extent{}
}

func (this *Extent) Read(io *kaitai.Stream, parent interface{}, root *Extent) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp1, err := this._io.ReadBytes(int(2))
	if err != nil {
		return err
	}
	tmp1 = tmp1
	this.Signature = tmp1
	if !(bytes.Equal(this.Signature, []uint8{10, 243})) {
		return kaitai.NewValidationNotEqualError([]uint8{10, 243}, this.Signature, this._io, "/seq/0")
	}
	tmp2, err := this._io.ReadU2le()
	if err != nil {
		return err
	}
	this.ValidEntriesNum = uint16(tmp2)
	tmp3, err := this._io.ReadU2le()
	if err != nil {
		return err
	}
	this.MaxEntriesNum = uint16(tmp3)
	tmp4, err := this._io.ReadU2le()
	if err != nil {
		return err
	}
	this.Depth = uint16(tmp4)
	tmp5, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.Generation = uint32(tmp5)
	var tmp6 uint16
	if this.Depth == 0 {
		tmp6 = 0
	} else {
		tmp6 = this.ValidEntriesNum
	}
	this.InternalNodes = make([]*Extent_InternalNode, tmp6)
	for i := range this.InternalNodes {
		tmp7 := NewExtent_InternalNode()
		err = tmp7.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InternalNodes[i] = tmp7
	}
	var tmp8 uint16
	if this.Depth != 0 {
		tmp8 = 0
	} else {
		tmp8 = this.ValidEntriesNum
	}
	this.LeafNodes = make([]*Extent_LeafNode, tmp8)
	for i := range this.LeafNodes {
		tmp9 := NewExtent_LeafNode()
		err = tmp9.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.LeafNodes[i] = tmp9
	}
	return err
}

type Extent_InternalNode struct {
	Block              uint32
	LeafBlockPtrLower  uint32
	LeafBlockPtrHigher uint16
	Unused             []byte
	_io                *kaitai.Stream
	_root              *Extent
	_parent            *Extent
}

func NewExtent_InternalNode() *Extent_InternalNode {
	return &Extent_InternalNode{}
}

func (this *Extent_InternalNode) Read(io *kaitai.Stream, parent *Extent, root *Extent) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp10, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.Block = uint32(tmp10)
	tmp11, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.LeafBlockPtrLower = uint32(tmp11)
	tmp12, err := this._io.ReadU2le()
	if err != nil {
		return err
	}
	this.LeafBlockPtrHigher = uint16(tmp12)
	tmp13, err := this._io.ReadBytes(int(2))
	if err != nil {
		return err
	}
	tmp13 = tmp13
	this.Unused = tmp13
	return err
}

type Extent_LeafNode struct {
	FirstFileBlock   uint32
	CoveredBlocks    uint16
	FirstBlockHigher uint16
	FirstBlockLower  uint32
	_io              *kaitai.Stream
	_root            *Extent
	_parent          *Extent
}

func NewExtent_LeafNode() *Extent_LeafNode {
	return &Extent_LeafNode{}
}

func (this *Extent_LeafNode) Read(io *kaitai.Stream, parent *Extent, root *Extent) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp14, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.FirstFileBlock = uint32(tmp14)
	tmp15, err := this._io.ReadU2le()
	if err != nil {
		return err
	}
	this.CoveredBlocks = uint16(tmp15)
	tmp16, err := this._io.ReadU2le()
	if err != nil {
		return err
	}
	this.FirstBlockHigher = uint16(tmp16)
	tmp17, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.FirstBlockLower = uint32(tmp17)
	return err
}
