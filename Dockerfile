FROM golang:1.22.4 AS build
WORKDIR /src
COPY go.mod go.sum .
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/vpn-admin ./main.go

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /bin/vpn-admin /bin/vpn-admin
ENTRYPOINT ["/bin/vpn-admin"]
