package models

import (
	"bloc/database"
	"context"
	"log"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

type Access struct {
	Id            string `db:"access_id"      json:"accessId"`
	IsFolder      bool   `db:"is_folder"      json:"isFolder"`
	State         string `db:"access_state"   json:"accessState"`
	SharedBy      string `db:"shared_by"      json:"sharedBy"`
	SharedTo      string `db:"shared_to"      json:"sharedTo"`
	FileId        string `db:"file_id"        json:"fileId"`
	Favorite      bool   `db:"favorite"       json:"favorite"`
	EncryptionKey []byte `db:"encryption_key" json:"-"`
}

func AccessCreate(access Access) error {
	// Link the file to the user through the file_access table
	_, err := database.DB.Exec(context.Background(), `INSERT INTO
	file_access(
		id,
		is_folder,
		access_state,
		f_shared_by,
		f_shared_to,
		f_file,
		encryption_key
	) VALUES ($1, $2, $3, $4, $5, $6, $7)`, access.Id, access.IsFolder, access.State, access.SharedBy, access.SharedTo, access.FileId, access.EncryptionKey)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func AccessDelete(id string, owner string) error {
	// Delete file access from database
	_, err := database.DB.Exec(context.Background(), `DELETE FROM file_access WHERE f_shared_by = $1 AND f_file = $2`, owner, id)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	return err
}

func AccessGet(id string) (Access, error) {
	var access Access
	err := pgxscan.Get(context.Background(), database.DB, &access, `SELECT
	id           AS access_id,
	is_folder,
	access_state AS access_state,
	f_shared_by  AS shared_by,
	f_shared_to  AS shared_to,
	f_file       AS file_id,
	favorite,
	encryption_key
		FROM file_access
			WHERE id = $1`, id)

	if err != nil {
		log.Println(err.Error())
		return Access{}, err
	}

	return access, nil
}

func AccessExist(sharedTo string, fileId string) bool {
	var access string
	err := pgxscan.Get(context.Background(), database.DB, &access, `SELECT id FROM file_access WHERE f_shared_to = $1 AND f_file = $2`, sharedTo, fileId)
	if err == pgx.ErrNoRows {
		return false
	}

	if access != "" {
		return true
	}

	return false
}
