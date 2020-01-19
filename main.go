package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

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

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

//Принимает длинный url ,возвращает короткий
func postShortURL(w http.ResponseWriter, request *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	l := LongU{}
	s := ShortU{}
	_ = json.NewDecoder(request.Body).Decode(&l)

	validURL := govalidator.IsURL(l.HUrl)
	if validURL != true {
		l.HUrl = "incorrect url" //вывод ошибки о некоректности URL
	}

	resporse := URLgenerate()
	s.HUrl = resporse

	json.NewEncoder(w).Encode(s)
}

//URLgenerate = Генерация короткого URL
func URLgenerate() string {
	shortURL := "http://qwe.com/"

	rand.Seed(time.Now().UTC().UnixNano())
	var bytes int
	bytes = rand.Intn(100000000)

	var s string
	s = strconv.Itoa(bytes)

	md5String := GetMD5Hash(s)
	result := md5String[0:5]

	shortURL += result

	return shortURL
}

//Принимает Короткий url ,возвращает Длинный
func postLongURL(w http.ResponseWriter, request *http.Request) {

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

	router.HandleFunc("/short", postShortURL).Methods("POST")
	router.HandleFunc("/long", postShortURL).Methods("POST")

	log.Fatal(http.ListenAndServe(":8001", router))
}
