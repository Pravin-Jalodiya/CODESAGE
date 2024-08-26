package data_cleaning

import "strings"

func CleanTags(tags string) []string {
	tagList := strings.Split(tags, ",")
	for i, tag := range tagList {
		tagList[i] = CleanString(tag)
	}
	return tagList
}
