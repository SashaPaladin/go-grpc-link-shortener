package repository

import (
	"errors"
)

type inmemoryRepo struct {
	dbOrig  map[string]string
	dbShort map[string]string
}

func InmemoryDbCreate() Repository {
	dbOrig := make(map[string]string)
	dbShort := make(map[string]string)
	return inmemoryRepo{
		dbOrig:  dbOrig,
		dbShort: dbShort,
	}
}

func (i inmemoryRepo) CreateNewLink(originalUrl string, shortUrl string) error {
	i.dbOrig[originalUrl] = shortUrl
	i.dbShort[shortUrl] = originalUrl
	return nil
}

func (i inmemoryRepo) FindOriginalUrl(shortUrl string) (string, error) {
	res, ok := i.dbShort[shortUrl]
	if ok {
		return res, nil
	}
	return "", errors.New("FindOriginUrlError")
}

func (i inmemoryRepo) FindShortUrl(originalUrl string) (string, error) {
	res, ok := i.dbOrig[originalUrl]
	if ok {
		return res, nil
	}
	return "", errors.New("FindShortUrlError")
}
