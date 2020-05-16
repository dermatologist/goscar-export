package main

import (
	"github.com/intervention-engine/fhir/server"
)

func main() {

	// set up the commandline flags (-mongo and -pgurl)
	// reqLog := flag.Bool("reqlog", false, "Enables request logging -- do NOT use in production")
	// serverURL := flag.String("server", "", "The full URL for the root of the server")
	// dbName := flag.String("dbname", "fhir", "Mongo database name")
	// idxConfigPath := flag.String("idxconfig", "config/indexes.conf", "Path to the indexes config file")
	// mongoHost := flag.String("mongohost", "localhost", "the hostname of the mongo database")
	// readOnly := flag.Bool("readonly", false, "Run the API in read-only mode (no creates, updates, or deletes allowed)")

	// flag.Parse()

	// If using meteor, then meteor port + 1
	// 	mongoHost := "mongodb://127.0.0.1:8086/fhir"
	mongoHost := "mongodb://127.0.0.1:27017/fhir"
	s := server.NewServer(mongoHost)

	config := server.DefaultConfig
	config.ServerURL = "http://localhost:7001/"

	// if *serverURL != "" {
	// 	config.ServerURL = *serverURL
	// }

	// if *dbName != "" {
	// 	config.DatabaseName = *dbName
	// }

	// if *idxConfigPath != "" {
	// 	config.IndexConfigPath = *idxConfigPath
	// }

	// if *reqLog {
	// 	s.Engine.Use(server.RequestLoggerHandler)
	// }

	// if *readOnly {
	// 	s.Engine.Use(ReadOnlyHandler)
	// }

	s.Run(config)

}

// ReadOnlyHandler makes the API read-only and responds to any requests that are not
// GET, HEAD, or OPTIONS with a 403 Forbidden error.
// func ReadOnlyHandler(c *gin.Context) {

// 	method := c.Request.Method
// 	switch method {
// 	// allowed methods:
// 	case "GET", "HEAD", "OPTIONS":
// 		return
// 	// all other methods:
// 	default:
// 		c.AbortWithStatus(403)
// 	}
// }
