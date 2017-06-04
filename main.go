package main

import (
	"net/http"
	"html/template"
	"io/ioutil"
	"encoding/json"
	"regexp"
	"strconv"
	"errors"
)

type Gif struct {
	Filename string
}

type JsonResponse struct {
	GifList []Gif
	IsLastPage bool
}

var tpl *template.Template
var gifList []Gif
var itemsPerPage int = 6

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))

	jsonData, err := ioutil.ReadFile("./data/gif.json")
	if err != nil {
		panic("Failed to read JSON data.")
	}
	err = json.Unmarshal(jsonData, &gifList)
	if err != nil {
		panic("Failed to unmarshal JSON data.")
	}
}

func main() {

	http.HandleFunc("/", index)

	http.HandleFunc("/api/gif", list)

	// static contents
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("./public"))))

	http.ListenAndServe(":9876", nil)
}

func index(w http.ResponseWriter, req *http.Request)  {
	tpl.ExecuteTemplate(w, "index.html", gifList[:itemsPerPage])
}

func list(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		code := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(code), code)
		return
	}

	page := getPageParam(req)

	gif, isLastPage, err := getPageGif(page)
	if err != nil {
		http.NotFound(w, req)
		return
	}

	jsonText, err := json.Marshal(JsonResponse{
		GifList: gif,
		IsLastPage: isLastPage,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonText)
}

func getPageParam(req *http.Request) int {
	pageParam := req.URL.Query()["page"]

	if pageParam == nil {
		return 1
	}

	page := pageParam[0]

	if page == "" {
		return 1
	}

	// 正規表現で半角数字にマッチさせる
	pattern := regexp.MustCompile(`\d`)
	if !pattern.MatchString(page) {
		return 1
	}

	// 数値に変換
	res, _ := strconv.Atoi(pageParam[0])

	if res == 0 {
		return 1
	}

	return res
}

func getPageGif(page int) (gif []Gif, isLastPage bool, err error) {
	startIndex := itemsPerPage * (page - 1)
	lastIndex := startIndex + itemsPerPage

	actualLastIndex := len(gifList) - 1

	if startIndex > actualLastIndex {
		return nil, false, errors.New("Page not found.")
	} else if lastIndex >= actualLastIndex {
		return gifList[startIndex:], true, nil
	} else {
		return gifList[startIndex:lastIndex], false, nil
	}
}
