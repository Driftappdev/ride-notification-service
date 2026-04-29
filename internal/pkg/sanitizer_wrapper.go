package pkg

import (
	gosanitizer "github.com/driftappdev/libpackage/gosanitizer"
	"strings"
)

func CleanString(v string) string {
	return strings.TrimSpace(gosanitizer.CollapseWhitespace(gosanitizer.NormalizeUnicode(gosanitizer.RemoveNullBytes(v))))
}
func CleanEmail(v string) string {
	return strings.ToLower(strings.TrimSpace(gosanitizer.RemoveNullBytes(v)))
}
func CleanPhone(v string) string { return gosanitizer.SanitizePhone(gosanitizer.RemoveNullBytes(v)) }
func CleanPath(v string) string {
	return gosanitizer.SanitizePath(gosanitizer.RemoveNullBytes(strings.TrimSpace(v)))
}
func CleanURL(v string) string {
	return gosanitizer.SanitizeURL(gosanitizer.RemoveNullBytes(strings.TrimSpace(v)))
}
