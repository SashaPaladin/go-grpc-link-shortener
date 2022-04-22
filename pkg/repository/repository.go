package repository

//go:generate mockgen -source repository.go -destination mocks/mock.go

type Repository interface {
	CreateNewLink(OriginalUrl string, ShortUrl string) error
	FindOriginalUrl(shortUrl string) (string, error)
	FindShortUrl(originalUrl string) (string, error)
}
