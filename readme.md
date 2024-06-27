# 🇪🇪 Estonian Dictionary Telegram Bot
[![Build status](https://dl.circleci.com/status-badge/img/gh/msergo/eki_telegram_bot/tree/master.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/msergo/eki_telegram_bot/tree/master)

## About

Telegram bot [@eki_ee_bot](https://t.me/eki_ee_bot) for querying translations from official [Estonian dictionaries](http://eki.ee/). Since there is no official API available, the bot parses HTML pages to fetch translations. Currently, it supports Estonian-Russian and Russian-Estonian translations, identifying the direction based on the charset. The bot does not have a permanent storage solution but utilizes Redis to cache already fetched articles. Messages are returned with inline keyboard for quick switching. 

## Motivation
The motivation behind this project was the need for a convenient tool to obtain translations quickly. The user interface of the eki.ee website is far away from being a user-friendly. As Telegram is my daily messenger, choosing it as the platform was an obvious decision. By that time, I had barely any knowledge of Go, so it was a good chance to learn by developing. Surprisingly, the bot works quite well and requires almost no maintenance.

![Screen](./screen.gif)

## Local run
Submit BOT_TOKEN and WEBHOOK_ADDRESS to env in the docker-compose.yaml and run it with `docker-compose up`. 
Now you're able to emulate webhook requests with `POST /localhost:8083`