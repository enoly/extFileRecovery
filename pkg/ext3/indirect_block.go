// This is a generated file! Please edit source .ksy file and use kaitai-struct-compiler to rebuild

package ext3

import "github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"

type IndirectBlock struct {
	BlocksPtrs []*IndirectBlock_BlockPtr
	_io        *kaitai.Stream
	_root      *IndirectBlock
	_parent    interface{}
}

func NewIndirectBlock() *IndirectBlock {
	return &IndirectBlock{}
}

func (this *IndirectBlock) Read(io *kaitai.Stream, parent interface{}, root *IndirectBlock) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	for i := 1; ; i++ {
		tmp1, err := this._io.EOF()
		if err != nil {
			return err
		}
		if tmp1 {
			break
		}
		tmp2 := NewIndirectBlock_BlockPtr()
		err = tmp2.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.BlocksPtrs = append(this.BlocksPtrs, tmp2)
	}
	return err
}

type IndirectBlock_BlockPtr struct {
	Ptr     uint32
	_io     *kaitai.Stream
	_root   *IndirectBlock
	_parent *IndirectBlock
}

func NewIndirectBlock_BlockPtr() *IndirectBlock_BlockPtr {
	return &IndirectBlock_BlockPtr{}
}

func (this *IndirectBlock_BlockPtr) Read(io *kaitai.Stream, parent *IndirectBlock, root *IndirectBlock) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp3, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.Ptr = uint32(tmp3)
	return err
}
