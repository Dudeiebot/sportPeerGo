package queries

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/dudeiebot/sportPeerGo/pkg/adapter/dbs"
	"github.com/dudeiebot/sportPeerGo/pkg/user/model"
)

func RegisterQuery(u model.User, d *dbs.Service) int {
	queri := `INSERT INTO users (username, email, phone, password, verification_token, bio)VALUES (?, ?, ?, ?, ?, ?)`

	res, err := d.DB.Exec(queri, u.Username,
		u.Email,
		u.Phone,
		u.Password,
		u.VerificationToken,
		u.Bio)
	if err != nil {
		log.Printf("error inserting into db: %s\n", err)
	}

	id, _ := res.LastInsertId()
	u.ID = int(id)
	return u.ID
}

func GetHashedAuth(ctx context.Context, c model.Credentials, d *dbs.Service) (*model.User, error) {
	var user model.User
	queri := `SELECT id, email, phone, password, is_verified FROM users WHERE email = ? OR phone = ? LIMIT 1`
	err := d.DB.QueryRowContext(ctx, queri, c.Access, c.Access).
		Scan(&user.ID, &user.Email, &user.Phone, &user.Password, &user.IsVerified)
	if err != nil {
		if err == sql.ErrNoRows {
			// have your custom error here
			return nil, errors.New("User Not Found")
		}
		return nil, err
	}
	return &user, nil
}

func VerifyEmailQueries(ctx context.Context, d *dbs.Service, token string) (sql.Result, error) {
	queri := `UPDATE users SET is_verified = TRUE, verification_token = NULL WHERE verification_token = ?`

	res, err := d.DB.ExecContext(ctx, queri, token)
	return res, err
}
