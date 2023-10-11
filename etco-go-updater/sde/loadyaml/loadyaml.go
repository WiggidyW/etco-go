package loadyaml

import (
	"os"

	"gopkg.in/yaml.v3"
)

// func Load[T any](
// 	root string,
// 	affix string,
// ) (t T, err error) {
// 	reader, err := os.Open(fmt.Sprintf("%s/%s", root, affix))
// 	if err != nil {
// 		return t, err
// 	}
// 	err = yaml.NewDecoder(reader).Decode(&t)
// 	if err != nil {
// 		return t, err
// 	}
// 	return t, nil
// }

func LoadYaml[T any](path string) (t T, err error) {
	reader, err := os.Open(path)
	if err != nil {
		return t, err
	}
	defer reader.Close()
	err = yaml.NewDecoder(reader).Decode(&t)
	if err != nil {
		return t, err
	}
	return t, nil
}
