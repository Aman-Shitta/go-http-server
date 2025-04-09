package utils

import (
	"fmt"
	"strings"
)

func BreakRequestData(s string) (string, string, error) {
	d := strings.Split(s, " ")

	if len(d) != 3 {
		return "", "", fmt.Errorf("args not proper : %s", d)
	}

	if strings.Trim(strings.ToLower(d[2]), "\r\n") != "http/1.1" {
		return "", "", fmt.Errorf("http version not supported : |%s|", strings.Trim(strings.ToLower(d[2]), "\n"))
	}

	return strings.ToUpper(d[0]), d[1], nil
}
