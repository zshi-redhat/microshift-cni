package assets

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"path/filepath"
	"text/template"
)

var templateFuncs = map[string]interface{}{
	"Dir":       filepath.Dir,
	"Sha256sum": func(s string) string { return fmt.Sprintf("%x", sha256.Sum256([]byte(s))) },
}

func renderParamsFromConfig(extra RenderParams) RenderParams {
	params := map[string]interface{}{}
	for k, v := range extra {
		params[k] = v
	}
	return params
}

func renderTemplate(tb []byte, data RenderParams) ([]byte, error) {
	tmpl, err := template.New("").Option("missingkey=error").Funcs(templateFuncs).Parse(string(tb))
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
