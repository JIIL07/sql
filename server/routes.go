package server

import (
	"encoding/json"
	"net/http"

	cloudfiles "github.com/JIIL07/cloudFiles-manager/client"
	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/user", app.UserHandler)
	router.HandlerFunc(http.MethodPost, "/adduser", app.AddUser)
	router.HandlerFunc(http.MethodDelete, "/deleteuser", app.DeleteUser)
	router.HandlerFunc(http.MethodGet, "/", app.TextHandler)
	router.HandlerFunc(http.MethodGet, "/files", app.SetFilesHandler)
	router.HandlerFunc(http.MethodPost, "/addfiles", app.AddFileHandler)
	router.HandlerFunc(http.MethodGet, "/deletedfiles", app.SetDeletedFilesHandler)

	return router
}

func (app *application) TextHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(`Welcome to the cloudfiles API. You can use the following endpoints:
GET /files
GET /deletedfiles
POST /addfiles
POST /deletefiles
POST /updatefiles
`))

	app.logger.Printf("Server detected / entering")
}

func (app *application) SetFilesHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) SetDeletedFilesHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) AddFileHandler(w http.ResponseWriter, r *http.Request) {
	var fileContext cloudfiles.FileContext
	if err := json.NewDecoder(r.Body).Decode(&fileContext); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		app.errlogger.Printf("error decoding: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	app.logger.Printf("Successfully written file")
}
