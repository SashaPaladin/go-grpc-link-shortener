# go-grpc-link-shortener

### Для запуска приложения:

```
make build
make run_postgres
```

Если приложение запускается без базы данных:

```
make run_inmemory
```

Если приложение запускается впервые, необходимо применить миграции к базе данных:

```
make migrate
```

## Укорачиватель ссылок

### Сервис написано на Go и принимает следующие запросы по gRPC:
* Метод Create, который сохраняет оригинальный URL в базе и возвращает сокращённый URL
* Метод Get, который принимает сокращённый URL и возвращать оригинальный URL

### Сервис создаёт ссылки следующего формата:
* Ссылка уникальная и на один оригинальный URL ссылается только одна сокращенная ссылка.
* Ссылка длинной 10 символов (изменяется в конфигурации)
* Ссылка состоит из символов латинского алфавита в нижнем и верхнем регистре, цифр и символа _ (подчеркивание)

### Пример работы сервиса

```
message OriginalUrl {
  string url = "http://www.megairiginalurl.com/kekW";
}

// return "localhost:8000/NEcU5AgZiB"
rpc Create (OriginalUrl) returns (ShortUrl) {}

message ShortUrl {
  string url = "localhost:8000/NEcU5AgZiB";
}

//return "http://www.megairiginalurl.com/kekW"
rpc Get (ShortUrl) returns (OriginalUrl) {}
```