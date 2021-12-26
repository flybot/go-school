package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	wr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	home(wr, req)
	if wr.Code != http.StatusOK {
		t.Errorf("got HTTP status code %d, expected 200", wr.Code)
	}

	if !strings.Contains(wr.Body.String(), "Welcome to the Club!") {
		t.Errorf(
			`response body "%s" does not contain "Welcome to the Club!"`,
			wr.Body.String(),
		)
	}
}

func TestValidator(t *testing.T) {
	var name string
	var email string

	name = "Good name."
	email = "good@email.com"
	errs := validateData(name, email)
	if len(errs) > 0 {
		t.Errorf("Validate data failed with %s, %s", name, email)
	}

	email = "good#email.com"
	errs = validateData(name, email)
	if len(errs) != 1 {
		t.Errorf("Validate data failed with %s, %s", name, email)
	}

	name = "Name2 & wrong"
	errs = validateData(name, email)
	if len(errs) != 2 {
		t.Errorf("Validate data failed with %s, %s", name, email)
	}
}

func TestMemberExists(t *testing.T) {
	membersList = append(membersList, Member{Name: "Abcd", Email: "mmmm@gmail.com", RegDate: time.Now()})

	email := "mmmm@gmail.com"
	r := memberExists(email)
	if r != true {
		t.Errorf("Function memberExists failed with %s", email)
	}
	email = "newmail@gmail.fr"
	r = memberExists(email)
	if r == true {
		t.Errorf("Function memberExists failed with %s", email)
	}
}
