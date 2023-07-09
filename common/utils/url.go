package utils

import (
	"net/url"
	"strings"
)

func ParseHost(target string) string {
	target = strings.TrimSpace(target)
	if strings.Contains(target, "http") {
		u, err := url.Parse(target)
		if err != nil {
			return ""
		}
		return u.Hostname()
	} else {
		return strings.TrimSpace(strings.Trim(target, "/"))
	}
}
