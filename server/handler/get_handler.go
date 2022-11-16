package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/nubesk/binn"
	"github.com/nubesk/binn-server/spool"
)

const (
	binnUUIDKey = "Binn-UUID"
)

func GetHandlerFunc(bn *binn.Binn, logger *log.Logger) http.HandlerFunc {
	spl := spool.New()
	return func(w http.ResponseWriter, r *http.Request) {
		initFunc := func() {
			u := uuid.New()
			us := u.String()
			cookie := &http.Cookie{
				Name:  "Binn-UUID",
				Value: us,
			}
			http.SetCookie(w, cookie)
			spl.Reset(us)
			err := bn.Subscribe(func(b *binn.Bottle) bool {
				if err := spl.Publish(us, b); err != nil {
					return false
				}
				return true
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write([]byte("[]"))
			return
		}
		cookie, err := r.Cookie(binnUUIDKey)
		if err != nil {
			initFunc()
			return
		}
		us := cookie.Value
		bs, err := spl.Get(us)
		if err != nil {
			initFunc()
			return
		}
		spl.Reset(us)
		res := bottlesToResponse(bs)
		if bytes, err := json.Marshal(res); err == nil {
			if _, err := w.Write(bytes); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
