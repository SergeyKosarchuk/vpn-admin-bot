FROM golang:1.20.4 as build

WORKDIR /src
COPY go.mod go.sum ./
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /admin-bot


FROM ubuntu:latest
COPY --from=build /admin-bot /admin-bot
CMD ["/admin-bot"]
