package zebra

import (
	"errors"
	"strings"
)

type Component struct {
	Name         string
	TemplatePath string
}

type Page struct {
	Name               string
	TemplatePath       string
	URL                string
	PathVariables      []string
	LayoutTemplatePath string
	Components         []Component
}

func (z *Zebra) GetPageByURL(url string) (Page, error) {
	for _, page := range z.Pages {
		requestedParts := strings.Split(trimTrailingSlash(url), "/")
		pageUrlParts := strings.Split(trimTrailingSlash(page.URL), "/")

		if len(requestedParts) != len(pageUrlParts) {
			continue
		}

		equal := true
		for i, part := range requestedParts {
			if strings.HasPrefix(pageUrlParts[i], "{") && strings.HasSuffix(pageUrlParts[i], "}") {
				continue
			}

			if part != pageUrlParts[i] {
				equal = false
				break
			}
		}

		if !equal {
			continue
		}

		return page, nil
	}

	return Page{}, errors.New("page not found")
}
