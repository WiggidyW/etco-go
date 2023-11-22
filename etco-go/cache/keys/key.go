package keys

import (
	"crypto/md5"
)

type Key struct {
	parent *Key
	Text   string
	Buf    [16]byte
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
		Text:   pfx,
		Buf:    md5.Sum([]byte(pfx)),
	}
}

func newKey(
	parent Key,
	parts ...string,
) Key {
	buf := parent.Buf[:]
	for _, part := range parts {
		buf = append(buf, []byte(part)...)
	}
	return Key{
		parent: &parent,
		Buf:    md5.Sum(buf),
		Text:   string(buf),
	}
}

func NewTypeStr(
	buf []byte,
	desc string,
) Key {
	buf = append(buf, []byte(desc)...)
	return Key{
		parent: nil,
		Buf:    md5.Sum(buf),
		Text:   string(buf),
	}
}
