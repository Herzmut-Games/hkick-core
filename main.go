package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"
)

var (
	addr = flag.String("addr", ":3000", "http service address")
)

func main() {
	uri, err := url.Parse("mqtt://172.30.1.32:1883")
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/goals", goalCount)

	client = connect("pub", uri)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
