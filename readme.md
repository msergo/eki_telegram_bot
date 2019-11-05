### About
Telegram bot to query words from eki.ee

#### Env vars
BOT_TOKEN - telegram token
WEBHOOK_ADDRESS - host address of bot for setting webhook


## ToDo:
 * automatic deployment
 * generate nice messages
 * optimize code
 * move redis to another image
 * run via docker-compose
 * add logs with ELK stack

### Starting the bot
```bash
docker run -p 8083:8083 -e PORT=$PORT -e WEBHOOK_ADDRESS=$WEBHOOK_ADDRESS -e BOT_TOKEN=$BOT_TOKEN msergo/sergo_bot:eki_ee_bot
```