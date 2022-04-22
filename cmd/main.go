package main

import (
	"flag"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"link-shortener/pkg/repository"
	"link-shortener/pkg/server"
	pb "link-shortener/proto"
	"log"
	"net"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("faliled to initialize configs: %s", err)
	}

	servCfg := server.Config{
		ServerAddress: viper.GetString("server.host") + ":" + viper.GetString("server.port"),
		LinkLength:    viper.GetInt("server.linkLength"),
	}

	lis, err := net.Listen("tcp", ":"+viper.GetString("server.port"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	flagPtr := flag.String("db", "", "a string")
	flag.Parse()

	var s *server.UrlManagementServer

	switch *flagPtr {
	case "postgres":
		db, err := repository.OpenConnection(repository.Config{
			Host:     viper.GetString("db.host"),
			Port:     viper.GetString("db.port"),
			Username: viper.GetString("db.username"),
			Password: viper.GetString("db.password"),
			DBName:   viper.GetString("db.dbname"),
			SSLMode:  viper.GetString("db.sslmode"),
		})
		if err != nil {
			log.Fatalf("failed to initialize db: %s", err)
		}
		s = &server.UrlManagementServer{R: db, Cfg: servCfg}
	case "inmemory":
		s = &server.UrlManagementServer{R: repository.InmemoryDbCreate(), Cfg: servCfg}
	default:
		log.Fatalf("faliled to initialize server: Enter flag -db=postgres or -db=inmemory")
	}

	grpcServer := grpc.NewServer()

	pb.RegisterUrlManagementServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
}

func initConfig() error {
	viper.AddConfigPath("cfg")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
