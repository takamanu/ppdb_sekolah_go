FROM golang:alpine

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o ppdb_sekolah

EXPOSE 8000

ENTRYPOINT ["/app/ppdb_sekolah"]
