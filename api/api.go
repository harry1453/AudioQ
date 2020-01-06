package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/harry1453/audioQ/project"
	"log"
	"net/http"
	"strconv"
)

func Initialize() {
	router := mux.NewRouter()
	router.Handle("/", http.RedirectHandler("/web", http.StatusMovedPermanently))
	router.Handle("/web", http.RedirectHandler("/web/", http.StatusMovedPermanently))
	router.PathPrefix("/web/").Handler(http.StripPrefix("/web/", http.FileServer(http.Dir("web/"))))
	router.HandleFunc("/api/getProject", getProject).Methods("GET")
	router.HandleFunc("/api/addCue", addCue).Methods("POST")
	router.HandleFunc("/api/removeCue/{cueNumber}", removeCue).Methods("POST")
	router.HandleFunc("/api/renameCue/{cueNumber}/{cueName}", renameCue).Methods("POST")
	router.HandleFunc("/api/moveCue/{from}/{to}", moveCue).Methods("GET", "POST") // TODO post only?
	router.HandleFunc("/api/playNext", playNext).Methods("POST")
	router.HandleFunc("/api/jumpTo/{cueNumber}", jumpTo).Methods("POST")
	router.HandleFunc("/api/stopPlaying", stopPlaying).Methods("POST")
	router.HandleFunc("/api/updateProjectName/{name}", updateProjectName).Methods("POST")
	router.HandleFunc("/api/updateProjectSettings", updateProjectSettings).Methods("POST")
	router.HandleFunc("/api/loadProject", loadProject).Methods("POST")
	router.HandleFunc("/api/saveProject", saveProject).Methods("GET")
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
	if project.Instance != nil {
		return true
	} else {
		writer.WriteHeader(http.StatusNotFound)
		return false
	}
}

func getProject(writer http.ResponseWriter, request *http.Request) {
	if checkProject(writer) {
		json.NewEncoder(writer).Encode(project.Instance.GetInfo())
	}
}

func addCue(writer http.ResponseWriter, request *http.Request) {
	if project.Instance != nil {
		if err := request.ParseMultipartForm(32 << 20); err != nil {
			sendError(writer, err)
			return
		}
		file, header, err := request.FormFile("audioFile")
		if err != nil {
			sendError(writer, err)
			return
		}
		if err := project.Instance.AddCue(request.FormValue("cueName"), header.Filename, file); err != nil {
			sendError(writer, err)
			return
		}
		sendOK(writer)
	}
}

func removeCue(writer http.ResponseWriter, request *http.Request) {
	if project.Instance != nil {
		cueNumber, err := strconv.Atoi(mux.Vars(request)["cueNumber"])
		if err != nil {
			sendError(writer, err)
			return
		}
		if err := project.Instance.RemoveCue(cueNumber); err != nil {
			sendError(writer, err)
		} else {
			sendOK(writer)
		}
	}
}

func renameCue(writer http.ResponseWriter, request *http.Request) {
	if checkProject(writer) {
		vars := mux.Vars(request)
		cueNumber, err := strconv.Atoi(vars["cueNumber"])
		if err != nil {
			sendError(writer, err)
			return
		}
		if err := project.Instance.RenameCue(cueNumber, vars["cueName"]); err != nil {
			sendError(writer, err)
		} else {
			sendOK(writer)
		}
	}
}

func moveCue(writer http.ResponseWriter, request *http.Request) {
	if checkProject(writer) {
		vars := mux.Vars(request)
		from, err := strconv.Atoi(vars["from"])
		if err != nil {
			sendError(writer, err)
			return
		}
		to, err := strconv.Atoi(vars["to"])
		if err != nil {
			sendError(writer, err)
			return
		}
		if err := project.Instance.MoveCue(from, to); err != nil {
			sendError(writer, err)
		} else {
			sendOK(writer)
		}
	}
}

func playNext(writer http.ResponseWriter, request *http.Request) {
	if checkProject(writer) {
		if err := project.Instance.PlayNext(); err != nil {
			sendError(writer, err)
		} else {
			sendOK(writer)
		}
	}
}

func jumpTo(writer http.ResponseWriter, request *http.Request) {
	if checkProject(writer) {
		if cueNumber, err := strconv.Atoi(mux.Vars(request)["cueNumber"]); err != nil {
			sendError(writer, err)
		} else {
			if err := project.Instance.JumpTo(cueNumber); err != nil {
				sendError(writer, err)
			} else {
				sendOK(writer)
			}
		}
	}
}

func stopPlaying(writer http.ResponseWriter, request *http.Request) {
	if checkProject(writer) {
		project.Instance.StopPlaying()
		sendOK(writer)
	}
}

func updateProjectName(writer http.ResponseWriter, request *http.Request) {
	if checkProject(writer) {
		project.Instance.SetName(mux.Vars(request)["name"])
		sendOK(writer)
	}
}

func updateProjectSettings(writer http.ResponseWriter, request *http.Request) {
	if checkProject(writer) {
		bufferSize, err := strconv.Atoi(request.FormValue("BufferSize"))
		if err != nil {
			sendError(writer, err)
			return
		}
		if bufferSize < 0 {
			sendError(writer, fmt.Errorf("buffer size cannot be less than 0: %d", bufferSize))
		}
		project.Instance.SetSettings(project.Settings{
			BufferSize: uint(bufferSize),
		})
		sendOK(writer)
	}
}

func loadProject(writer http.ResponseWriter, request *http.Request) {
	if checkProject(writer) {
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
		project.Instance.Close()
		project.Instance = newProject
		sendOK(writer)
	}
}

func saveProject(writer http.ResponseWriter, request *http.Request) {
	if checkProject(writer) {
		writer.Header().Set("Content-Disposition", "attachment; filename="+project.Instance.Name+".audioq")
		writer.Header().Set("Content-Type", request.Header.Get("Content-Type"))
		if err := json.NewEncoder(writer).Encode(*project.Instance); err != nil {
			sendError(writer, err)
			return
		}
	}
}
