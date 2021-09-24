package handlers

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

//WriteHTTPError - default write error func
func WriteHTTPError(w http.ResponseWriter, httpStatusCode int, err error) {
	w.WriteHeader(httpStatusCode)
	_, errWrite := w.Write([]byte("{\"Message\":\"" + err.Error() + "\"}"))
	if errWrite != nil {
		log.Printf("[ERROR] Not Writing to ResponseWriter error %s due: %s", err.Error(), errWrite.Error())
	}
}
