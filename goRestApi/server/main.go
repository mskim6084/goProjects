package main

import (
	"encoding/json"
	"errors"
	"fmt"
	albumData "goRestApi/server/albumData"
	"net/http"
)

func getAllAlbums(response http.ResponseWriter, request *http.Request) {
	fmt.Printf("got /getAllAlbums request\n")

	response.Header().Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(albumData.Albums)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = response.Write(jsonData)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
}

func addNewAlbum(response http.ResponseWriter, request *http.Request) {

}

func main() {
	http.HandleFunc("/getAllAlbums", getAllAlbums)
	http.HandleFunc("/addNewAlbum", addNewAlbum)

	err := http.ListenAndServe(":3333", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server closed \n")
	} else if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
