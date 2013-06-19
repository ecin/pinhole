package main

import (
  "os/exec"
  "net/http"
  "log"
  "os"
  "path"
//  "fmt"
)

const (
  SCREENSHOT_PATH = "/tmp"
  TWITTER_CARD_URL = "http://localhost:3000/twitter_cards/"
)

func screenshotPath(opaqueId string) string {
  return path.Join(SCREENSHOT_PATH, opaqueId + ".png")
}

func twitterCardHandler(response http.ResponseWriter, request *http.Request) {
  opaqueId := request.URL.Path[len("/cards/"):]
  outputFile := screenshotPath(opaqueId)
  url := TWITTER_CARD_URL + opaqueId

  cached := true

  if _, err := os.Stat(screenshotPath(opaqueId)); os.IsNotExist(err) {
    cached = false
  }

  // redirect to file if everything is fine
  if !cached {
    go createThumbnail(url, outputFile)
    response.WriteHeader(http.StatusAccepted)
  } else {
    response.Header().Set("Location", "/static/" + opaqueId + ".png")
    response.WriteHeader(http.StatusFound)
  }
}

func createThumbnail(url string, outputFile string) {
  phantomJS := exec.Command("bin/screenshot", url, outputFile)
  phantomJS.Run()
}

func Log(handler http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
    handler.ServeHTTP(w, r)
  })
}

func main() {
  http.HandleFunc("/cards/", twitterCardHandler)
  http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir(SCREENSHOT_PATH))))
  http.ListenAndServe(":8080", Log(http.DefaultServeMux))
}
