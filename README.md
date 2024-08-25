# VPN admin bot

This is simple Telegram bot to manage clients from [wg-easy](https://github.com/wg-easy/wg-easy). Currently it expects only one person - an admin and uses undocumented REST API. Last supported version 


## How to run your own bot

1. Register new bot using [BotFather](https://core.telegram.org/bots/features#creating-a-new-bot) and save api token somewhere.
2. Run wg-easy with docker or docker compose and make sure admin dashboard port is reachable from bot. Tested and supported version is [14](https://github.com/wg-easy/wg-easy/pkgs/container/wg-easy/255975322?tag=14).
3. Provide these environment variables before run.

    ```bash
    TELEGRAM_API_TOKEN=token from botfather
    TELEGRAM_USERNAME=case sensetive telegram username without@
    VPN_ADMIN_URL=http://<host>:<port>/
    VPN_ADMIN_PASSWORD_HASH=🚨YOUR_ADMIN_PASSWORD_HASH
    ```

4. Start bot `go run main.go`.

## Project structure
All initialization is placed in `main.go` file and business logic is in `pkg` folder.

## Development
Use command `go test ./...` in root folder to run tests and `golangci-lint run` to run linter.

## Supported commands

* /ping - PONG
* /list - Show all clients.
* /create - Provide name and create a new client.
* /enable - Enable client.
* /disable - Disable client(disables VPN access).
* /delete - Delete client. Warning, client can not be restored.
* /config - Get client config file. Can be shared with other person to provide VPN access.
* /qrcode - Get client config encoded in qrcode image with PNG format.
