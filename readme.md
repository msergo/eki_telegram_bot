# Estonian Dictionary Telegram Bot

## About
Telegram bot for querying translations from official [Estonian dictionaries](http://eki.ee/). Since there is no official API available, the bot parses HTML pages to fetch translations. Currently, it supports Estonian-Russian and Russian-Estonian translations, identifying the direction based on the charset. The bot does not have a permanent storage solution but utilizes Redis to cache already fetched articles. Messages are returned with inline keyboard for quick switching. 

## Motivation
The motivation behind this project was the need for a convenient tool to obtain translations quickly. The user interface of the eki.ee website is far away from being a user-friendly. As Telegram is my daily messenger, choosing it as the platform was an obvious decision. By that time, I had barely any knowledge of Go, so it was a good chance to learn by developing. Surprisingly, the bot works quite well and requires almost no maintenance.


![Screen](./screen.gif)