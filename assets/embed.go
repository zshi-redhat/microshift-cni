package embedded

import (
	"embed"
	"path/filepath"
)

//go:embed ovn-kubernetes
var content embed.FS

func Asset(name string) ([]byte, error) {
	return content.ReadFile(name)
}

func List() ([]string, error) {
	out := make([]string, 0)
	files, err := content.ReadDir("ovn-kubernetes")
	if err != nil {
		return out, err
	}
	for _, f := range files {
		fp := filepath.Join("ovn-kubernetes", f.Name())
		out = append(out, fp)
	}
	return out, nil
}

func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}
