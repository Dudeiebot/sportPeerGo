package queries

import (
	"log"

	"github.com/dudeiebot/sportPeerGo/pkg/dbs"
	"github.com/dudeiebot/sportPeerGo/pkg/user/model"
)

func RegisterQuery(u model.User, d *dbs.Service) int {
	queri := `INSERT INTO users (username, email, phone, password, verificationToken, bio)VALUES (?, ?, ?, ?, ?, ?)`

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
