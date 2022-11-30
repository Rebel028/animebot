package downloadhelper

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
)

func DownloadImageAsBase64(url string) string {
	bytes := DownloadFile(url)
	return toBase64(bytes)
}

func DownloadFile(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
