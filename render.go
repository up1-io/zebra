package zebra

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func (z *Zebra) render404(w http.ResponseWriter) {
	tmpl, err := template.ParseFiles(filepath.Join(z.RootDir, pagesFolderName, "_404.gohtml"))
	if err != nil {
		z.render500(w)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	if err := tmpl.Execute(w, nil); err != nil {
		z.render500(w)
	}
}

func (z *Zebra) render500(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	_, err := w.Write([]byte("500 Internal server error"))
	if err != nil {
		panic(err)
	}
}

func (z *Zebra) renderTemplate(w http.ResponseWriter, result Result, p *Page) {
	var paths []string
	for _, component := range p.Components {
		paths = append(paths, component.TemplatePath)
	}

	paths = append([]string{p.LayoutTemplatePath}, paths...)
	paths = append(paths, p.TemplatePath)

	tmpl, err := template.ParseFiles(paths...)
	if err != nil {
		println(err.Error())
		z.render500(w)
		return
	}

	if err := tmpl.Execute(w, result.Data); err != nil {
		println(err.Error())
		z.render500(w)
		return
	}
}
