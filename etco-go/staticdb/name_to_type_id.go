package staticdb

import (
	kvreader_ "github.com/WiggidyW/etco-go/staticdb/kvreaders_"
)

func NameToTypeId(name string) (
	typeId int32,
	ok bool,
) {
	return kvreader_.KVReaderNameToTypeId.Get(name)
}
