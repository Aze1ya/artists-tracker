package handlers

import (
	"html/template"
	"log"
	"net/http"
)

type ErrorData struct {
	Errortxt    string
	Errorstatus int
}

var (
	a         ApiClient
	allArtist []Artists
)

func Home(w http.ResponseWriter, r *http.Request) {
	if status := checkErrHome(r); status != 0 {
		errorPage(w, http.StatusText(status), status)
		// log.Print(http.StatusText(status))
		return
	}

	if status, err := Convert(); status != 0 {
		errorPage(w, http.StatusText(status), status)
		log.Print(err)
		return
	}

	tmpl, err := template.ParseFiles("ui/html/content.html", "ui/html/base.html", "ui/html/site.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(http.StatusText(http.StatusInternalServerError))
		return
	}

	tmpl.ExecuteTemplate(w, "site", allArtist)
}

func Convert() (int, error) {
	allArtist1, status, err := a.ConvertAllArtist()
	if status != 0 {
		return 500, err
	}
	allArtist = allArtist1
	return 0, nil
}

func Artistdata(w http.ResponseWriter, r *http.Request) {
	if status := checkErrArtist(r); status != 0 {
		errorPage(w, http.StatusText(status), status)
		log.Print(http.StatusText(status))
		return
	}

	artistId := r.URL.Path[8:]

	oneArtist, status, err := a.ConvertOneArtist(artistId)
	if status != 0 {
		errorPage(w, http.StatusText(status), status)
		log.Print(err)
		return
	}

	tmpl, err := template.ParseFiles("ui/html/artist.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(http.StatusText(http.StatusInternalServerError))
		return
	}
	tmpl.Execute(w, oneArtist)
}

func Filterdata(w http.ResponseWriter, r *http.Request) {
	if status := checkErrFilter(r); status != 0 {
		errorPage(w, http.StatusText(status), status)
		log.Print(http.StatusText(status))
		return
	}

	filteredallArtist, status, err := Filter(r, allArtist)
	if status != 0 {
		errorPage(w, http.StatusText(status), status)
		log.Print(err)
		return
	}

	tmpl, err := template.ParseFiles("ui/html/content.html", "ui/html/base.html", "ui/html/site.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(http.StatusText(http.StatusInternalServerError))
		return
	}

	tmpl.ExecuteTemplate(w, "site", filteredallArtist)
}

func errorPage(w http.ResponseWriter, Errortxt string, Errorstatus int) {
	newErrorData := new(ErrorData)
	newErrorData.Errortxt = Errortxt
	newErrorData.Errorstatus = Errorstatus

	w.WriteHeader(Errorstatus)

	tmpl, err := template.ParseFiles("ui/html/error.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(http.StatusText(http.StatusInternalServerError))
		return
	}
	tmpl.Execute(w, newErrorData)
}

func checkErrHome(r *http.Request) int {
	if r.URL.Path != "/" {
		return http.StatusNotFound
	}
	if r.Method != "GET" && r.Method != "OPTIONS" && r.Method != "HEAD" {
		return http.StatusMethodNotAllowed
	}
	return 0
}

func checkErrArtist(r *http.Request) int {
	if r.Method != "GET" && r.Method != "OPTIONS" && r.Method != "HEAD" {
		return http.StatusMethodNotAllowed
	}
	return 0
}

func checkErrFilter(r *http.Request) int {
	if r.URL.Path != "/filters/" {
		return http.StatusNotFound
	}
	if r.Method != "POST" && r.Method != "OPTIONS" && r.Method != "HEAD" {
		return http.StatusMethodNotAllowed
	}
	return 0
}
