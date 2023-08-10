FROM golang:alpine

WORKDIR /go-pvpc
COPY go.mod go.sum ./
RUN go mod download

COPY cmd internal pkg ./

RUN go build -o ./bin/api ./cmd/api \
  && go build -o ./bin/migrate ./cmd/migrate

CMD ["/go-pvpc/bin/api"]
EXPOSE 8080