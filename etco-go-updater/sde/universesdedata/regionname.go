package universesdedata

import (
	"regexp"
)

var RE_REGION_NAME = regexp.MustCompile(`([a-z])([A-Z])`)

func fixRegionName(name string) string {
	return RE_REGION_NAME.ReplaceAllString(name, `${1} ${2}`)
}
