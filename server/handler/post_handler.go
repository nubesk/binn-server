package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/nubesk/binn"
)

func PostHandlerFunc(bn *binn.Binn, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var body request
		if err := json.Unmarshal(bytes, &body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		b := requestToBottle(&body)
		err = bn.Publish(b)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		if bytes, err := json.Marshal(body); err == nil {
			if _, err := w.Write(bytes); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if logger != nil {
				logger.Printf("[in] {\"message\": \"%s\"}\n", b.Msg)
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
