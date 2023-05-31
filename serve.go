package zebra

import (
	"github.com/davecgh/go-spew/spew"
	"log"
	"net/http"
	"path/filepath"
)

const startText = `
███████╗███████╗██████╗ ██████╗  █████╗ 
╚══███╔╝██╔════╝██╔══██╗██╔══██╗██╔══██╗
  ███╔╝ █████╗  ██████╔╝██████╔╝███████║
 ███╔╝  ██╔══╝  ██╔══██╗██╔══██╗██╔══██║
███████╗███████╗██████╔╝██║  ██║██║  ██║
╚══════╝╚══════╝╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝ 
`

func (z *Zebra) ServeMux(r Router) *http.ServeMux {
	z.router = r

	mux := http.NewServeMux()

	for _, page := range z.Pages {
		mux.HandleFunc(page.URL, z.withMiddleware())
	}

	filePath := filepath.Join(z.RootDir, "public")
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir(filePath))))

	return mux
}

func (z *Zebra) ListenAndServe(addr string, r Router) error {
	z.router = r

	mux := http.NewServeMux()

	for _, page := range z.Pages {
		mux.HandleFunc(page.URL, z.withMiddleware())
	}

	filePath := filepath.Join(z.RootDir, "public")
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir(filePath))))

	println(startText)
	log.Printf("Starting Zebra server at %s\n", addr)

	return http.ListenAndServe(addr, mux)
}

func (z *Zebra) withMiddleware() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := z.GetPageByURL(r.URL.Path)
		if err != nil {
			z.render404(w)
			return
		}

		req := Request{
			Request:       *r,
			PathVariables: make(map[string]string),
		}

		result := Result{
			Data: make(map[string]interface{}),
		}

		if p.PathVariables != nil {
			req.PathVariables = getPathVars(p.URL, r.URL.Path)
			spew.Dump(req.PathVariables)
		}

		middleware := z.router.getMiddlewareByURL(p.URL)
		if middleware != nil {
			middleware(req, func(err error, res Result) {
				if err != nil {
					panic(err)
				}

				if res.Redirect != "" {
					http.Redirect(w, r, res.Redirect, http.StatusTemporaryRedirect)
					return
				}

				result.Data = res.Data
			})
		}

		z.renderTemplate(w, result, &p)
	}
}
