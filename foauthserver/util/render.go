package util

import (
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/oxtoacart/bpool"
)

var (
	templates map[string]*template.Template
	bufpool   *bpool.BufferPool
)

// renderTemplate is a wrapper around template.ExecuteTemplate.
// It writes into a bytes.Buffer before writing to the http.ResponseWriter to catch
// any errors resulting from populating the template.
func RenderTemplate(c *gin.Context, name string, data map[string]interface{}) error {
	// Ensure the template exists in the map.
	tmpl, ok := templates[name]
	if !ok {
		return fmt.Errorf("the template (%s) does not exist", name)
	}

	// Create a buffer to temporarily write to and check if any errors were encounterd.
	buf := bufpool.Get()
	defer bufpool.Put(buf)

	err := tmpl.ExecuteTemplate(buf, "base", data)
	if err != nil {
		return err
	}

	// The X-Frame-Options HTTP response header can be used to indicate whether
	// or not a browser should be allowed to render a page in a <frame>,
	// <iframe> or <object> . Sites can use this to avoid clickjacking attacks,
	// by ensuring that their content is not embedded into other sites.
	c.Header("X-Frame-Options", "deny")
	// Set the header and write the buffer to the http.ResponseWriter
	c.Header("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(c.Writer)
	return nil
}

func LoadTemplates(basePath string) {
	templates = make(map[string]*template.Template)

	bufpool = bpool.NewBufferPool(64)

	layoutTemplates := map[string][]string{
		basePath + "template/layouts/outside.html": {
			basePath + "template/includes/register.html",
			basePath + "template/includes/login.html",
		},
		basePath + "template/layouts/inside.html": {
			basePath + "template/includes/authorize.html",
		},
	}

	for layout, includes := range layoutTemplates {
		for _, include := range includes {
			files := []string{include, layout}
			templates[filepath.Base(include)] = template.Must(template.ParseFiles(files...))
		}
	}
}
