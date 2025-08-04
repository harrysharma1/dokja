package util

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

func Sluggify(s string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatalf("Error with sluggifying string %s: %s", s, reg)
	}

	processedString := reg.ReplaceAllString(s, " ")

	processedString = strings.TrimSpace(processedString)

	slug := strings.ReplaceAll(processedString, " ", "-")

	slug = strings.ToLower(slug)
	slug = "/" + slug
	return slug

}

func ParseInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("Error converting string to integer: %s", err)
	}
	return n
}
