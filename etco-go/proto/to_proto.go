package proto

type P0Protoable[
	PROTO any,
] interface {
	ToProto() PROTO
}

type P1Protoable[
	PROTO any,
	P0 any,
] interface {
	ToProto(P0) PROTO
}

func P0ToProtoMany[
	PROTO any,
	T P0Protoable[PROTO],
](
	entries []T,
) (
	protoEntries []PROTO,
) {
	protoEntries = make([]PROTO, len(entries))
	for i, entry := range entries {
		protoEntries[i] = entry.ToProto()
	}
	return protoEntries
}

func P1ToProtoMany[
	PROTO any,
	P0 any,
	T P1Protoable[
		PROTO,
		P0,
	],
](
	entries []T,
	p0 P0,
) (
	protoEntries []PROTO,
) {
	protoEntries = make([]PROTO, len(entries))
	for i, entry := range entries {
		protoEntries[i] = entry.ToProto(p0)
	}
	return protoEntries
}
