FROM golang:1.21.4 as build

WORKDIR /src
COPY go.mod go.sum ./
COPY main.go ./
COPY pkg ./pkg
RUN CGO_ENABLED=0 GOOS=linux go build -o /admin-bot

FROM ubuntu:latest
RUN apt update && apt upgrade -y
COPY --from=build /admin-bot /admin-bot
CMD ["/admin-bot"]
