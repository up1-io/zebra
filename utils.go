package zebra

import (
	"fmt"
	"strings"
)

func dupCheck(items []string) error {
	for i, v1 := range items {
		for j, v2 := range items {
			if i == j {
				continue
			}

			if v1 == v2 {
				return fmt.Errorf("duplicate item found %s", v1)
			}
		}
	}

	return nil
}

func convertFilePathToURL(filepath string) string {
	out := strings.Replace(filepath, ".gohtml", "", 1)
	out = strings.Split(out, fmt.Sprintf("%s/", pagesFolderName))[1]
	out = strings.Replace(out, "_index", "", -1)

	return fmt.Sprintf("/%s", out)
}

func trimTrailingSlash(url string) string {
	if strings.HasSuffix(url, "/") {
		return url[:len(url)-1]
	}

	return url
}
