package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
)

//LongU = long Url
type LongU struct {
	HUrl string `json:"url"`
}

//ShortU = short url
type ShortU struct {
	HUrl string `json:"url"`
}

var longUrls []LongU
var shortUrls []ShortU

//Принимает длинный url ,возвращает короткий
func postShortURL(w http.ResponseWriter, request *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	l := LongU{}
	_ = json.NewDecoder(request.Body).Decode(&l)

	validURL := govalidator.IsURL(l.HUrl)
	if validURL != true {
		//обработка неверного URL
	}

	longUrls = append(longUrls, l)
}

func getShortURL(w http.ResponseWriter, request *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	l := LongU{}
	_ = json.NewDecoder(request.Body).Decode(&l)
	longUrls = append(longUrls, l)
}

func main() {
	router := mux.NewRouter()

	shortUrls = append(shortUrls, ShortU{HUrl: "short-1"})
	router.HandleFunc("/short", getShortURL).Methods("GET")
	router.HandleFunc("/short", postShortURL).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", router))
}
