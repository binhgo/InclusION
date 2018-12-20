package util

import (
	"crypto/sha256"
	"fmt"
	"net/http"

)

func Hash(input string) string {
	h := sha256.New()

	h.Write([]byte(input))
	//fmt.Printf("%x", h.Sum(nil))
	hashString := fmt.Sprintf("%x", h.Sum(nil))
	return hashString
}


func CheckBodyNil(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
}
