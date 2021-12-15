// This is a generated file! Please edit source .ksy file and use kaitai-struct-compiler to rebuild

package ext3

import (
	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"
)

type DeletedFiles struct {
	Entries []*DeletedFiles_DeletedEntry
	_io     *kaitai.Stream
	_root   *DeletedFiles
	_parent interface{}
}

func NewDeletedFiles() *DeletedFiles {
	return &DeletedFiles{}
}

func (this *DeletedFiles) Read(io *kaitai.Stream, parent interface{}, root *DeletedFiles) (err error) {
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
		tmp2 := NewDeletedFiles_DeletedEntry()
		err = tmp2.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Entries = append(this.Entries, tmp2)
	}
	return err
}

type DeletedFiles_DeletedEntry_FileTypeEnum int

const (
	DeletedFiles_DeletedEntry_FileTypeEnum__Unknown DeletedFiles_DeletedEntry_FileTypeEnum = 0
	DeletedFiles_DeletedEntry_FileTypeEnum__RegFile DeletedFiles_DeletedEntry_FileTypeEnum = 1
	DeletedFiles_DeletedEntry_FileTypeEnum__Dir     DeletedFiles_DeletedEntry_FileTypeEnum = 2
	DeletedFiles_DeletedEntry_FileTypeEnum__Chrdev  DeletedFiles_DeletedEntry_FileTypeEnum = 3
	DeletedFiles_DeletedEntry_FileTypeEnum__Blkdev  DeletedFiles_DeletedEntry_FileTypeEnum = 4
	DeletedFiles_DeletedEntry_FileTypeEnum__Fifo    DeletedFiles_DeletedEntry_FileTypeEnum = 5
	DeletedFiles_DeletedEntry_FileTypeEnum__Sock    DeletedFiles_DeletedEntry_FileTypeEnum = 6
	DeletedFiles_DeletedEntry_FileTypeEnum__Symlink DeletedFiles_DeletedEntry_FileTypeEnum = 7
)

type DeletedFiles_DeletedEntry struct {
	InodePtr uint32
	RecLen   uint16
	NameLen  uint8
	FileType DeletedFiles_DeletedEntry_FileTypeEnum
	Name     string
	Padding  []byte
	_io      *kaitai.Stream
	_root    *DeletedFiles
	_parent  *DeletedFiles
}

func NewDeletedFiles_DeletedEntry() *DeletedFiles_DeletedEntry {
	return &DeletedFiles_DeletedEntry{}
}

func (this *DeletedFiles_DeletedEntry) Read(io *kaitai.Stream, parent *DeletedFiles, root *DeletedFiles) (err error) {
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
	this.FileType = DeletedFiles_DeletedEntry_FileTypeEnum(tmp6)
	tmp7, err := this._io.ReadBytes(int(this.NameLen))
	if err != nil {
		return err
	}
	tmp7 = tmp7
	this.Name = string(tmp7)
	var tmp8 int8
	tmp9 := (this.NameLen + 8) % 4
	if tmp9 < 0 {
		tmp9 += 4
	}
	if (4 - tmp9) == 4 {
		tmp8 = 0
	} else {
		tmp10 := (this.NameLen + 8) % 4
		if tmp10 < 0 {
			tmp10 += 4
		}
		tmp8 = int8(4 - tmp10)
	}
	tmp11, err := this._io.ReadBytes(int(tmp8))
	if err != nil {
		return err
	}
	tmp11 = tmp11
	this.Padding = tmp11
	return err
}
