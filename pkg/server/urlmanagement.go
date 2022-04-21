package server

import (
	"context"
	"link-shortener/pkg/repository"
	"link-shortener/pkg/util"
	pb "link-shortener/proto"
)

type UrlManagementServer struct {
	R   repository.Repository
	Cfg Config
	pb.UnimplementedUrlManagementServer
}

type Config struct {
	ServerAddress string
	LinkLength    int
}

func (s UrlManagementServer) Create(ctx context.Context, in *pb.OriginalUrl) (*pb.ShortUrl, error) {
	url, err := s.R.FindShortUrl(in.Url)
	if err != nil {
		var newShortUrl string
		for true {
			newShortUrl = generateNewShortUrl(s.Cfg.ServerAddress, s.Cfg.LinkLength)
			_, err := s.R.FindOriginalUrl(in.Url)
			if err != nil {
				break
			}
		}
		err := s.R.CreateNewLink(in.Url, newShortUrl)
		if err != nil {
			return nil, err
		}
		return &pb.ShortUrl{Url: newShortUrl}, nil
	}
	return &pb.ShortUrl{Url: url}, nil
}

func (s UrlManagementServer) Get(ctx context.Context, in *pb.ShortUrl) (*pb.OriginalUrl, error) {
	url, err := s.R.FindOriginalUrl(in.Url)
	if err != nil {
		return nil, err
	}
	return &pb.OriginalUrl{Url: url}, nil
}

func generateNewShortUrl(serverAddress string, linkLength int) string {
	return serverAddress + "/" + util.GenCode(linkLength)
}
