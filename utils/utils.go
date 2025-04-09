package utils

import (
	"fmt"
	"slices"
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

func RespBody(headers map[string]string, statusCode int, content string, ctype string) string {
	var statusText string

	switch statusCode {
	case 200:
		statusText = "OK"
	case 400:
		statusText = "Bad Request"
	case 404:
		statusText = "Not Found"
	default:
		statusCode = 400
		statusText = "Bad Request"
	}

	encoding := ""

	encodingHeaders, ok := headers["Accept-Encoding"]

	if ok {
		encodings := strings.Split(encodingHeaders, ", ")

		if len(encodings) > 0 {
			if slices.Contains(encodings, "gzip") {
				encoding = "gzip"
				return fmt.Sprintf(
					"HTTP/1.1 %d %s\r\nContent-Type: %s\r\nContent-Encoding: %s\r\nContent-Length: %d\r\n\r\n%s",
					statusCode,
					statusText,
					ctype,
					encoding,
					len(content),
					content,
				)
			}
		}
	}

	return fmt.Sprintf(
		"HTTP/1.1 %d %s\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s",
		statusCode,
		statusText,
		ctype,
		len(content),
		content,
	)
}
