package zebra

import (
	"errors"
	"strings"
)

// Page is a struct that represents a page.
type Page struct {
	Name               string
	URL                string
	PathVariables      []string
	Components         []Component
	TemplatePath       string
	LayoutTemplatePath string
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
