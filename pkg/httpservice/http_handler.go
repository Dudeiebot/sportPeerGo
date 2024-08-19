package httpservice

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
)

func NewHandler[IN, OUT any](
	s *Server,
	targetFunc func(context.Context, *Server, IN) (OUT, error),
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

		out, err := targetFunc(r.Context(), s, in)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(out)
	}
}
