package zebra

import (
	"path/filepath"
)

type Zebra struct {
	RootDir string
	Pages   []Page
	Router  Router
}

type Option func(*Zebra)

const defaultRootDir = "."

func WithRootDir(dir string) Option {
	return func(z *Zebra) {
		z.RootDir = dir
	}
}

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
