package main

import (
	"LaserTaskSystem/internal/api/api_handler_v1"
	"LaserTaskSystem/internal/api/api_v1"
	"LaserTaskSystem/internal/config"
	"LaserTaskSystem/internal/task"
	"embed"
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

func main() {
	// Get config
	conf := config.NewConfig()
	err := conf.Load("./config.yaml", true)
	if err != nil {
		log.Fatalln(err)
	}

	// Init task managers
	managers := make(map[string]*task.Manager)
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

	// Create API v1 handler
	handler := api_handler_v1.NewHandler(managers)

	// Init API v1 server with handler
	ogenServer, err := api_v1.NewServer(handler)
	if err != nil {
		log.Fatalln(err)
	}

	// Create http multiplexer
	mux := http.NewServeMux()

	// Register ogen server
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", ogenServer))

	// Add file server for static
	subFS, err := getStatic()
	if err != nil {
		log.Fatal(err)
	}
	staticHandler := http.FileServer(http.FS(subFS))
	mux.Handle("/", http.StripPrefix("/", staticHandler))

	// Print server port
	log.Println("Listening on port " + strconv.Itoa(conf.Server.Port) + "...")

	// Start server
	err = http.ListenAndServe(":"+strconv.Itoa(conf.Server.Port), mux)
	if err != nil {
		log.Fatal(err)
	}

}
