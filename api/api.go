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

func getProject(writer http.ResponseWriter, request *http.Request) {
	if true {
		json.NewEncoder(writer).Encode(project.GetInfo())
	}
}

func addCue(writer http.ResponseWriter, request *http.Request) {
	if true {
		if err := request.ParseMultipartForm(32 << 20); err != nil {
			sendError(writer, err)
			return
		}
		file, header, err := request.FormFile("audioFile")
		if err != nil {
			sendError(writer, err)
			return
		}
		if err := project.AddCue(request.FormValue("cueName"), header.Filename, file); err != nil {
			sendError(writer, err)
			return
		}
		sendOK(writer)
	}
}

func removeCue(writer http.ResponseWriter, request *http.Request) {
	if true {
		cueNumber, err := strconv.Atoi(mux.Vars(request)["cueNumber"])
		if err != nil {
			sendError(writer, err)
			return
		}
		if err := project.RemoveCue(cueNumber); err != nil {
			sendError(writer, err)
		} else {
			sendOK(writer)
		}
	}
}

func renameCue(writer http.ResponseWriter, request *http.Request) {
	if true {
		vars := mux.Vars(request)
		cueNumber, err := strconv.Atoi(vars["cueNumber"])
		if err != nil {
			sendError(writer, err)
			return
		}
		if err := project.RenameCue(cueNumber, vars["cueName"]); err != nil {
			sendError(writer, err)
		} else {
			sendOK(writer)
		}
	}
}

func moveCue(writer http.ResponseWriter, request *http.Request) {
	if true {
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
		if err := project.MoveCue(from, to); err != nil {
			sendError(writer, err)
		} else {
			sendOK(writer)
		}
	}
}

func playNext(writer http.ResponseWriter, request *http.Request) {
	if true {
		if err := project.PlayNext(); err != nil {
			sendError(writer, err)
		} else {
			sendOK(writer)
		}
	}
}

func jumpTo(writer http.ResponseWriter, request *http.Request) {
	if true {
		if cueNumber, err := strconv.Atoi(mux.Vars(request)["cueNumber"]); err != nil {
			sendError(writer, err)
		} else {
			if err := project.JumpTo(cueNumber); err != nil {
				sendError(writer, err)
			} else {
				sendOK(writer)
			}
		}
	}
}

func stopPlaying(writer http.ResponseWriter, request *http.Request) {
	if true {
		project.StopPlaying()
		sendOK(writer)
	}
}

func updateProjectName(writer http.ResponseWriter, request *http.Request) {
	if true {
		project.SetName(mux.Vars(request)["name"])
		sendOK(writer)
	}
}

func updateProjectSettings(writer http.ResponseWriter, request *http.Request) {
	if true {
		bufferSize, err := strconv.Atoi(request.FormValue("BufferSize"))
		if err != nil {
			sendError(writer, err)
			return
		}
		if bufferSize < 0 {
			sendError(writer, fmt.Errorf("buffer size cannot be less than 0: %d", bufferSize))
		}
		project.SetSettings(project.Settings{
			BufferSize: uint(bufferSize),
		})
		sendOK(writer)
	}
}

func loadProject(writer http.ResponseWriter, request *http.Request) {
	if true {
		if err := request.ParseMultipartForm(32 << 20); err != nil {
			sendError(writer, err)
			return
		}
		file, _, err := request.FormFile("audioqProject")
		if err != nil {
			sendError(writer, err)
			return
		}
		if err := project.LoadProject(file); err != nil {
			sendError(writer, err)
			return
		}
		sendOK(writer)
	}
}

func saveProject(writer http.ResponseWriter, request *http.Request) {
	if true {
		writer.Header().Set("Content-Disposition", "attachment; filename="+project.GetName()+".audioq")
		writer.Header().Set("Content-Type", request.Header.Get("Content-Type"))
		if err := project.SaveProject(writer); err != nil {
			sendError(writer, err)
			return
		}
	}
}
