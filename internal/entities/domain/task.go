package domain

import "strings"

type Task struct {
	id          int64
	title       string
	description string
	category    string
	tags        []string
}

func NewTask(title, description, category string, tags []string) Task {
	normalizedTags := make([]string, len(tags))
	for i, tag := range tags {
		normalizedTags[i] = strings.TrimSpace(tag)
	}
	return Task{
		title:       normalizeString(title),
		description: normalizeString(description),
		category:    normalizeString(category),
		tags:        normalizedTags,
	}
}

func normalizeString(s string) string {
	return strings.TrimSpace(s)
}
