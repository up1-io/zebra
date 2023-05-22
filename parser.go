package zebra

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func (z *Zebra) loadPagesFromDir(path string) error {
	pages, err := parsePageDir(path)
	if err != nil {
		return err
	}

	for _, page := range pages {
		dir := strings.Replace(page.TemplatePath, page.Name, "", 1)
		layoutTemplatePath, err := findRelatedLayoutTemplate(dir, &page)
		if err != nil {
			return err
		}

		page.LayoutTemplatePath = layoutTemplatePath

		components, err := findRequiredComponents(page.TemplatePath)
		if err != nil {
			return err
		}

		layoutComponents, err := findRequiredComponents(page.LayoutTemplatePath)

		components = append(components, layoutComponents...)
		for _, component := range components {
			page.Components = append(page.Components, Component{
				Name:         component,
				TemplatePath: filepath.Join(z.RootDir, componentsFolderName, fmt.Sprintf("%s.gohtml", component)),
			})
		}

		z.Pages = append(z.Pages, page)
	}

	return nil
}

func findRequiredComponents(filePath string) ([]string, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	components := regexp.MustCompile(`{{\s*template\s*"([a-zA-Z0-9-_/]+)"\s*[.a-zA-Z0-9]*?\s*}}`).FindAllStringSubmatch(string(b), -1)
	spew.Dump(components)

	if len(components) == 0 {
		return []string{}, nil
	}

	var out []string
	for _, component := range components {
		if component[1] == "content" {
			continue
		}

		out = append(out, component[1])
	}

	return out, nil
}

func parsePathVariables(url string) ([]string, error) {
	pathParams := regexp.MustCompile(`\{[a-zA-Z0-9]+}`).FindAllString(url, -1)

	if err := dupCheck(pathParams); err != nil {
		return nil, fmt.Errorf("%s in url %s", err.Error(), url)
	}

	return pathParams, nil
}

func parsePageDir(path string) (out []Page, err error) {
	err = filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !isGoHtmlFile(info.Name()) || isZebraSysFile(info.Name()) {
			return nil
		}

		url := convertFilePathToURL(path)
		pathParams, err := parsePathVariables(url)
		if err != nil {
			return err
		}

		out = append(out, Page{
			Name:          info.Name(),
			TemplatePath:  path,
			URL:           url,
			PathVariables: pathParams,
		})

		return nil
	})

	return out, err
}

func isZebraSysFile(filename string) bool {
	return filename == "_layout.gohtml" || filename == "_404.gohtml"
}

func isGoHtmlFile(filename string) bool {
	return strings.HasSuffix(filename, ".gohtml")
}

func findRelatedLayoutTemplate(path string, page *Page) (string, error) {
	dir, err := os.ReadDir(path)
	if err != nil {
		return "", fmt.Errorf("failed to read directory: %s", err)
	}

	for _, file := range dir {
		if file.Name() == "_layout.gohtml" {
			return filepath.Join(path, file.Name()), nil
		}
	}

	if page.LayoutTemplatePath == "" {
		upDirPath := filepath.Dir(path)
		if upDirPath == path {
			return "", fmt.Errorf("template layout not found for %s", page.TemplatePath)
		}
		return findRelatedLayoutTemplate(upDirPath, page)
	}

	return "", fmt.Errorf("layout template not found for %s", page.TemplatePath)
}
