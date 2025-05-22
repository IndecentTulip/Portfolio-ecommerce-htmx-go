package temp_render

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

// for Renderer interface 
type Templates struct{
  templates *template.Template
}

func NewTemplate() *Templates {
  return &Templates{
    // path is determined from root ./src
    templates: template.Must(template.ParseGlob("../views/*.html")),
  }
}

// the Render method internally calls ExecuteTemplate, which is defined on *template.Template.
func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
  // apperently *template.Template implements "ExecuteTemplate" 
  // that is needed by Renderer interface
  return t.templates.ExecuteTemplate(w, name, data) 
}


