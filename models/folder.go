package models

import (
	"bloc/database"
	"context"
	"log"

	"github.com/georgysavva/scany/pgxscan"
)

type Folder struct {
	Id    string `db:"id"      json:"id"`
	Name  string `db:"name"    json:"name"`
	Owner string `db:"f_owner" json:"owner"`
	Path  string `db:"path"    json:"path"`
}

func FolderCreate(folder Folder) error {
	// Create new folder in database
	_, err := database.DB.Exec(context.Background(), `INSERT INTO folders(id, name, f_owner, path) VALUES($1, $2, $3, $4)`, folder.Id, folder.Name, folder.Owner, folder.Path)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func FolderDelete(id string, owner string) error {
	// Delete folder from database
	_, err := database.DB.Exec(context.Background(), `DELETE FROM folders WHERE id = $1 AND f_owner = $2`, id, owner)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return err
}

func FolderList(owner string) ([]*Folder, error) {
	rows, err := database.DB.Query(context.Background(), `SELECT id, name, path FROM folders WHERE f_owner = $1`, owner)
	if err != nil {
		log.Println(err)
	}

	var folder []*Folder
	err = pgxscan.ScanAll(&folder, rows)
	if err != nil {
		log.Println(err.Error())
		return []*Folder{}, err
	}

	return folder, err
}
