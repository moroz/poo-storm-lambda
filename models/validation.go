package models

import (
	"regexp"
)

var emailRegexp = regexp.MustCompile("^(([^<>()\\[\\]\\\\.,;:\\s@\"]+(\\.[^<>()\\[\\]\\\\.,;:\\s@\"]+)*)|(\".+\"))@((\\[[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}])|(([a-zA-Z\\-0-9]+\\.)+[a-zA-Z]{2,}))$")

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

	return errors
}
