package main

import (
	"flag"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/zantabri/ss-service/handlers"
	"github.com/zantabri/ss-service/store"
)

func main() {

	var sd = flag.String("sd", "", "secret directory required")
	flag.Parse()

	if len(*sd) == 0 {
		panic("sd : secret directory is required")
	}

	store, err := store.NewFileStore(sd)

	if err != nil {
		panic(err.Error())
	}

	handlers := handlers.New(&store)
	router := httprouter.New()

	if err != nil {
		panic(err.Error())
	}

	router.GET("/health", handlers.HealthCheck)
	router.POST("/", handlers.AddSecret)
	router.GET("/", handlers.GetSecret)
	http.ListenAndServe(":8080", router)

}
