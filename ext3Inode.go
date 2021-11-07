// This is a generated file! Please edit source .ksy file and use kaitai-struct-compiler to rebuild
package main

import "github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"

type Ext3Inode struct {
	Mode                     uint16
	Uid                      uint16
	Size                     uint32
	Atime                    uint32
	Ctime                    uint32
	Mtime                    uint32
	Dtime                    uint32
	Gid                      uint16
	LinksCount               uint16
	Blocks                   uint32
	Flags                    uint32
	Osd1                     uint32
	DirectBlocks             []*Ext3Inode_BlockPtr
	FirstLevelIndirectBlock  *Ext3Inode_BlockPtr
	SecondLevelIndirectBlock *Ext3Inode_BlockPtr
	ThirdLevelIndirectBlock  *Ext3Inode_BlockPtr
	Generation               uint32
	FileAcl                  uint32
	DirAcl                   uint32
	Faddr                    uint32
	Osd2                     []byte
	_io                      *kaitai.Stream
	_root                    *Ext3Inode
	_parent                  interface{}
}

func NewExt3Inode() *Ext3Inode {
	return &Ext3Inode{}
}

func (this *Ext3Inode) Read(io *kaitai.Stream, parent interface{}, root *Ext3Inode) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp1, err := this._io.ReadU2le()
	if err != nil {
		return err
	}
	this.Mode = uint16(tmp1)
	tmp2, err := this._io.ReadU2le()
	if err != nil {
		return err
	}
	this.Uid = uint16(tmp2)
	tmp3, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.Size = uint32(tmp3)
	tmp4, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.Atime = uint32(tmp4)
	tmp5, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.Ctime = uint32(tmp5)
	tmp6, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.Mtime = uint32(tmp6)
	tmp7, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.Dtime = uint32(tmp7)
	tmp8, err := this._io.ReadU2le()
	if err != nil {
		return err
	}
	this.Gid = uint16(tmp8)
	tmp9, err := this._io.ReadU2le()
	if err != nil {
		return err
	}
	this.LinksCount = uint16(tmp9)
	tmp10, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.Blocks = uint32(tmp10)
	tmp11, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.Flags = uint32(tmp11)
	tmp12, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.Osd1 = uint32(tmp12)
	this.DirectBlocks = make([]*Ext3Inode_BlockPtr, 12)
	for i := range this.DirectBlocks {
		tmp13 := NewExt3Inode_BlockPtr()
		err = tmp13.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.DirectBlocks[i] = tmp13
	}
	tmp14 := NewExt3Inode_BlockPtr()
	err = tmp14.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.FirstLevelIndirectBlock = tmp14
	tmp15 := NewExt3Inode_BlockPtr()
	err = tmp15.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.SecondLevelIndirectBlock = tmp15
	tmp16 := NewExt3Inode_BlockPtr()
	err = tmp16.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.ThirdLevelIndirectBlock = tmp16
	tmp17, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.Generation = uint32(tmp17)
	tmp18, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.FileAcl = uint32(tmp18)
	tmp19, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.DirAcl = uint32(tmp19)
	tmp20, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.Faddr = uint32(tmp20)
	tmp21, err := this._io.ReadBytes(int(12))
	if err != nil {
		return err
	}
	tmp21 = tmp21
	this.Osd2 = tmp21
	return err
}

type Ext3Inode_BlockPtr struct {
	Ptr     uint32
	_io     *kaitai.Stream
	_root   *Ext3Inode
	_parent *Ext3Inode
}

func NewExt3Inode_BlockPtr() *Ext3Inode_BlockPtr {
	return &Ext3Inode_BlockPtr{}
}

func (this *Ext3Inode_BlockPtr) Read(io *kaitai.Stream, parent *Ext3Inode, root *Ext3Inode) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp22, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.Ptr = uint32(tmp22)
	return err
}
