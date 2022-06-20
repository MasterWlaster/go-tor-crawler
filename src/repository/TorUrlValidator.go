package repository

import (
	"regexp"
	"strings"
)

type TorUrlValidator struct {
}

func NewTorUrlValidator() *TorUrlValidator {
	return &TorUrlValidator{}
}

func (t *TorUrlValidator) IsValid(url string) (bool, error) {
	return regexp.MatchString(`https?://(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.onion\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`, url)
}

func (t *TorUrlValidator) Relative(relative string, url string) (string, error) {
	if r, err := t.isRelative(relative); !r {
		return relative, err
	}

	rp := strings.HasPrefix(relative, "/")
	us := strings.HasSuffix(url, "/")

	switch {
	case !rp && !us:
		return url + "/" + relative, nil
	case rp && us:
		return url + strings.TrimPrefix(relative, "/"), nil
	case !rp && us:
		fallthrough
	case rp && !us:
		return url + relative, nil
	}

	return "", nil
}

func (t *TorUrlValidator) isRelative(url string) (bool, error) {
	return !strings.Contains(url, "http"), nil
}
