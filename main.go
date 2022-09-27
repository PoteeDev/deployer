package main

import(
	"fmt"
	"net/http"
	"log"
)

func sayhello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Привет!")
}


func main() {
	http.HandleFunc("/", sayhello)
	err := http.ListenAndServe(":8080", nil)
	
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

