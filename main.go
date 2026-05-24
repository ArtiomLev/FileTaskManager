package main

import (
	"LaserTaskSystem/config"
	"LaserTaskSystem/task"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"strconv"
)

//go:embed web/static/*
var staticFiles embed.FS

func getStatic() (fs.FS, error) {
	return fs.Sub(staticFiles, "web/static")
}

var conf = config.NewConfig()
var managers = make(map[string]*task.Manager)

var tasks []task.Task

func main() {
	// Get config
	err := conf.Load("./config.yaml", true)
	if err != nil {
		log.Fatalln(err)
	}

	// Init task managers
	for _, managerConfig := range conf.TaskManagers {
		manager, err := task.NewManager(
			managerConfig.Name,
			managerConfig.DisplayName,
			managerConfig.ActivePath,
			managerConfig.CompletedPath)
		if err != nil {
			log.Fatalln(err)
		}
		managers[managerConfig.Name] = manager
	}

	http.HandleFunc("/task", tasksHandler)
	http.HandleFunc("/health", healthCheckHandler)

	subFS, err := getStatic()
	if err != nil {
		log.Fatal(err)
	}
	staticHandler := http.FileServer(http.FS(subFS))
	http.Handle("/", http.StripPrefix("/", staticHandler))

	err = http.ListenAndServe(":"+strconv.Itoa(conf.Server.Port), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Listening on port 8080")

	fmt.Println("")

}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTask(w, r)
	case http.MethodPost:
		addTask(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func getTask(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(tasks)
	if err != nil {
		return
	}
}

func addTask(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Не удалось прочитать форму: "+err.Error(), http.StatusBadRequest)
		return
	}

	contract := r.FormValue("ContractNumber")
	name := r.FormValue("Name")
	if contract == "" || name == "" {
		http.Error(w, "Номер договора и имя задания обязательны", http.StatusBadRequest)
		return
	}

	filesHeader := r.MultipartForm.File["files"]
	if len(filesHeader) == 0 {
		http.Error(w, "Не выбрано ни одного файла", http.StatusBadRequest)
		return
	}

	var fileClosers []io.Closer
	defer func() {
		for _, closer := range fileClosers {
			err := closer.Close()
			if err != nil {
				continue
			}
		}
	}()

	var inputs []task.FileInput
	for _, fh := range filesHeader {
		src, err := fh.Open()
		if err != nil {
			http.Error(w, "Ошибка открытия файла: "+err.Error(), http.StatusInternalServerError)
			return
		}

		fileClosers = append(fileClosers, src)
		inputs = append(inputs, task.FileInput{
			Name:   fh.Filename,
			Reader: src,
		})
	}

	createdTask, err := managers["tube"].CreateTask(contract, name, inputs)
	if err != nil {
		switch {
		case errors.Is(err, task.ErrTaskExists):
			http.Error(w, "Задание уже существует", http.StatusConflict)
		case errors.Is(err, task.ErrInvalidContract), errors.Is(err, task.ErrInvalidName):
			http.Error(w, "Неверный номер договора или имя задания", http.StatusBadRequest)
		case errors.Is(err, task.ErrNoFiles):
			http.Error(w, "Не выбрано ни одного файла", http.StatusBadRequest)
		case errors.Is(err, task.ErrCannotMoveTempToTask):
			http.Error(w, "Ошибка сохранения задания", http.StatusInternalServerError)
		default:
			http.Error(w, "Ошибка создания задания: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = fmt.Fprintf(w, "Задание \"%s/%s\" успешно создано с %d файлами",
		createdTask.ContractNumber, createdTask.Name, len(inputs))
	if err != nil {
		return
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Service is working!")
	if err != nil {
		return
	}
}
