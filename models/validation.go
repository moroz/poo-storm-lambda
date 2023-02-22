package models

import (
	"net/url"
	"regexp"
)

var emailRegexp = regexp.MustCompile("^(([^<>()\\[\\]\\\\.,;:\\s@\"]+(\\.[^<>()\\[\\]\\\\.,;:\\s@\"]+)*)|(\".+\"))@((\\[[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}])|(([a-zA-Z\\-0-9]+\\.)+[a-zA-Z]{2,}))$")

func isValidURL(input string) bool {
	u, err := url.Parse(input)
	if err != nil {
		return false
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	return true
}

func (p *CreateCommentInput) Validate() []string {
	errors := make([]string, 0)

	if p.Body == "" {
		errors = append(errors, "Body can't be blank")
	}

	if p.Signature == "" {
		errors = append(errors, "Signature can't be blank")
	}

	if !emailRegexp.Match([]byte(p.Email)) {
		errors = append(errors, "Email is invalid")
	}

	if p.Website != nil && *p.Website != "" && !isValidURL(*p.Website) {
		errors = append(errors, "Website URL must be a valid URL with http or https scheme")
	}

	return errors
}
