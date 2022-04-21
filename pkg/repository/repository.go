package repository

type Repository interface {
	CreateNewLink(OriginalUrl string, ShortUrl string) error
	FindOriginalUrl(shortUrl string) (string, error)
	FindShortUrl(originalUrl string) (string, error)
}
