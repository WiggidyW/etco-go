package etcogobucket

import "encoding/binary"

func Int64ToBytes(i int64) [8]byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], uint64(i))
	return buf
}

func Int32PairToBytes(a, b int32) [8]byte {
	var buf [8]byte
	binary.BigEndian.PutUint32(buf[:4], uint32(a))
	binary.BigEndian.PutUint32(buf[4:], uint32(b))
	return buf
}

func BytesToInt32Pair(buf [8]byte) (a, b int32) {
	a = int32(binary.BigEndian.Uint32(buf[:4]))
	b = int32(binary.BigEndian.Uint32(buf[4:]))
	return a, b
}

func Uint16PairToBytes(u1, u2 uint16) [4]byte {
	var buf [4]byte
	binary.BigEndian.PutUint16(buf[:2], u1)
	binary.BigEndian.PutUint16(buf[2:], u2)
	return buf
}

func BytesToUint16Pair(buf [4]byte) (u1, u2 uint16) {
	u1 = binary.BigEndian.Uint16(buf[:2])
	u2 = binary.BigEndian.Uint16(buf[2:])
	return u1, u2
}

func Uint16ToBytes(u uint16) [2]byte {
	var buf [2]byte
	binary.BigEndian.PutUint16(buf[:], u)
	return buf
}

func BytesToUint16(buf [2]byte) uint16 {
	return binary.BigEndian.Uint16(buf[:])
}
