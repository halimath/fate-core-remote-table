package web

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/halimath/httputils/requesturi"
)

var (
	//go:embed public
	staticFiles embed.FS
)

func Provide() http.Handler {
	staticFilesFS, err := fs.Sub(staticFiles, "public")
	if err != nil {
		panic(err)
	}

	pathRewriter, err := requesturi.RewritePath(map[string]string{
		"/join/*":    "/",
		"/session/*": "/",
	})
	if err != nil {
		panic(err)
	}

	return requesturi.Middleware(http.FileServer(http.FS(staticFilesFS)), pathRewriter)
}
