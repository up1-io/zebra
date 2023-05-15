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

func removeRootDir(path string) string {
	return strings.Split(path, "pages/")[1]
}

func convertToURL(filePath string) string {
	out := strings.Replace(filePath, ".gohtml", "", 1)
	out = removeRootDir(out)

	out = strings.Replace(out, "_index", "", -1)
	return fmt.Sprintf("/%s", out)
}

func trimTrailingSlash(url string) string {
	if strings.HasSuffix(url, "/") {
		return url[:len(url)-1]
	}

	return url
}
