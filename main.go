package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

type Member struct {
	Name    string
	Email   string
	RegDate time.Time
}

type ViewData struct {
	Title   string
	Warn    bool
	Errors  []string
	Members []Member
}

var membersList []Member

func memberExists(email string) bool {
	for _, a := range membersList {
		if a.Email == email {
			return true
		}
	}
	return false
}

func validateData(name string, email string) []string {
	errs := []string{}
	regexpEmail := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	regexpName := regexp.MustCompile("^[a-zA-Z .]*$")
	// check the name
	if !regexpName.Match([]byte(name)) {
		errs = append(errs, "The name field must contain only English letters, spaces and dots")
	}
	// check is email valid
	if !regexpEmail.Match([]byte(email)) {
		errs = append(errs, "The email field should be a valid email address!")
	}
	return errs
}

func logging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		f(w, r)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	errs := []string{}
	memberAlreadyExists := false
	if r.Method == http.MethodPost {
		m := Member{
			Name:    r.FormValue("username"),
			Email:   r.FormValue("email"),
			RegDate: time.Now(),
		}
		if memberExists(m.Email) {
			memberAlreadyExists = true
		} else {
			errs = validateData(m.Name, m.Email)
			if len(errs) < 1 {
				membersList = append(membersList, m)
			}
		}
	}

	data := ViewData{
		Title:   "Member Club",
		Warn:    memberAlreadyExists,
		Errors:  errs,
		Members: membersList,
	}
	//Log returned data
	log.Printf("%v", data)
	//process template file
	tmpl, _ := template.ParseFiles("page.html")
	tmpl.Execute(w, data)
}

func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/", logging(home))
	fmt.Println("Server is listening...")
	http.ListenAndServe(":"+port, nil)
}
