package main

import (
	"database/sql"
	"emailparser/internals/models"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/DusanKasan/parsemail"
	_ "modernc.org/sqlite"
)

// This is the core object for the application
// It coordinates everything, such as routing, database handles, etc.
type App struct {
	DB *sql.DB
}

// GetAttachment handles requests from Twilio, when sending SMS/MMS, to download
// an individual attachment It attempts to retrieve the attachment based on its
// ID (id) in the route segment.  If the attachment is available, it's returned,
// otherwise an HTTP 404 Not Found is returned.
func (a *App) GetAttachment(w http.ResponseWriter, r *http.Request) {
	// retrieve the id from the route segment
	id := r.PathValue("id")
	if id == "" {
		log.Println(fmt.Errorf("attachment id was empty"))
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// convert the id to an integer equivalent
	attachmentID, err := strconv.Atoi(id)
	if err != nil {
		log.Println(fmt.Errorf("invalid attachment id"))
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	am := models.AttachmentDataModel{DB: a.DB}
	attachment, err := am.Get(int64(attachmentID))
	if err != nil {
		log.Println(fmt.Errorf("matching attachment could not be found: %v", err))
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	w.Header().Set("content-disposition", fmt.Sprintf("attachment; filename=%s", attachment.Filename))
	w.Write(attachment.File)
}

// HandleSendGridWebhook webhook POST requests from SendGrid's Inbound Parser, in
// response to incoming emails
func (a *App) HandleSendGridWebhook(w http.ResponseWriter, r *http.Request) {
	// Retrieve the POST body's payload
	err := r.ParseForm()
	if err != nil {
		serverError(w, err)
		return
	}

	// Check if the subject line is valid and retrieve the reference ID
	subject := r.PostForm.Get("subject")
	if subject == "" {
		serverError(w, fmt.Errorf("subject not provided"))
		return
	}
	refID, err := getRefIDFromSubject(subject)
	if err != nil {
		serverError(w, err)
		return
	}

	ur := models.UserReferenceDataModel{DB: a.DB}
	// Attempt to retrieve a user with a matching reference ID
	user, err := ur.GetUserFromReference(refID)
	if err != nil {
		serverError(w, err)
		return
	}

	// Instantiate an email object from the POST body
	emailString := r.PostForm.Get("email")
	if emailString == "" {
		serverError(w, fmt.Errorf("email body not provided"))
		return
	}

	reader := strings.NewReader(emailString)
	email, err := parsemail.Parse(reader)
	if err != nil {
		log.Printf("could not parse email")
		serverError(w, err)
	}
	log.Println(email.TextBody)

	// Create a new note on the User's account
	ndm := models.NoteDataModel{DB: a.DB}
	note, err := ndm.Create(email.TextBody, user)
	if err != nil {
		serverError(w, err)
	}
	adm := models.AttachmentDataModel{DB: a.DB}
	for _, a := range email.Attachments {
		_, err := adm.Create(note, a)
		if err != nil {
			log.Printf("could not attach attachment to note. reason: %s", err)
		}
	}

	// Send a new notification
	io.WriteString(w, fmt.Sprintf("Retrieved user %s from ref %s.\n", user.Name, refID))
}

// main is the core point of the application
func main() {
	// open a connection to the database
	db, err := sql.Open("sqlite", "app.db")
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := App{DB: db}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /attachment/{id}", app.GetAttachment)
	mux.HandleFunc("POST /", app.HandleSendGridWebhook)

	log.Print("Starting server on :4000")
	err = http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

// serverError is a utility function that handles returning server errors from
// the application;
func serverError(w http.ResponseWriter, err error) {
	fmt.Println(err.Error())
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// getRefIDFromSubject attempts to extract a reference (ref) ID from an email's subject
// It does this by first matching the subject string against a regular expression.
// If the expression matches, then it extracts the reference id and returns it.
// If either of two previous checks fails, an error is returned.
func getRefIDFromSubject(subject string) (string, error) {
	pattern := "^(?i:Ref(?i:erence)? ID: )(?P<refid>[0-9a-zA-Z]{14})$"
	subjectRegex := regexp.MustCompile(pattern)
	if !subjectRegex.MatchString(subject) {
		return "", fmt.Errorf("%s does not match subject pattern", subject)
	}
	log.Printf("%s matches pattern %s\n", subject, pattern)
	matches := subjectRegex.FindStringSubmatch(subject)
	return matches[1], nil
}
