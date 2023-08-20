FROM golang:alpine

WORKDIR /pvpc-backend
COPY go.mod go.sum ./
RUN go mod download

COPY cmd internal pkg ./

RUN go build -o ./bin/api ./cmd/api \
  && go build -o ./bin/migrate ./cmd/migrate

CMD ["/pvpc-backend/bin/api"]
EXPOSE 8080