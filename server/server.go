package main

import (
	"context"
	"flag"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	pb "link-shortener/gen/proto"
	"link-shortener/server/repo"
	"link-shortener/server/util"
	"log"
	"net"
)

const (
	serverAddress = ":8080"
	linkLength    = 10
)

type UrlManagementServer struct {
	r repo.Repository
	pb.UnimplementedUrlManagementServer
}

func (s UrlManagementServer) Create(ctx context.Context, in *pb.OriginalUrl) (*pb.ShortUrl, error) {
	url, err := s.r.FindShortUrl(in.Url)
	if err != nil {
		url = serverAddress + "/" + util.GenCode(linkLength)
		s.r.CreateNewLink(in.Url, url)
		return &pb.ShortUrl{Url: url}, nil
	}
	return &pb.ShortUrl{Url: url}, nil
}

func (s UrlManagementServer) Get(ctx context.Context, in *pb.ShortUrl) (*pb.OriginalUrl, error) {
	url, err := s.r.FindOriginalUrl(in.Url)
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

	flagPtr := flag.String("db", "", "a string")
	flag.Parse()

	var s *UrlManagementServer

	switch *flagPtr {
	case "postgres":
		s = &UrlManagementServer{r: repo.OpenConnection()}
	case "inmemory":
		s = &UrlManagementServer{r: repo.InmemoryDbCreate()}
	default:
		s = &UrlManagementServer{r: repo.OpenConnection()}
	}

	grpcServer := grpc.NewServer()

	pb.RegisterUrlManagementServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
}
