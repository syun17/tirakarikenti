FROM golang:1.24-alpine

WORKDIR /app

# Install air for hot reload
RUN go install github.com/air-verse/air@v1.61.7

# ★ go.sumもコピーする
COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]