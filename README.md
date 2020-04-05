# Daisy Mae

## Animal Crossing New Horizons Turnip Discord Bot

**Daisy Mae** is a discord bot specifically for Animal Crossing New Horizons made with [DiscordGo](https://github.com/bwmarrin/discordgo). The goal is to make turnip trading
and notifications for a discord server eaiser! (It's also a good way for me to practice Go ;) )

**NOTE:** I am not planning on making Daisy Mae into a moderation bot as there are a lot of good ones existing
at the moment! Thanks!

I will not be providing a server for this bot. If you would like to use Daisy Mae, please check out the documentation
below.

---

## Preliminary Requirements

This bot uses [PostgreSQL](https://www.postgresql.org/) to store data on _New Horizons'_ bug and fish. Make sure to set the database
up or the bot will show error messages.

To setup postgres, [download](https://www.postgresql.org/download/) the binary for your system, and setup a user account.

After, import the data using `setBugFish.psql` and `data.csv` (A easy way of doing this is using `pwd` and then append `data.csv`). If the database works by then, the code will handle everything else.

**Please NOTE:** The bot is still under development so files are continually changing; Current files are subject to change.

---

## Usage

1. Clone or Download Repository or run `go get github.com/yiping-allison/daisymae`
2. Configure example.config file and rename the file to: `.config`
3. Run `go build`
4. Run executable

## Bugs & Contributing

This is my first time writing a discord bot! I welcome any help or bug reports!

## Contributions

* Developers of [DiscordGo](https://github.com/bwmarrin/discordgo)
