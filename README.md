# IoTeXGroupITBot
A custom Telegram bot written in Go Lang to manage the IoTeXGroupIT Telegram Group, especially to filter out chinese fake users which insert the spam message embedded in the First/Last Name fields.

## Features
- Kicks (but doesn't ban) a user from a group when he joins with **First Name** longer than a specific value
- Kicks (but doesn't ban) a user from a group when he joins with **Full Name** longer than a specific value

## Instructions
- Install Go Lang
- Build the bot with go lang (crosscompile command for Linux *buildLinux* is provided)
- Set the Environment variable BOTAPIKEY with your Telegram Bot Api Key or
- optionally, put your Api Key into the configuration.json file
- Edit the options in configuration.json file according to your spam filter preferences
- run with ./bot
