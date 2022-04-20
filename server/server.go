package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	pb "link-shortener/gen/proto"
	"link-shortener/server/util"
	"log"
	"net"
)

const (
	serverAddress = ":8080"
	linkLength    = 10
	dbHost        = "localhost"
	dbPort        = 5432
	dbUser        = "postgres"
	dbPassword    = "postgres"
	dbName        = "postgres"
)

type Link struct {
	ID          int    `json:"id"`
	OriginalUrl string `json:"original_url"`
	ShortUrl    string `json:"short_url"`
}

type UrlManagementServer struct {
	db *sql.DB
	pb.UnimplementedUrlManagementServer
}

func (s UrlManagementServer) Create(ctx context.Context, in *pb.OriginalUrl) (*pb.ShortUrl, error) {
	url, err := FindShortUrl(s.db, in.Url)
	if err != nil {
		url = serverAddress + "/" + util.GenCode(linkLength)
		CreateNewLink(s.db, in.Url, url)
		return &pb.ShortUrl{Url: url}, nil
	}
	return &pb.ShortUrl{Url: url}, nil
}

func (s UrlManagementServer) Get(ctx context.Context, in *pb.ShortUrl) (*pb.OriginalUrl, error) {
	url, err := FindOriginalUrl(s.db, in.Url)
	if err != nil {
		return nil, err
	}
	return &pb.OriginalUrl{Url: url}, nil
}

func main() {
	lis, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := &UrlManagementServer{db: OpenConnection()}

	grpcServer := grpc.NewServer()

	pb.RegisterUrlManagementServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
}

func OpenConnection() *sql.DB {
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

	return db
}

func CreateNewLink(db *sql.DB, OriginalUrl string, ShortUrl string) error {
	sqlStatement := `INSERT INTO links (original_url, short_url) VALUES ($1, $2)`
	_, err := db.Exec(sqlStatement, OriginalUrl, ShortUrl)
	if err != nil {
		return err
	}
	return nil
}

func FindOriginalUrl(db *sql.DB, shortUrl string) (string, error) {
	sqlStatement := `SELECT * FROM links WHERE short_url = $1`
	var res Link
	err := db.QueryRow(sqlStatement, shortUrl).Scan(&res.ID, &res.OriginalUrl, &res.ShortUrl)
	if err != nil {
		return "", err
	}
	return res.ShortUrl, nil
}

func FindShortUrl(db *sql.DB, originalUrl string) (string, error) {
	sqlStatement := `SELECT * FROM links WHERE original_url = $1`
	var res Link
	err := db.QueryRow(sqlStatement, originalUrl).Scan(&res.ID, &res.OriginalUrl, &res.ShortUrl)
	if err != nil {
		return "", err
	}
	return res.ShortUrl, nil
}
