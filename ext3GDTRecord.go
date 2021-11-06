// This is a generated file! Please edit source .ksy file and use kaitai-struct-compiler to rebuild
package main

import "github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"

type Ext3GdtRecord struct {
	BlockBitmapBlock uint32
	InodeBitmapBlock uint32
	InodeTableBlock  uint32
	FreeBlocksCount  uint16
	FreeInodesCount  uint16
	UsedDirsCount    uint16
	_io              *kaitai.Stream
	_root            *Ext3GdtRecord
	_parent          interface{}
}

func NewExt3GdtRecord() *Ext3GdtRecord {
	return &Ext3GdtRecord{}
}

func (this *Ext3GdtRecord) Read(io *kaitai.Stream, parent interface{}, root *Ext3GdtRecord) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp1, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.BlockBitmapBlock = uint32(tmp1)
	tmp2, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.InodeBitmapBlock = uint32(tmp2)
	tmp3, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.InodeTableBlock = uint32(tmp3)
	tmp4, err := this._io.ReadU2le()
	if err != nil {
		return err
	}
	this.FreeBlocksCount = uint16(tmp4)
	tmp5, err := this._io.ReadU2le()
	if err != nil {
		return err
	}
	this.FreeInodesCount = uint16(tmp5)
	tmp6, err := this._io.ReadU2le()
	if err != nil {
		return err
	}
	this.UsedDirsCount = uint16(tmp6)
	return err
}
