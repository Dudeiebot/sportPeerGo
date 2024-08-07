package queries

import (
	"database/sql"
	"log"

	"github.com/dudeiebot/sportPeerGo/pkg/dbs"
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

func VerifyEmailQueries(d *dbs.Service, token string) (sql.Result, error) {
	queri := `UPDATE users SET is_verified = TRUE, verification_token = NULL WHERE verification_token = ?`

	res, err := d.DB.Exec(queri, token)
	return res, err
}
