package zebra

import (
	"path/filepath"
)

// Zebra is the main struct of the framework.
type Zebra struct {
	RootDir string
	Pages   []Page
	router  Router
}

// Option is a function that can be passed to New to configure the Zebra instance.
type Option func(*Zebra)

const defaultRootDir = "."

// New creates a new Zebra instance. It will load all pages from the pages directory.
// The pages directory is relative to the root directory.
func New(opt ...Option) (*Zebra, error) {
	r := NewRouter()
	z := &Zebra{
		RootDir: defaultRootDir,
		router:  r,
	}

	for _, o := range opt {
		o(z)
	}

	dir := filepath.Join(z.RootDir, pagesFolderName)
	err := z.loadPagesFromDir(dir)
	if err != nil {
		return nil, err
	}

	return z, nil
}
