package main

import (
	"flag"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/zantabri/ss-service/handlers"
)


func main() {

	var sd = flag.String("sd", "", "secret directory required")
	flag.Parse()
	
	if len(*sd) == 0 {
		panic("sd : secret directory is required")
	}

	router := httprouter.New()
	handlers, err := handlers.New(*sd)

	if err != nil {
		panic(err.Error())
	}
	


	router.GET("/health", handlers.HealthCheck)
	router.POST("/", handlers.AddSecret)
	router.GET("/", handlers.GetSecret)
	http.ListenAndServe(":8080", router)

}
