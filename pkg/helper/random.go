package helper

import "github.com/labstack/gommon/random"

// Random returns a random string with given length.
// charsets is a list of character sets used in random string.
func Random(len uint8, charsets ...string) string {
	if len <= 0 {
		return ""
	}
	return random.String(len, charsets...)
}
