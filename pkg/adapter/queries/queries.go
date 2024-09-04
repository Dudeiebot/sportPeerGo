package queries

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dudeiebot/sportPeerGo/pkg/adapter/dbs"
	"github.com/dudeiebot/sportPeerGo/pkg/user/model"
)

func RegisterQuery(u model.User, d *dbs.Service) int {
	// add ctx
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

func GetOtpQueries(
	ctx context.Context,
	d *dbs.Service,
	otp, email string,
) (*model.ForgetPass, error) {
	queri := `
		SELECT otp, otp_expiration
		FROM users
		WHERE email = ?
	`
	var forgetPass model.ForgetPass
	err := d.DB.QueryRowContext(ctx, queri, email).Scan(&forgetPass.Otp, &forgetPass.ExpirationTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no OTP found for this email")
		}
		return nil, fmt.Errorf("error querying database: %w", err)
	}
	return &forgetPass, nil
}

func UpdatePasswordQueries(ctx context.Context, d *dbs.Service, f model.ForgetPass) error {
	queri := `
		UPDATE users
		SET password = ?
		WHERE email = ?
	`
	_, err := d.DB.ExecContext(ctx, queri, f.NewPass, f.Email)
	if err != nil {
		return fmt.Errorf("error updating password: %w", err)
	}
	return nil
}

func ClearOtpQueries(ctx context.Context, d *dbs.Service, email string) error {
	queri := `
		UPDATE users
		SET otp = NULL, otp_expiration = NULL
		WHERE email = ?
	`
	_, err := d.DB.ExecContext(ctx, queri, email)
	if err != nil {
		return fmt.Errorf("error clearing OTP: %w", err)
	}
	return nil
}

func UsernameQueries(ctx context.Context, d *dbs.Service, u model.User) (sql.Result, error) {
	queri := `UPDATE users SET username = ? WHERE id = ?`

	res, err := d.DB.ExecContext(ctx, queri, u.Username, u.ID)
	return res, err
}

func EmailQueries(ctx context.Context, d *dbs.Service, u model.User) (sql.Result, error) {
	queri := `UPDATE users SET email = ? WHERE id = ?`

	res, err := d.DB.ExecContext(ctx, queri, u.Email, u.ID)
	return res, err
}

func ForgetPassQueries(ctx context.Context, d *dbs.Service, f model.ForgetPass) error {
	queri := `
		UPDATE users
		SET password = ?
		WHERE email = ? 
		  AND otp = ?
		  AND otp_expiration > ?
	`

	_, err := d.DB.ExecContext(ctx, queri, f.NewPass, f.Email, f.Otp, time.Now())
	if err != nil {
		return fmt.Errorf("error updating password: %w", err)
	}

	return nil
}

func StoreOtpQueries(ctx context.Context, d *dbs.Service, f model.ForgetPass) error {
	queri := `
        UPDATE users
        SET otp = ?, otp_expiration = ?
        WHERE email = ?
    `
	result, err := d.DB.ExecContext(ctx, queri, f.Otp, f.ExpirationTime, f.Email)
	if err != nil {
		return fmt.Errorf("error updating user's OTP in the database: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no user found with email: %s", f.Email)
	}

	return nil
}
