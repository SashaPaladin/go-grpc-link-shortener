package repository

import (
	"database/sql"
	"fmt"
	"link-shortener"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type postgresRepo struct {
	db *sql.DB
}

func OpenConnection(cfg Config) (Repository, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open("postgres", psqlInfo)

	err = db.Ping()

	return &postgresRepo{db: db}, err
}

func (r *postgresRepo) CreateNewLink(OriginalUrl string, ShortUrl string) error {
	sqlStatement := `INSERT INTO links (original_url, short_url) VALUES ($1, $2)`
	_, err := r.db.Exec(sqlStatement, OriginalUrl, ShortUrl)
	return err
}

func (r *postgresRepo) FindOriginalUrl(shortUrl string) (string, error) {
	sqlStatement := `SELECT * FROM links WHERE short_url = $1`
	var res go_grpc_link_shortener.Link
	err := r.db.QueryRow(sqlStatement, shortUrl).Scan(&res.ID, &res.OriginalUrl, &res.ShortUrl)
	return res.OriginalUrl, err
}

func (r *postgresRepo) FindShortUrl(originalUrl string) (string, error) {
	sqlStatement := `SELECT * FROM links WHERE original_url = $1`
	var res go_grpc_link_shortener.Link
	err := r.db.QueryRow(sqlStatement, originalUrl).Scan(&res.ID, &res.OriginalUrl, &res.ShortUrl)
	return res.ShortUrl, err
}
