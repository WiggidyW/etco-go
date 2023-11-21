package bucket

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/proto"

	b "github.com/WiggidyW/etco-go-bucket"
)

const (
	AUTH_HASH_SET_EXPIRES_IN time.Duration = 24 * time.Hour
	AUTH_HASH_SET_BUF_CAP    int           = 0
)

func init() {
	keys.TypeStrAuthHashSet = cache.RegisterType[b.AuthHashSet]("authhashset", AUTH_HASH_SET_BUF_CAP)
}

func GetAuthHashSet(
	x cache.Context,
	domain string,
) (
	rep b.AuthHashSet,
	expires time.Time,
	err error,
) {
	return get(
		x,
		func(ctx context.Context) (b.AuthHashSet, error) {
			return client.ReadAuthHashSet(ctx, domain)
		},
		keys.CacheKeyAuthHashSet(domain),
		keys.TypeStrAuthHashSet,
		AUTH_HASH_SET_EXPIRES_IN,
		nil,
	)
}

func ProtoGetAuthHashSet(
	x cache.Context,
	domain string,
) (
	rep *proto.CfgAuthList,
	expires time.Time,
	err error,
) {
	var authHashSet b.AuthHashSet
	authHashSet, expires, err = GetAuthHashSet(x, domain)
	if err == nil {
		rep = AuthHashSetToProto(authHashSet)
	}
	return rep, expires, err
}

func ProtoGetUserAuthHashSet(x cache.Context) (
	rep *proto.CfgAuthList,
	expires time.Time,
	err error,
) {
	return ProtoGetAuthHashSet(x, "user")
}

func ProtoGetAdminAuthHashSet(x cache.Context) (
	rep *proto.CfgAuthList,
	expires time.Time,
	err error,
) {
	return ProtoGetAuthHashSet(x, "admin")
}

func SetAuthHashSet(
	x cache.Context,
	domain string,
	rep b.AuthHashSet,
) (
	err error,
) {
	return set(
		x,
		func(ctx context.Context, rep b.AuthHashSet) error {
			return client.WriteAuthHashSet(ctx, rep, domain)
		},
		keys.CacheKeyAuthHashSet(domain),
		keys.TypeStrAuthHashSet,
		AUTH_HASH_SET_EXPIRES_IN,
		rep,
		nil,
	)
}

func ProtoSetAuthHashSet(
	x cache.Context,
	domain string,
	rep *proto.CfgAuthList,
) (
	err error,
) {
	return SetAuthHashSet(x, domain, AuthHashSetFromProto(rep))
}

func ProtoSetUserAuthHashSet(
	x cache.Context,
	rep *proto.CfgAuthList,
) (
	err error,
) {
	return ProtoSetAuthHashSet(x, "user", rep)
}

func ProtoSetAdminAuthHashSet(
	x cache.Context,
	rep *proto.CfgAuthList,
) (
	err error,
) {
	return ProtoSetAuthHashSet(x, "admin", rep)
}

// To Proto

func keysToSlice[K comparable, V any](m map[K]V) []K {
	s := make([]K, len(m))
	i := 0
	for id := range m {
		s[i] = id
		i++
	}
	return s
}

func sliceToSet(s []int32) map[int32]struct{} {
	m := make(map[int32]struct{}, len(s))
	for _, id := range s {
		m[id] = struct{}{}
	}
	return m
}

func AuthHashSetToProto(ahs b.AuthHashSet) *proto.CfgAuthList {
	return &proto.CfgAuthList{
		PermitCharacterIds:   keysToSlice(ahs.PermitCharacterIds),
		BannedCharacterIds:   keysToSlice(ahs.BannedCharacterIds),
		PermitCorporationIds: keysToSlice(ahs.PermitCorporationIds),
		BannedCorporationIds: keysToSlice(ahs.BannedCorporationIds),
		PermitAllianceIds:    keysToSlice(ahs.PermitAllianceIds),
	}
}

func AuthHashSetFromProto(cfg *proto.CfgAuthList) b.AuthHashSet {
	if cfg == nil {
		return b.AuthHashSet{}
	} else {
		return b.AuthHashSet{
			PermitCharacterIds:   sliceToSet(cfg.PermitCharacterIds),
			BannedCharacterIds:   sliceToSet(cfg.BannedCharacterIds),
			PermitCorporationIds: sliceToSet(cfg.PermitCorporationIds),
			BannedCorporationIds: sliceToSet(cfg.BannedCorporationIds),
			PermitAllianceIds:    sliceToSet(cfg.PermitAllianceIds),
		}
	}
}
