# About

This is a small Golang program that is intended to be executed by p4 triggers. When a new changelist is submitted, changelist details (change #, files changed, the user who submitted the change) are posted to a Discord channel via webhook.

My personal use case was to broadcast changelist submits to a Perforce depot that is being used by myself and a friend for a Unity game project.

# Requirements

* Perforce server that runs on some flavor of Linux. This was tested on CentOS 7.
* A webhook URL for a Discord channel. See here [this page](https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks&amp?page=3) for information on how to generate one
* A JSON-formatted configuration file that contains one key-value pair that will store the webhook URL. See `config-example.json` for an example. By default, the program will reference a config file located at `/etc/discord-trigger.conf` but other configs can be passed in by running the program with the `-conf` parameter: `-conf=/path/to/config.json`
