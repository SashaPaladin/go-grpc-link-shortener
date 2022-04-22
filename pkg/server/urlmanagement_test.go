package server

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	mockrepository "link-shortener/pkg/repository/mocks"
	pb "link-shortener/proto"
	"log"
	"net"
	"testing"
)

func TestCreate(t *testing.T) {
	type mockBehavior func(r *mockrepository.MockRepository, str1 string, str2 string)
	tests := []struct {
		name         string
		req          *pb.OriginalUrl
		res          *pb.ShortUrl
		mockBehavior mockBehavior
		errMsg       string
	}{
		{
			"valid request, when the short url in the database is found",
			&pb.OriginalUrl{Url: "originalurlindb.com/kekW"},
			&pb.ShortUrl{Url: "localhost/kekW_typogen"},
			func(r *mockrepository.MockRepository, reqUrl string, resUrl string) {
				r.EXPECT().FindShortUrl(reqUrl).Return(resUrl, nil)
			},
			"",
		},
		{
			"valid request, when a short URL is not found in the database",
			&pb.OriginalUrl{Url: "originalurlindb.com/kekW"},
			&pb.ShortUrl{Url: "/"},
			func(r *mockrepository.MockRepository, reqUrl string, resUrl string) {
				r.EXPECT().FindShortUrl(reqUrl).Return("", errors.New("sql: no rows in result set"))
				r.EXPECT().FindOriginalUrl(reqUrl).Return("", errors.New("sql: no rows in result set"))
				r.EXPECT().CreateNewLink(reqUrl, "/").Return(nil)
			},
			"",
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mockrepository.NewMockRepository(gomock.NewController(t))
			tt.mockBehavior(m, tt.req.Url, tt.res.Url)
			conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer(m)))
			if err != nil {
				log.Fatal(err)
			}
			client := pb.NewUrlManagementClient(conn)
			defer conn.Close()

			request := &pb.OriginalUrl{Url: tt.req.Url}

			response, err := client.Create(ctx, request)
			if response != nil {
				if response.Url != tt.res.Url {
					t.Error("response: expected", tt.res.Url, "received", response.Url)
				}
			}
			if err != nil {
				if er, ok := status.FromError(err); ok {
					if er.Message() != tt.errMsg {
						t.Error("error: expected", tt.errMsg, "received", er.Message())
					}
				}
			}
		})
	}
}

func TestGet(t *testing.T) {
	type mockBehavior func(r *mockrepository.MockRepository, str1 string, str2 string)
	tests := []struct {
		name         string
		req          *pb.ShortUrl
		res          *pb.OriginalUrl
		mockBehavior mockBehavior
		errMsg       string
	}{
		{
			"invalid request",
			&pb.ShortUrl{Url: ""},
			&pb.OriginalUrl{Url: ""},
			func(r *mockrepository.MockRepository, reqUrl string, resUrl string) {
				r.EXPECT().FindOriginalUrl(reqUrl).Return(resUrl, errors.New("sql: no rows in result set"))
			},
			"sql: no rows in result set",
		},
		{
			"valid request",
			&pb.ShortUrl{Url: "localhost/kekW_typogen"},
			&pb.OriginalUrl{Url: "originalurlindb.com/kekW"},
			func(r *mockrepository.MockRepository, reqUrl string, resUrl string) {
				r.EXPECT().FindOriginalUrl(reqUrl).Return(resUrl, nil)
			},
			"",
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mockrepository.NewMockRepository(gomock.NewController(t))
			tt.mockBehavior(m, tt.req.Url, tt.res.Url)
			conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer(m)))
			if err != nil {
				log.Fatal(err)
			}
			client := pb.NewUrlManagementClient(conn)
			defer conn.Close()

			request := &pb.ShortUrl{Url: tt.req.Url}

			response, err := client.Get(ctx, request)
			if response != nil {
				if response.Url != tt.res.Url {
					t.Error("response: expected", tt.res.Url, "received", response.Url)
				}
			}
			if err != nil {
				if er, ok := status.FromError(err); ok {
					if er.Message() != tt.errMsg {
						t.Error("error: expected", tt.errMsg, "received", er.Message())
					}
				}
			}
		})
	}
}

func dialer(repository *mockrepository.MockRepository) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	pb.RegisterUrlManagementServer(server, &UrlManagementServer{R: repository})

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}
