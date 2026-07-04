package main

import (
	"log"
	"net/http"
	"os"
	"golang.org/x/net/webdav"
)

func main() {
	// /data 디렉토리를 WebDAV 루트로 지정 (도커 볼륨 매핑 경로)
	fs := webdav.Dir("/data")

	// 🌟 핵심: 파일 잠금 장치를 메모리(In-Memory)에만 저장하여 .DAV 폴더 생성을 원천 차단
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

	// .env 파일에서 주입될 인증 정보 로드
	username := os.Getenv("WEBDAV_USER")
	password := os.Getenv("WEBDAV_PASSWORD")

	if username == "" || password == "" {
		log.Fatal("🚨 환경변수 WEBDAV_USER 또는 WEBDAV_PASSWORD가 설정되지 않았습니다.")
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

	log.Println("🚀 지저분한 찌꺼기 없는 Go WebDAV Server가 80포트에서 가동됩니다.")
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}