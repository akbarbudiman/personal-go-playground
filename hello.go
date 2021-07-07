package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Document struct {
	Id         int       `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	LastUpdate time.Time `json:"lastUpdate"`
}

var DOCUMENTS []Document

func main() {
	RunRedisExample()
	DOCUMENTS = []Document{
		{Id: 1, Title: "Daftar Belanja", Content: "Telur, Roti tawar, Laundry", LastUpdate: time.Now()},
		{Id: 2, Title: "Agenda", Content: "Cari Kontrakan", LastUpdate: time.Now()},
	}
	handleRequests()
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/documents", getAllDocument).Methods("GET")
	myRouter.HandleFunc("/documents", addDocument).Methods("POST")
	myRouter.HandleFunc("/documents/{id}", getCertainDocument).Methods("GET")
	myRouter.HandleFunc("/documents/{id}", removeCertainDocument).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method + " " + r.RequestURI)
	fmt.Fprintf(w, "Welcome to the HomePage!")
}

func getAllDocument(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method + " " + r.RequestURI)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(DOCUMENTS)
}

func addDocument(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method + " " + r.RequestURI)

	reqBody, _ := ioutil.ReadAll(r.Body)
	var newDocument Document
	json.Unmarshal(reqBody, &newDocument)
	newDocument.LastUpdate = time.Now()

	for _, data := range DOCUMENTS {
		if data.Id == newDocument.Id {
			w.WriteHeader(http.StatusConflict)
			return
		}
	}

	DOCUMENTS = append(DOCUMENTS, newDocument)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(DOCUMENTS)
}

func getCertainDocument(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method + " " + r.RequestURI)
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	idStrToSearch := vars["id"]
	idIntToSearch, _ := strconv.ParseInt(idStrToSearch, 10, 0)

	for _, data := range DOCUMENTS {
		if data.Id == int(idIntToSearch) {
			json.NewEncoder(w).Encode(data)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func removeCertainDocument(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method + " " + r.RequestURI)

	vars := mux.Vars(r)
	idStrToSearch := vars["id"]
	idIntToSearch, _ := strconv.ParseInt(idStrToSearch, 10, 0)

	foundIndex := -1
	for index, data := range DOCUMENTS {
		if data.Id == int(idIntToSearch) {
			foundIndex = index
			break
		}
	}

	if foundIndex == -1 {
		w.WriteHeader(http.StatusNotFound)
	} else {
		DOCUMENTS = append(DOCUMENTS[:foundIndex], DOCUMENTS[foundIndex+1:]...)
		w.WriteHeader(http.StatusNoContent)
	}

}
