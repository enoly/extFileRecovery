// This is a generated file! Please edit source .ksy file and use kaitai-struct-compiler to rebuild

package ext3

import "github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"

type Directory struct {
	Entries []*Directory_DirEntry
	_io     *kaitai.Stream
	_root   *Directory
	_parent interface{}
}

func NewDirectory() *Directory {
	return &Directory{}
}

func (this *Directory) Read(io *kaitai.Stream, parent interface{}, root *Directory) (err error) {
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
		tmp2 := NewDirectory_DirEntry()
		err = tmp2.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Entries = append(this.Entries, tmp2)
	}
	return err
}

type Directory_DirEntry_FileTypeEnum int

const (
	Directory_DirEntry_FileTypeEnum__Unknown Directory_DirEntry_FileTypeEnum = 0
	Directory_DirEntry_FileTypeEnum__RegFile Directory_DirEntry_FileTypeEnum = 1
	Directory_DirEntry_FileTypeEnum__Dir     Directory_DirEntry_FileTypeEnum = 2
	Directory_DirEntry_FileTypeEnum__Chrdev  Directory_DirEntry_FileTypeEnum = 3
	Directory_DirEntry_FileTypeEnum__Blkdev  Directory_DirEntry_FileTypeEnum = 4
	Directory_DirEntry_FileTypeEnum__Fifo    Directory_DirEntry_FileTypeEnum = 5
	Directory_DirEntry_FileTypeEnum__Sock    Directory_DirEntry_FileTypeEnum = 6
	Directory_DirEntry_FileTypeEnum__Symlink Directory_DirEntry_FileTypeEnum = 7
)

type Directory_DirEntry struct {
	InodePtr uint32
	RecLen   uint16
	NameLen  uint8
	FileType Directory_DirEntry_FileTypeEnum
	Name     string
	Padding  []byte
	_io      *kaitai.Stream
	_root    *Directory
	_parent  *Directory
}

func NewDirectory_DirEntry() *Directory_DirEntry {
	return &Directory_DirEntry{}
}

func (this *Directory_DirEntry) Read(io *kaitai.Stream, parent *Directory, root *Directory) (err error) {
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
	this.FileType = Directory_DirEntry_FileTypeEnum(tmp6)
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
