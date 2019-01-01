package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/harry1453/audioQ/project"
	"log"
	"net/http"
)

var mProject *project.Project

func init() {
	mProject = new(project.Project)
	mProject.Init()
	go initialize()
}

func initialize() {
	router := mux.NewRouter()
	router.Handle("/", http.RedirectHandler("/web", http.StatusMovedPermanently))
	router.Handle("/web", http.RedirectHandler("/web/", http.StatusMovedPermanently))
	router.PathPrefix("/web/").Handler(http.StripPrefix("/web/", http.FileServer(http.Dir("web/"))))
	router.HandleFunc("/api/getProject", getProject).Methods("GET")
	router.HandleFunc("/api/playNext", playNext).Methods("GET") // TODO post?
	router.HandleFunc("/api/loadFile", loadFile).Methods("POST")
	router.HandleFunc("/api/saveFile", saveFile).Methods("GET")
	log.Print(http.ListenAndServe(":8888", router))
}

type CommandResponse struct {
	OK    bool
	Error string
}

func sendOK(writer http.ResponseWriter) {
	json.NewEncoder(writer).Encode(CommandResponse{true, ""})
}

func sendError(writer http.ResponseWriter, err error) {
	json.NewEncoder(writer).Encode(CommandResponse{false, err.Error()})
}

func checkProject(writer http.ResponseWriter) bool {
	if mProject != nil {
		return true
	} else {
		writer.WriteHeader(http.StatusNotFound)
		return false
	}
}

func getProject(writer http.ResponseWriter, request *http.Request) {
	if checkProject(writer) {
		json.NewEncoder(writer).Encode(mProject.GetInfo())
	}
}

func playNext(writer http.ResponseWriter, request *http.Request) {
	if checkProject(writer) {
		err := mProject.PlayNext()
		var response CommandResponse
		if err != nil {
			response.Error = err.Error()
		}
		response.OK = err == nil
		json.NewEncoder(writer).Encode(response)
	}
}

func loadFile(writer http.ResponseWriter, request *http.Request) {
	if mProject != nil {
		if err := request.ParseMultipartForm(32 << 20); err != nil {
			sendError(writer, err)
			return
		}
		file, _, err := request.FormFile("audioqProject")
		if err != nil {
			sendError(writer, err)
			return
		}
		newProject := new(project.Project)
		if err := json.NewDecoder(file).Decode(newProject); err != nil {
			sendError(writer, err)
			return
		}
		if err := newProject.Init(); err != nil {
			sendError(writer, err)
			return
		}
		mProject.Close()
		mProject = newProject
		sendOK(writer)
	}
}

func saveFile(writer http.ResponseWriter, request *http.Request) {
	if checkProject(writer) {
		writer.Header().Set("Content-Disposition", "attachment; filename="+mProject.Name+".audioq")
		writer.Header().Set("Content-Type", request.Header.Get("Content-Type"))
		if err := json.NewEncoder(writer).Encode(*mProject); err != nil {
			sendError(writer, err)
			return
		}
	}
}
