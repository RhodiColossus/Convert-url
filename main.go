package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

//URLsTb = таблица с данными

type URLsTb struct {
	lond  string
	short string
	date  string
}

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
var db *sql.DB

//
func oldURLchecker() {
	dt := time.Now()
	date := dt.Format("01-02-2006")

	db.Exec("DELETE FROM urlstb WHERE date < $1", date)
}

//setURL = запись длинного url в базу и проверка жизни коротких url
func setURL(longURL LongU, shortURL ShortU) string {
	oldURLchecker()
	dt := time.Now()
	date := dt.Format("01-02-2006")
	_, err := db.Exec("INSERT INTO urlstb VALUES($1, $2, $3)", longURL.HUrl, shortURL.HUrl, date)

	if err != nil {
		result := db.QueryRow("SELECT short FROM urlstb WHERE long = $1", longURL.HUrl)

		err := result.Scan(&shortURL.HUrl)
		if err != nil {
			//ошибка
		}
		return shortURL.HUrl
	}

	return shortURL.HUrl
}

//GetMD5Hash = делает мд5 хеш
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

//Принимает Короткий url ,возвращает Длинный
func postShortURL(w http.ResponseWriter, request *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	l := LongU{}
	s := ShortU{}
	_ = json.NewDecoder(request.Body).Decode(&l)

	validURL := govalidator.IsURL(l.HUrl)
	if validURL != true {
		s.HUrl = "incorrect url" //вывод ошибки о некоректности URL
	} else {
		resporse := URLgenerate()
		s.HUrl = resporse
		s.HUrl = setURL(l, s)
	}
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

//Принимает Длинный url ,возвращает Короткий
func postLongURL(w http.ResponseWriter, request *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	s := ShortU{}
	_ = json.NewDecoder(request.Body).Decode(&s)

	validURL := govalidator.IsURL(s.HUrl)
	if validURL != true {
		s.HUrl = "incorrect url" //вывод ошибки о некоректности URL
	} else {
		result := db.QueryRow("SELECT long FROM urlstb WHERE short = $1", s.HUrl)
		l := LongU{}
		err := result.Scan(&l.HUrl)
		if err != nil {
			//ошибка
		}
		json.NewEncoder(w).Encode(l)
	}
}

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://rodp:qwe@localhost/urls")
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/short", postShortURL).Methods("POST")
	router.HandleFunc("/long", postLongURL).Methods("POST")

	log.Fatal(http.ListenAndServe(":8004", router))
}
