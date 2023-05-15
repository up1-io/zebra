package zebra

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/gotailwindcss/tailwind/twembed"
	"github.com/gotailwindcss/tailwind/twhandler"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

func (z *Zebra) ListenAndServe(addr string) error {
	mux := http.NewServeMux()

	for _, page := range z.Pages {
		middleware := z.Router.getMiddlewareByURL(page.URL)
		println(page.URL, middleware)
		mux.HandleFunc(page.URL, z.withMiddleware())
	}

	filePath := filepath.Join(z.RootDir, "public")
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir(filePath))))
	mux.Handle("/css/", twhandler.New(http.Dir("public/css"), "/css", twembed.New()))

	return http.ListenAndServe(addr, mux)
}

func (z *Zebra) withMiddleware() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := z.GetPageByURL(r.URL.Path)
		if err != nil {
			z.render404(w)
			return
		}

		ctx := Request{
			Request:       *r,
			PathVariables: make(map[string]string),
		}

		result := Result{
			Data: make(map[string]interface{}),
		}

		if p.PathVariables != nil {
			ctx.PathVariables = z.getPathVars(p.URL, r.URL.Path)
			spew.Dump(ctx.PathVariables)
		}

		middleware := z.Router.getMiddlewareByURL(p.URL)
		if middleware != nil {
			middleware(ctx, func(err error, r Result) {
				if err != nil {
					panic(err)
				}

				result.Data = r.Data
			})
		}

		z.renderTemplate(w, result, &p)
	}
}

func (z *Zebra) getPathVars(url string, requestURL string) map[string]string {
	pathVars := make(map[string]string)

	urlParts := splitURL(url)
	requestURLParts := splitURL(requestURL)

	for i, part := range urlParts {
		if strings.HasPrefix(part, "{") {
			key := part[1 : len(part)-1]
			pathVars[key] = requestURLParts[i]
		}
	}

	return pathVars
}

func splitURL(url string) []string {
	return strings.Split(url, "/")
}

func (z *Zebra) renderTemplate(w http.ResponseWriter, result Result, p *Page) {
	tmpl, err := template.ParseFiles(p.LayoutTemplatePath, p.TemplatePath)
	if err != nil {
		z.render500(w)
		return
	}

	if err := tmpl.Execute(w, result.Data); err != nil {
		z.render500(w)
		return
	}
}
