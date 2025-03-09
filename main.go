// main.go
package main

import (
	"github.com/gin-contrib/cors"
)

func main() {
	dependencies := NewDependencies()


	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowCredentials = true

	dependencies.engine.Use(cors.New(corsConfig))
	dependencies.Run()
	//
	
}

