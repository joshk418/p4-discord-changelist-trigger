# About

This is a small Golang program that is intended to be executed by p4 triggers. When a new changelist is submitted, changelist details (change #, files changed, the user who submitted the change) are posted to a Discord channel via webhook.

My personal use case was to broadcast changelist submits to a Perforce depot that is being used by myself and a friend for a Unity game project.

# Requirements

* Perforce server
* A webhook URL for a Discord channel. See here [this page](https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks&amp?page=3) for information on how to generate one
