package handler

import (
	"crypto/rand"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"html/template"
	"math/big"
	"net/http"
	"strings"
	"urlShortener/cache"
)

const baseUrl = "localhost:8080/"

type shortenUrlReqBody struct {
	Url string `json:"url"`
}

type shortenUrlRespBody struct {
	Url string `json:"url"`
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, nil)
}

func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody shortenUrlReqBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userUrl := ensureHTTPS(reqBody.Url)
	newPath := generateUrlPath()
	cache.Cache[newPath] = userUrl
	cache.WriteToDb(newPath, userUrl)

	shortenedURL := baseUrl + newPath
	respBody := shortenUrlRespBody{shortenedURL}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(respBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	url := chi.URLParam(r, "url")
	redirectUrl, ok := cache.Cache[url]
	if !ok {
		http.NotFound(w, r)
	}
	http.Redirect(w, r, redirectUrl, http.StatusFound)
}

func ensureHTTPS(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "https://" + url
	}
	return url
}

func generateUrlPath() string {
	i := 1
	for {
		path, _ := generateRandomString(i)
		if _, ok := cache.Cache[path]; !ok {
			return path
		}
		i += 1
	}
}

func generateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[randomIndex.Int64()]
	}

	return string(result), nil
}
