package assets

import (
	"io/ioutil"
	"math/rand"
	"path/filepath"
)

type Icon = []byte

var (
	icons []Icon
)

func Init() error {
	matches, err := filepath.Glob("../assets/*.png")
	if err != nil {
		return err
	}

	for _, filename := range matches {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}
		icons = append(icons, b)
	}

	return nil
}

func GetIcon() Icon {
	return icons[rand.Intn(len(icons))]
}
