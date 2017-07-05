package main

import (
	"agaza/api"
	"agaza/config"
	"agaza/logger"
	"os"
)


func main() {
	logger.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	configuration := config.LoadConfiguration()
	logger.Trace.Println("Server Started on http://localhost:", configuration.APIserverport)
	api.GetFastHTTPServer().StartAndServeAPIs()
}
