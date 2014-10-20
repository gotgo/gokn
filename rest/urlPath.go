package rest

import (
	"bytes"
	"fmt"
	"html/template"
	"net/url"
	"regexp"
	"strings"

	"github.com/fatih/structs"
	"github.com/google/go-querystring/query"
)

var paramRegex *regexp.Regexp

func init() {
	paramRegex, _ = regexp.Compile("\\{([a-zA-Z0-9]+)\\}")
}

type UrlPath struct {
	resourceT        string
	keys             []string
	compiledTemplate *template.Template
}

func NewUrlPath(resourceT string) *UrlPath {
	return &UrlPath{
		resourceT: resourceT,
	}
}

func (up *UrlPath) Path(args interface{}) string {
	if args == nil {
		return up.resourceT
	}

	up.compile()
	toClean := structs.Map(args)
	clean := make(map[string]string)
	for k, v := range toClean {
		clean[k] = url.QueryEscape(fmt.Sprintf("%v", v))
	}
	buff := bytes.NewBufferString("")
	up.compiledTemplate.Execute(buff, clean)
	path := buff.String()

	//remove used keys
	for _, k := range up.keys {
		delete(clean, k)
	}

	return queryParams(path, clean)
}

func queryParams(path string, args map[string]string) string {
	if len(args) == 0 {
		return path
	}

	if v, err := query.Values(args); err != nil {
		panic(err)
	} else {
		return path + "?" + v.Encode()
	}
}

func (up *UrlPath) compile() {
	if up.compiledTemplate == nil {
		captures := paramRegex.FindAllStringSubmatch(up.resourceT, -1)
		keys := make([]string, len(captures))
		for i := range captures {
			keys[i] = captures[i][0]
		}
		up.keys = keys
		templateName := up.resourceT
		resourceTemplate := up.resourceAsTemplate()
		compiled := template.Must(template.New(templateName).Parse(resourceTemplate))
		up.compiledTemplate = compiled
	}
}

// prepare converts the typical url template /url/{param} to the html.templates of
// /url/{{.param}}
func (up *UrlPath) resourceAsTemplate() string {
	urlT := up.resourceT
	halfway := strings.Replace(urlT, "{", "{{.", -1)
	final := strings.Replace(halfway, "}", "}}", -1)
	return final
}
