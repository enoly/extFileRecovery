// This is a generated file! Please edit source .ksy file and use kaitai-struct-compiler to rebuild

package ext3

import "github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"

type Bitmap struct {
	Flags   []bool
	_io     *kaitai.Stream
	_root   *Bitmap
	_parent interface{}
}

func NewBitmap() *Bitmap {
	return &Bitmap{}
}

func (this *Bitmap) Read(io *kaitai.Stream, parent interface{}, root *Bitmap) (err error) {
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
		tmp2, err := this._io.ReadBitsIntBe(1)
		if err != nil {
			return err
		}
		this.Flags = append(this.Flags, tmp2 != 0)
	}
	return err
}
