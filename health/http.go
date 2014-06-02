package health

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ninjablocks/mqtt-proxy/conf"
)

type HeathServer struct {
}

func StartHealthServer(conf *conf.Configuration) {
	http.HandleFunc("/health", HomeHandler)
	log.Printf("[health] listening %s", ":1880")
	log.Fatal(http.ListenAndServe(":1880", nil))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}
