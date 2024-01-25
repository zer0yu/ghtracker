package url

import (
	"net/url"
	"path"
	"strings"
)

func GetRelativeURL(inputURL string) (string, error) {
	if !strings.HasPrefix(inputURL, "http://") && !strings.HasPrefix(inputURL, "https://") {
		inputURL = "https://" + inputURL
	}

	u, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}

	// 清理路径并删除开头的斜杠
	relativeURL := path.Clean(u.Path)
	relativeURL = strings.TrimPrefix(relativeURL, "/")

	return relativeURL, nil
}
