[![CircleCI](https://circleci.com/gh/IoTeXGroupIT/IoTeXGroupITBot.svg?style=svg)](https://circleci.com/gh/IoTeXGroupIT/IoTeXGroupITBot)
# IoTeXGroupITBot
A custom Telegram bot written in Go Lang to manage the [IoTeXGroupIT](http://t.me/IoTeXGroupIT) Telegram Group, especially to filter out chinese fake users which insert the spam message embedded in the First/Last Name fields.

## Features

### Spam Filter
- Kicks (but doesn't ban) a user from a group when he joins with **First Name** longer than a specific value
- Kicks (but doesn't ban) a user from a group when he joins with **Full Name** longer than a specific value

### Welcome message
A welcome message is posted when a new user joins the group if the corresponding settings are enabled in the configuration file. The welcome message has disabled notification, and each new welcome message deletes the previous one to avoid group cluttering 

### Commands
- **/help** Display the commands list
- **/roadmap** Display the current IoTeX Roadmap (image)

## Instructions
- Install Go Lang
- Build the bot with go lang (the executable `buildLinux` is provided to crosscompile for Linux )
- Copy the bot executable file to your server  
- On your server, set the Environment variable BOTAPIKEY with your Telegram Bot Api Key 
- Edit the bot options in configuration.json file according to your preferences
- run with the bot on the server (in Linux run with `./bot 2> bot_log &`)
