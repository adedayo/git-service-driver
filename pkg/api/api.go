package api

import (
	"fmt"
	"log"
	"net/http"

	model "github.com/adedayo/git-service-driver/pkg"
	"github.com/gorilla/handlers"
)

var (
	configManager model.ConfigManager
)

//ServeAPI serves the analysis service on the specified port
func ServeAPI(config Config) {
	hostPort := "localhost:%d"
	if !config.Local {
		// not localhost electron app
		hostPort = ":%d"
	}

	corsOptions = append(corsOptions, handlers.AllowedOrigins(allowedOrigins))
	log.Printf("Running Git Service API on %s", fmt.Sprintf(hostPort, config.ApiPort))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(hostPort, config.ApiPort), handlers.CORS(corsOptions...)(getRoutes((config.CodeBaseDir)))))
}
