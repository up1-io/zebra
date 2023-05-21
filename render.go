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
