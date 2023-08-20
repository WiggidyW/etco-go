package jwt

import (
	"fmt"
	"strconv"
)

type JWTResponse struct {
	CharacterId        *int32 // only nil if an error is also returned
	NativeRefreshToken string // the user of this server must receive this token at all costs
}

type jWTClaims struct {
	Audience    string `json:"aud"`
	Issuer      string `json:"iss"`
	Sub         string `json:"sub"`
	CharacterID int32  `json:"-"`
}

func (clm *jWTClaims) Valid() error {
	// validate the parsed claims
	if clm.Issuer != "login.eveonline.com" &&
		clm.Issuer != "https://login.eveonline.com" &&
		clm.Issuer != "http://login.eveonline.com" {
		return fmt.Errorf("jwt: invalid issuer")
	}
	if clm.Audience != "EVE Online" {
		return fmt.Errorf("jwt: invalid audience")
	}
	if clm.Sub == "" {
		return fmt.Errorf("jwt: subject missing or empty")
	}

	// extract the characterID from the subject
	var i int = len(clm.Sub) - 1
	for ; i >= 0; i-- {
		if clm.Sub[i] == ':' {
			break
		}

	}
	id64, err := strconv.ParseInt(clm.Sub[i+1:], 10, 32)
	if err != nil /* || id64 > 2147483647 || id64 < 0 */ {
		return fmt.Errorf("jwt: invalid subject")
	}
	clm.CharacterID = int32(id64)

	// return valid
	return nil
}
