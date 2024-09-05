package validation

import (
	"cli-project/pkg/utils/data_cleaning"
	"errors"
	"net/url"
	"strings"
)

func ValidateQuestionLink(link string) (string, error) {
	lowerLink := data_cleaning.CleanString(link)
	parsedURL, err := url.Parse(lowerLink)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" || !strings.Contains(parsedURL.Host, "leetcode.com") {
		return "", errors.New("invalid question link: must be a valid Leetcode link")
	}
	return lowerLink, nil
}
