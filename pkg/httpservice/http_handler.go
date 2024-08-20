package httpservice

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
)

func NewHandler[IN, OUT any](
	targetFunc func(context.Context, IN) (OUT, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in IN
		if reflect.TypeOf(in) != reflect.TypeOf((*http.Request)(nil)) {
			if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		} else {
			in = any(r).(IN)
		}

		out, err := targetFunc(r.Context(), in)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		switch v := any(out).(type) {
		case *LoginResponse:
			http.SetCookie(w, &http.Cookie{
				Name:     "token",
				Value:    v.Token,
				HttpOnly: true,
				Secure:   true,
				MaxAge:   3600,
				SameSite: http.SameSiteStrictMode,
			})
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(out)
	}
}
