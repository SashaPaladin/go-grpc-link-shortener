FROM golang:latest

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN apt-get update

RUN go mod download
RUN go build -o link-shortener ./cmd/main.go

EXPOSE 8000

CMD ["./link-shortener"]

