package serverdebug

import (
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"
)

type page struct {
	Path        string
	Description string
}

type indexPage struct {
	pages []page
}

func newIndexPage() *indexPage {
	return &indexPage{}
}

func (i *indexPage) addPage(path, description string) {
	i.pages = append(i.pages, page{Path: path, Description: description})
}

var indexTmpl = template.Must(template.New("index").Parse(`<html>
<head><title>Chat Service Debug</title></head>
<body>
<h2>Chat Service Debug</h2>
<ul>
{{range .}}<li><a href="{{.Path}}">{{.Path}}</a> {{.Description}}</li>
{{end}}
</ul>
</body>
</html>`))

func (i *indexPage) handler(eCtx echo.Context) error {
	eCtx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	eCtx.Response().WriteHeader(http.StatusOK)
	return indexTmpl.Execute(eCtx.Response(), i.pages)
}
