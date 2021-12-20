package models

import (
	"bloc/database"
	"context"
	"log"
	"time"

	"github.com/georgysavva/scany/pgxscan"
)

// File model
type File struct {
	Id       string    `db:"file_id"    json:"fileId"`    // UID of the filedb
	Owner    string    `db:"file_owner" json:"fileOwner"` // Username of the owner of the file
	Name     string    `db:"file_name"  json:"fileName"`  // Name of the file
	Type     string    `db:"file_type"  json:"fileType"`  // Mime type of the file
	Size     int64     `db:"file_size"  json:"fileSize"`  // Size of the file in bytes
	Chunk    int64     `db:"file_chunk" json:"-"`         // Size of the chunks of the encrypted file
	LastEdit time.Time `db:"last_edit"  json:"lastEdit"`  // Last time the file has been edited
}

// Format of file to send when getting a list
type FileFull struct {
	AccessId    string `db:"access_id"    json:"accessId"`
	AccessState string `db:"access_state" json:"accessState"`

	SharedBy string    `db:"shared_by" json:"sharedBy"`
	SharedTo string    `db:"shared_to" json:"sharedTo"`
	LastEdit time.Time `db:"last_edit" json:"lastEdit"`
	Favorite bool      `db:"favorite"  json:"favorite"`

	FileId   string `db:"file_id"   json:"fileId"`
	FileName string `db:"file_name" json:"fileName"`
	FileType string `db:"file_type" json:"fileType"`
	FileSize int64  `db:"file_size" json:"fileSize"`
}

type FileSharedList struct {
	AccessId string `db:"access_id" json:"accessId"`
	FileId   string `db:"file_id"   json:"fileId"`
	SharedBy string `db:"shared_by" json:"sharedBy"`
	SharedTo string `db:"shared_to" json:"sharedTo"`
}

// Add file in database
func FileCreate(file File) error {
	// Insert metadatas of the file
	_, err := database.DB.Exec(context.Background(), `INSERT INTO files(id, f_owner, name, type, size, chunk) VALUES($1, $2, $3, $4, $5, $6)`, file.Id, file.Owner, file.Name, file.Type, file.Size, file.Chunk)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

// Delete file from database
func FileDelete(id string, owner string) error {
	// Delete file access from database
	_, err := database.DB.Exec(context.Background(), `DELETE FROM files WHERE id = $1 AND f_owner = $2`, id, owner)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return err
}

func FileGet(id string) (File, error) {
	var file File
	err := pgxscan.Get(context.Background(), database.DB, &file, `SELECT
	id    AS file_id,
	f_owner AS file_owner,
	name  AS file_name,
	type  AS file_type,
	size  AS file_size,
	chunk AS file_chunk,
	last_edit
		FROM files
			WHERE id = $1`, id)

	if err != nil {
		log.Println(err.Error())
		return File{}, err
	}

	return file, nil
}

// Update favorite state
func FileUpdateFavorite(favorite bool, accessId string) error {
	_, err := database.DB.Exec(context.Background(), `UPDATE file_access SET favorite = $1 WHERE id = $2`, favorite, accessId)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

// Get size of a file
func FileGetSize(id string) (int64, error) {
	var size int64
	err := pgxscan.Get(context.Background(), database.DB, &size, `SELECT size FROM files WHERE id = $1`, id)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	return size, nil
}

/*************************************
 * FILE LIST
 * A bunch of function to list files with their access
 *************************************/

// List files or received files
func FileList(username string) ([]*FileFull, error) {
	req := `SELECT
	t1.id           AS access_id,
	t1.access_state AS access_state,
	t1.f_shared_by  AS shared_by,
	t1.f_shared_to  AS shared_to,
	t2.last_edit    AS last_edit,
	t1.favorite     AS favorite,
	t2.id           AS file_id,
	t2.name         AS file_name,
	t2.type         AS file_type,
	t2.size         AS file_size
		FROM file_access AS t1
			INNER JOIN files AS t2 ON t1.f_file = t2.id
				WHERE t1.f_shared_to = $1;`

	rows, err := database.DB.Query(context.Background(), req, username)
	if err != nil {
		log.Println(err)
	}

	var files []*FileFull
	err = pgxscan.ScanAll(&files, rows)
	if err != nil {
		log.Println(err.Error())
		return []*FileFull{}, err
	}

	return files, err
}

// List shared files
func FileListSharedBy(username string) ([]*FileSharedList, error) {
	rows, err := database.DB.Query(context.Background(), `SELECT
	t1.id           AS access_id,
	t2.id           AS file_id,
	t1.f_shared_by  AS shared_by,
	t1.f_shared_to  AS shared_to
		FROM file_access AS t1
			INNER JOIN files AS t2 ON t1.f_file = t2.id
				WHERE t1.f_shared_by = $1 AND t1.access_state = 'SHARED';`, username)

	if err != nil {
		log.Println(err)
	}

	var files []*FileSharedList
	err = pgxscan.ScanAll(&files, rows)
	if err != nil {
		log.Println(err.Error())
		return []*FileSharedList{}, err
	}

	return files, err
}
