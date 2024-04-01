# Telegram easy vpn admin

This is simple bot to manage clients from [wg-easy](https://github.com/wg-easy/wg-easy) with telegram. Currently it expects only one person - an admin and uses undocumented REST API.

## How to run your own bot

1. Register new bot using [BotFather](https://core.telegram.org/bots/features#creating-a-new-bot) and save api token somewhere.
2. Run wg-easy with docker or docker compose. Make sure to open dashboard port(51821 TCP by default) to your bot.
3. Provide these environment variables before run. You can add `export` to .bashrc or configure IDE.

    ```bash
    TELEGRAM_API_TOKEN=token from botfather
    TELEGRAM_USERNAME=case sensetive telegram username without@
    VPN_ADMIN_URL=http://<host>:<port>/
    VPN_ADMIN_PASSWORD=🚨YOUR_ADMIN_PASSWORD
    ```

4. Start bot `go run main.go`.

## Supported commands

* /ping - PONG
* /list - Show all clients.
* /create - Provide name and create a new client.
* /enable - Enable client.
* /disable - Disable client(disables VPN access).
* /delete - Delete client. Warning, client can not be restored.
* /config - Get config file for the device. Can be shared with other person to provide VPN access.
* ~~/code - Show QR Code to share config~~(command under development)
