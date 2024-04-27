FROM golang:1.22.2

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY main.go ./
COPY pkg ./pkg
RUN go build -o vpn-admin

CMD ["./vpn-admin"]
