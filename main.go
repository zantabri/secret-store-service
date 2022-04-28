package main

import (

	"net/http"

	"github.com/julienschmidt/httprouter"
)


func main() {

	router := httprouter.New()
	
	router.GET("/health", HealthCheck)
	router.POST("/", AddSecret)
	router.GET("/", GetSecret)
	http.ListenAndServe(":8080", router)

}
