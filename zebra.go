// Package zebra is a web framework for Go. It is designed to be simple and easy to use.
//
// To get started, create a new Zebra instance and serve it:
//
//	func main() {
//			app, err := zebra.New()
//			if err != nil {
//				log.Fatalln(err)
//			}
//
//			log.Fatalln(app.ListenAndServe(":8080"))
//	}
package zebra

import (
	"path/filepath"
)

// Zebra is the main struct of the framework.
type Zebra struct {
	RootDir string
	Pages   []Page
	Router  Router
}

// Option is a function that can be passed to New to configure the Zebra instance.
type Option func(*Zebra)

const defaultRootDir = "."

// WithRootDir sets the root directory of the Zebra instance.
func WithRootDir(dir string) Option {
	return func(z *Zebra) {
		z.RootDir = dir
	}
}

// New creates a new Zebra instance. It will load all pages from the pages directory.
// The pages directory is relative to the root directory.
func New(opt ...Option) (*Zebra, error) {
	r := NewRouter()
	z := &Zebra{
		RootDir: defaultRootDir,
		Router:  r,
	}

	for _, o := range opt {
		o(z)
	}

	dir := filepath.Join(z.RootDir, "pages")
	err := z.loadPagesFromDir(dir)
	if err != nil {
		return nil, err
	}

	return z, nil
}
