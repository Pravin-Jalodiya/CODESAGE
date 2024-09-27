package utils

func MergeTags(existingTags, newTags []string) []string {
	tagSet := make(map[string]struct{})

	for _, tag := range existingTags {
		tagSet[tag] = struct{}{}
	}

	for _, tag := range newTags {
		if _, found := tagSet[tag]; !found {
			existingTags = append(existingTags, tag)
		}
	}

	return existingTags
}
