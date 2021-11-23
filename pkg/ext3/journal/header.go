// This is a generated file! Please edit source .ksy file and use kaitai-struct-compiler to rebuild

package journal

import "github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"

type Header_BlockTypeEnum int

const (
	Header_BlockTypeEnum__DescriptorBlock Header_BlockTypeEnum = 1
	Header_BlockTypeEnum__CommitBlock     Header_BlockTypeEnum = 2
	Header_BlockTypeEnum__SuperblockV1    Header_BlockTypeEnum = 3
	Header_BlockTypeEnum__SuperblockV2    Header_BlockTypeEnum = 4
	Header_BlockTypeEnum__RevokeBlock     Header_BlockTypeEnum = 5
)

type Header struct {
	Signature    uint32
	BlockType    Header_BlockTypeEnum
	SerialNumber uint32
	_io          *kaitai.Stream
	_root        *Header
	_parent      interface{}
}

func NewHeader() *Header {
	return &Header{}
}

func (this *Header) Read(io *kaitai.Stream, parent interface{}, root *Header) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp1, err := this._io.ReadU4be()
	if err != nil {
		return err
	}
	this.Signature = uint32(tmp1)
	tmp2, err := this._io.ReadU4be()
	if err != nil {
		return err
	}
	this.BlockType = Header_BlockTypeEnum(tmp2)
	tmp3, err := this._io.ReadU4be()
	if err != nil {
		return err
	}
	this.SerialNumber = uint32(tmp3)
	return err
}
