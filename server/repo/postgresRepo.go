package repo

import (
	"database/sql"
	"fmt"
	"link-shortener/server/model"
	"log"
)

const (
	dbHost     = "localhost"
	dbPort     = 5432
	dbUser     = "postgres"
	dbPassword = "postgres"
	dbName     = "postgres"
)

type postgresRepo struct {
	db *sql.DB
}

func OpenConnection() Repository {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Failed to connected to DB: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to pinged to DB: %v", err)
	}

	return &postgresRepo{db: db}
}

func (r *postgresRepo) CreateNewLink(OriginalUrl string, ShortUrl string) error {
	sqlStatement := `INSERT INTO links (original_url, short_url) VALUES ($1, $2)`
	_, err := r.db.Exec(sqlStatement, OriginalUrl, ShortUrl)
	if err != nil {
		return err
	}
	return nil
}

func (r *postgresRepo) FindOriginalUrl(shortUrl string) (string, error) {
	sqlStatement := `SELECT * FROM links WHERE short_url = $1`
	var res model.Link
	err := r.db.QueryRow(sqlStatement, shortUrl).Scan(&res.ID, &res.OriginalUrl, &res.ShortUrl)
	if err != nil {
		return "", err
	}
	return res.OriginalUrl, nil
}

func (r *postgresRepo) FindShortUrl(originalUrl string) (string, error) {
	sqlStatement := `SELECT * FROM links WHERE original_url = $1`
	var res model.Link
	err := r.db.QueryRow(sqlStatement, originalUrl).Scan(&res.ID, &res.OriginalUrl, &res.ShortUrl)
	if err != nil {
		return "", err
	}
	return res.ShortUrl, nil
}
