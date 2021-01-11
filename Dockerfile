FROM golang:1.14

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go get github.com/githubnemo/CompileDaemon

COPY . .

ENTRYPOINT CompileDaemon --build="go build main.go" --command=./main
