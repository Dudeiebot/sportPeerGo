package httpservice

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/dudeiebot/sportPeerGo/pkg/adapter/dbs"
	"github.com/dudeiebot/sportPeerGo/pkg/user/model"
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
		case *LogoutResponse:
			http.SetCookie(w, &http.Cookie{
				Name:     "token",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			})
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(out)
	}
}

func NewUpdateHandler(
	s *Server,
	queryFunc func(ctx context.Context, dbs *dbs.Service, user model.User) (sql.Result, error),
	successMessage string,
) http.HandlerFunc {
	return NewHandler(func(ctx context.Context, req *http.Request) (*Response, error) {
		id := chi.URLParam(req, "id")
		userId := ctx.Value("userId").(int64)
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return &Response{Message: "invalid user ID format"}, nil
		}
		if idInt != userId {
			return &Response{Message: "Unauthorized"}, nil
		}
		var user model.User
		if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
			return &Response{Message: "invalid request body"}, nil
		}
		user.ID = int(userId)
		res, err := queryFunc(ctx, s.DBS, user)
		if err != nil {
			return &Response{Message: "Error executing DB query"}, nil
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return &Response{Message: "Error getting rows affected"}, err
		}
		if rowsAffected == 0 {
			return &Response{Message: "No rows affected"}, nil
		}
		return &Response{Message: successMessage}, nil
	})
}
