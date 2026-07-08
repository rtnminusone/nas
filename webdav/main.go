package main

import (
  "log"
  "net/http"
  "os"
  "golang.org/x/net/webdav"
)

func main() {
  // set data directory
  fs := webdav.Dir("/srv")

  // block .DAV
  ls := webdav.NewMemLS()

  handler := &webdav.Handler{
    FileSystem: fs,
    LockSystem: ls,
    Logger: func(r *http.Request, err error) {
      if err != nil {
        log.Printf("WEBDAV ERR: %s [%s]", r.URL.Path, err)
      }
    },
  }

  username := os.Getenv("WEBDAV_USER")
  password := os.Getenv("WEBDAV_PASSWORD")

  if username == "" || password == "" {
    log.Fatal("user setting fail")
  }

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    user, pass, ok := r.BasicAuth()
    if !ok || user != username || pass != password {
      w.Header().Set("WWW-Authenticate", `Basic realm="NAS WebDAV"`)
      w.WriteHeader(http.StatusUnauthorized)
      w.Write([]byte("401 Unauthorized\n"))
      return
    }
    handler.ServeHTTP(w, r)
  })

  log.Println("run Go WebDAV Server at port 80")
  if err := http.ListenAndServe(":80", nil); err != nil {
    log.Fatal(err)
  }
}