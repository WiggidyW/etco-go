package protoutil

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func NewPBObfuscateCharacterID(characterID int32) string {
	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%x", characterID)))
	return hex.EncodeToString(hasher.Sum(nil))
}
