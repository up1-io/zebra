package zebra

import (
	"fmt"
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
		dir := filepath.Dir(page.TemplatePath)
		layoutTemplatePath, err := findRelatedLayoutTemplate(dir)
		if err != nil {
			return err
		}

		page.LayoutTemplatePath = layoutTemplatePath

		templateComponents, err := z.findRequiredComponents(page.TemplatePath)
		if err != nil {
			return err
		}

		layoutComponents, err := z.findRequiredComponents(page.LayoutTemplatePath)
		if err != nil {
			return err
		}

		// ToDo: Check if component has subcomponents and search for circular dependencies. Add Behavior for duplicate components.
		page.Components = append(templateComponents, layoutComponents...)

		z.Pages = append(z.Pages, page)
	}

	return nil
}

func (z *Zebra) findRequiredComponents(path string) ([]Component, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	regex := regexp.MustCompile(`{{\s*template\s*"([a-zA-Z0-9-_/]+)"\s*[.a-zA-Z0-9]*?\s*}}`)
	allStringSubmatch := regex.FindAllStringSubmatch(string(b), -1)

	var out []Component
	for _, submatch := range allStringSubmatch {
		name := submatch[1]

		// Skip content template definition
		if name == "content" {
			continue
		}

		out = append(out, Component{
			Name:         name,
			TemplatePath: filepath.Join(z.RootDir, componentsFolderName, fmt.Sprintf("%s.gohtml", name)),
		})
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

func findRelatedLayoutTemplate(path string) (string, error) {
	var out string
	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() == "_layout.gohtml" {
			out = path
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to read directory: %s", err)
	}

	if out != "" {
		return out, nil
	}

	parentDirPath := filepath.Dir(path)
	if parentDirPath == path {
		return "", fmt.Errorf("template layout not found for %s", path)
	}

	return findRelatedLayoutTemplate(parentDirPath)
}
