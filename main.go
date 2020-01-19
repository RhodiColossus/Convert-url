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

//setURL = запись длинного url в базу и проверка жизни коротких url
func setURL(longURL LongU, shortURL ShortU) string {
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

//Принимает длинный url ,возвращает короткий
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

	router.HandleFunc("/long", postShortURL).Methods("POST")
	router.HandleFunc("/short", postLongURL).Methods("POST")

	log.Fatal(http.ListenAndServe(":8004", router))
}
