package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"rest_urlshort_postgres/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func (s *Storage) Ping() error {
	const op = "storage.sql.Ping"
	err := s.db.Ping()
	if err != nil {
		return fmt.Errorf("%s:::%w", op, err)
	}

	return nil
}
func New(storagePath string) (error, *Storage) {
	const op = "storage.sqlite.New"
	storage, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err), nil
	}

	q := `CREATE TABLE  IF NOT EXISTS url(
id INTEGER PRIMARY KEY,
alias TEXT NOT NULL UNIQUE,
url TEXT NOT NUll);
CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);`

	prepare, err := storage.Prepare(q)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err), nil
	}

	_, err = prepare.Exec()
	if err != nil {
		return fmt.Errorf("%s:%w", op, err), nil
	}

	return nil, &Storage{storage}
}

func (s *Storage) SaveUrl(urlTYoSave string, alias string) (int64, error) {
	const op = "storage.sql.SaveNew"
	q := `insert into url(url,alias) values(?,?)`

	stmt, err := s.db.Prepare(q)
	if err != nil {
		return -1, fmt.Errorf("%s:::%w", op, err)
	}

	exec, err := stmt.Exec(urlTYoSave, alias)
	if err != nil {
		if sqlerr, ok := err.(sqlite3.Error); ok && sqlerr.ExtendedCode == sqlite3.ErrConstraintUnique {

			return -1, fmt.Errorf("%s:::%w", op, storage.ErrUrlExist)

		}
		return -1, fmt.Errorf("%s:::%w", op, err)

	}
	id, err := exec.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("%s:>Fail to return last insert di<:%w", op, err)

	}
	return id, nil

}

func (s *Storage) GetUrl(alias string) (string, error) {

	const op = "storage.sqlite.GetUrl"

	q := `SELECT (url) FROM url where alias = ? `

	stmt, err := s.db.Prepare(q)
	if err != nil {

		return "", fmt.Errorf("%s :fail to create statement :::%w ", op, err)
	}

	var res string
	err = stmt.QueryRow(alias).Scan(&res)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrUrlNotFound

		}
		return "", fmt.Errorf("%s:::fail to get url from db ::: %w", op, err)

	}

	return res, nil

}
