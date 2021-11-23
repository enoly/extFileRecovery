// This is a generated file! Please edit source .ksy file and use kaitai-struct-compiler to rebuild

package ext3

import "github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"

type Inode struct {
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
	DirectBlocks             []*Inode_BlockPtr
	FirstLevelIndirectBlock  *Inode_BlockPtr
	SecondLevelIndirectBlock *Inode_BlockPtr
	ThirdLevelIndirectBlock  *Inode_BlockPtr
	Generation               uint32
	FileAcl                  uint32
	DirAcl                   uint32
	Faddr                    uint32
	Osd2                     []byte
	_io                      *kaitai.Stream
	_root                    *Inode
	_parent                  interface{}
}

func NewInode() *Inode {
	return &Inode{}
}

func (this *Inode) Read(io *kaitai.Stream, parent interface{}, root *Inode) (err error) {
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
	this.DirectBlocks = make([]*Inode_BlockPtr, 12)
	for i := range this.DirectBlocks {
		tmp13 := NewInode_BlockPtr()
		err = tmp13.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.DirectBlocks[i] = tmp13
	}
	tmp14 := NewInode_BlockPtr()
	err = tmp14.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.FirstLevelIndirectBlock = tmp14
	tmp15 := NewInode_BlockPtr()
	err = tmp15.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.SecondLevelIndirectBlock = tmp15
	tmp16 := NewInode_BlockPtr()
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

type Inode_BlockPtr struct {
	Ptr     uint32
	_io     *kaitai.Stream
	_root   *Inode
	_parent *Inode
}

func NewInode_BlockPtr() *Inode_BlockPtr {
	return &Inode_BlockPtr{}
}

func (this *Inode_BlockPtr) Read(io *kaitai.Stream, parent *Inode, root *Inode) (err error) {
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
