package models

import (
	"bloc/database"
	"context"
	"log"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

type User struct {
	Username   string `db:"username"`
	Password   string `db:"password"`
	PublicKey  []byte `db:"public_key"`
	PrivateKey []byte `db:"private_key"`
	Quota      int64  `db:"quota"`
}

func UserExist(username string) bool {
	usr, err := UserGet(username)
	if err == pgx.ErrNoRows {
		return false
	}

	if usr.Username == username {
		return true
	}

	return false
}

func UserCreate(username string, password string, pubkey []byte, pvkey []byte) error {
	_, err := database.DB.Exec(context.Background(), `INSERT INTO users(username, password, public_key, private_key) VALUES($1, $2, $3, $4)`, username, password, pubkey, pvkey)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func UserGet(username string) (User, error) {
	var user User
	err := pgxscan.Get(context.Background(), database.DB, &user, `SELECT username, password, public_key, private_key, quota FROM users WHERE username = $1`, username)
	if err != nil {
		log.Println(err.Error())
		return User{}, err
	}

	return user, nil
}

func UserUpdateQuota(username string, quota int64) error {
	_, err := database.DB.Exec(context.Background(), `UPDATE users SET quota = $1 WHERE username = $2`, quota, username)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func UserGetQuota(username string) (int64, error) {
	var quota int64
	err := pgxscan.Get(context.Background(), database.DB, &quota, `SELECT quota FROM users WHERE username = $1`, username)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	return quota, nil
}

func UserGetPubKey(username string) ([]byte, error) {
	var pubkey []byte
	err := pgxscan.Get(context.Background(), database.DB, &pubkey, `SELECT public_key FROM users WHERE username = $1`, username)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return pubkey, nil
}
