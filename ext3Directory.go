// This is a generated file! Please edit source .ksy file and use kaitai-struct-compiler to rebuild
package main

import "github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"

type Ext3Directory struct {
	Entries []*Ext3Directory_DirEntry
	_io     *kaitai.Stream
	_root   *Ext3Directory
	_parent interface{}
}

func NewExt3Directory() *Ext3Directory {
	return &Ext3Directory{}
}

func (this *Ext3Directory) Read(io *kaitai.Stream, parent interface{}, root *Ext3Directory) (err error) {
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
		tmp2 := NewExt3Directory_DirEntry()
		err = tmp2.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Entries = append(this.Entries, tmp2)
	}
	return err
}

type Ext3Directory_DirEntry_FileTypeEnum int

const (
	Ext3Directory_DirEntry_FileTypeEnum__Unknown Ext3Directory_DirEntry_FileTypeEnum = 0
	Ext3Directory_DirEntry_FileTypeEnum__RegFile Ext3Directory_DirEntry_FileTypeEnum = 1
	Ext3Directory_DirEntry_FileTypeEnum__Dir     Ext3Directory_DirEntry_FileTypeEnum = 2
	Ext3Directory_DirEntry_FileTypeEnum__Chrdev  Ext3Directory_DirEntry_FileTypeEnum = 3
	Ext3Directory_DirEntry_FileTypeEnum__Blkdev  Ext3Directory_DirEntry_FileTypeEnum = 4
	Ext3Directory_DirEntry_FileTypeEnum__Fifo    Ext3Directory_DirEntry_FileTypeEnum = 5
	Ext3Directory_DirEntry_FileTypeEnum__Sock    Ext3Directory_DirEntry_FileTypeEnum = 6
	Ext3Directory_DirEntry_FileTypeEnum__Symlink Ext3Directory_DirEntry_FileTypeEnum = 7
)

type Ext3Directory_DirEntry struct {
	InodePtr uint32
	RecLen   uint16
	NameLen  uint8
	FileType Ext3Directory_DirEntry_FileTypeEnum
	Name     string
	Padding  []byte
	_io      *kaitai.Stream
	_root    *Ext3Directory
	_parent  *Ext3Directory
}

func NewExt3Directory_DirEntry() *Ext3Directory_DirEntry {
	return &Ext3Directory_DirEntry{}
}

func (this *Ext3Directory_DirEntry) Read(io *kaitai.Stream, parent *Ext3Directory, root *Ext3Directory) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp3, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.InodePtr = uint32(tmp3)
	tmp4, err := this._io.ReadU2le()
	if err != nil {
		return err
	}
	this.RecLen = uint16(tmp4)
	tmp5, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.NameLen = tmp5
	tmp6, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.FileType = Ext3Directory_DirEntry_FileTypeEnum(tmp6)
	tmp7, err := this._io.ReadBytes(int(this.NameLen))
	if err != nil {
		return err
	}
	tmp7 = tmp7
	this.Name = string(tmp7)
	tmp8, err := this._io.ReadBytes(int(((this.RecLen - uint16(this.NameLen)) - 8)))
	if err != nil {
		return err
	}
	tmp8 = tmp8
	this.Padding = tmp8
	return err
}
