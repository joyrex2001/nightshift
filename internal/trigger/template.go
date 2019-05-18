package trigger

import (
	"bytes"
	"os"
	"text/template"
)

// RenderTemplate will render provided template. It will return an error if the
// rendering of the template fails.
func RenderTemplate(templ string, values interface{}) (string, error) {
	var funcs = template.FuncMap{
		"env": templateEnv,
	}
	// render template
	tobj, err := template.New("template").Funcs(funcs).Parse(templ)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = tobj.Execute(buf, values); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// templateEnv implements the {{ .env }} method which will return the value of
// the given environment variable.
func templateEnv(v string) string {
	return os.Getenv(v)
}
