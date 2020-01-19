package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
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
	id    int
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
func setURL(url LongU) {

	long := "i_am+long"
	short := "i_am_short"
	id := 1
	date := "10.10.2020"
	result, err := db.Exec("INSERT INTO urlstb VALUES($1, $2, $3, $4)", id, long, short, date)
	if err != nil {
		//обработка ошибки
		fmt.Println(result)
	}
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
		l.HUrl = "incorrect url" //вывод ошибки о некоректности URL
	}

	resporse := URLgenerate()
	s.HUrl = resporse
	setURL(l)
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

	log.Fatal(http.ListenAndServe(":8002", router))
}
