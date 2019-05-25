package trigger

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/golang/glog"
	"github.com/joyrex2001/nightshift/internal/scanner"
)

func mixinObjects(settings map[string]string, objects map[string]*scanner.Object) map[string]interface{} {
	templateValues := make(map[string]interface{}, len(settings)+1)
	for key, element := range settings {
		templateValues[key] = element
	}
	if len(objects) > 0 {
		templateValues["objects"] = make(map[string]scanner.Object, len(objects))
		objs, _ := templateValues["objects"].(map[string]scanner.Object)
		for key, element := range objects {
			objs[key] = *element
		}
	}
	return templateValues
}

// RenderTemplate will render provided template. It will return an error if the
// rendering of the template fails.
func RenderTemplate(templ string, values interface{}) (string, error) {
	var funcs = template.FuncMap{
		"env":  templateEnv,
		"add":  templateAdd,
		"now":  templateNow,
		"time": templateTime,
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

// templateNow will return the current time as an epoch.
func templateNow() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}

// templateAdd will add given values to eachother.
func templateAdd(a, b string) string {
	aa, err := strconv.ParseInt(a, 10, 64)
	if err != nil {
		glog.Errorf("invalid value given for add: %s", err)
	}
	bb, err := strconv.ParseInt(b, 10, 64)
	if err != nil {
		glog.Errorf("invalid value given for add: %s", err)
	}
	return fmt.Sprintf("%d", aa+bb)
}

// templateTime will render a timestring for given epoch to given template. If
// template contains  "rfc3339", "ansic" or "unixdate", it will format the
// timestring according to that standard. Otherwhise the template  should
// contain the desired timestring for the reference time:
// Mon Jan 2 15:04:05 -0700 MST 2006
func templateTime(template, epoch string) string {
	ep, err := strconv.ParseInt(epoch, 10, 64)
	if err != nil {
		glog.Errorf("invalid time given for time function: %s", err)
	}
	switch strings.ToLower(template) {
	case "rfc3339":
		template = time.RFC3339
	case "ansic":
		template = time.ANSIC
	case "unixdate":
		template = time.UnixDate
	}
	t := time.Unix(ep, 0)
	return t.Format(template)
}
