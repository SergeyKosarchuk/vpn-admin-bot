FROM golang:1.20.4

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /go-vnp-admin
CMD ["/go-vnp-admin"]
