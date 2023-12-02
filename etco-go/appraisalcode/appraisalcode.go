package appraisalcode

import (
	"regexp"
)

type CodeChar = byte

var (
	BUYBACK_CHAR CodeChar = 'u'
	SHOP_CHAR    CodeChar = 's'
	HAUL_CHAR    CodeChar = 'h'
)

var Re *regexp.Regexp = regexp.MustCompile("[ush]{1}[0-9a-f]{15}")

type CodeType uint8

const (
	UnknownCode CodeType = iota
	BuybackCode
	ShopCode
	HaulCode
)

// lowercase or bust
func ParseCode(txt string) (string, CodeType) {
	code := Re.FindString(txt)
	if code == "" {
		return "", UnknownCode
	} else if code[0] == BUYBACK_CHAR {
		return code, BuybackCode
	} else if code[0] == SHOP_CHAR {
		return code, ShopCode
	} else /* if code[0] == HAUL_CHAR */ {
		return code, HaulCode
	}
}
