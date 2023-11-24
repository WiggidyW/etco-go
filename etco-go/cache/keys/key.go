package keys

import (
	"crypto/md5"
)

type Key struct {
	parent *Key
	text   string
	buf    [16]byte
}

func (k Key) Bytes16() [16]byte { return k.buf }
func (k Key) Bytes() []byte     { return k.buf[:] }
func (k Key) String() string    { return string(k.buf[:]) }
func (k Key) FullString() string {
	s := k.text
	parent := k.parent
	for parent != nil {
		s = parent.text + s
		parent = parent.parent
	}
	return s
}
func (k Key) PrettyString() string {
	s := k.text
	parent := k.parent
	for parent != nil {
		s = parent.text + "-" + s
		parent = parent.parent
	}
	return s
}

func newPfxKey(
	parts ...string,
) Key {
	text := ""
	for _, part := range parts {
		text += part
	}
	return Key{
		parent: nil,
		text:   text,
		buf:    md5.Sum([]byte(text)),
	}
}

func newKey(
	parent Key,
	parts ...string,
) Key {
	text := ""
	for _, part := range parts {
		text += part
	}
	buf := append(parent.buf[:], text...)
	return Key{
		parent: &parent,
		text:   text,
		buf:    md5.Sum(buf),
	}
}

func NewTypeStr(
	buf []byte,
	desc string,
) Key {
	buf = append(buf, []byte(desc)...)
	return Key{
		parent: nil,
		text:   desc,
		buf:    md5.Sum(buf),
	}
}
