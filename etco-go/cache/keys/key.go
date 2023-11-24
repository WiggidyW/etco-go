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
	fullString := k.FullString()
	buf := make([]byte, len(fullString))
	var i int
	for i = len(fullString) - 1; i >= 0; i-- {
		char := fullString[i]
		if char > 127 {
			i++
			break
		} else {
			buf[i] = char
		}
	}
	return string(buf[i:])
}

func newPfxKey(
	parts ...string,
) Key {
	pfx := ""
	for _, part := range parts {
		pfx += part
	}
	return Key{
		parent: nil,
		text:   pfx,
		buf:    md5.Sum([]byte(pfx)),
	}
}

func newKey(
	parent Key,
	parts ...string,
) Key {
	buf := parent.buf[:]
	for _, part := range parts {
		buf = append(buf, []byte(part)...)
	}
	return Key{
		parent: &parent,
		buf:    md5.Sum(buf),
		text:   string(buf),
	}
}

func NewTypeStr(
	buf []byte,
	desc string,
) Key {
	buf = append(buf, []byte(desc)...)
	return Key{
		parent: nil,
		buf:    md5.Sum(buf),
		text:   string(buf),
	}
}
