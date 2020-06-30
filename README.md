# gitnotify-rss

Beware, the code is a mess. It works fine, I just don't feel like reorganizing it.

## How to build?
```sh
git clone https://github.com/toluschr/gitnotify-rss
cd gitnotify-rss
go build -o gitnotify-rss
```

## How to use?
 - [Create a new API-Token](https://help.github.com/en/github/authenticating-to-github/creating-a-personal-access-token)
with notification permissions.
 - Add the following feed to your rss reader `http://localhost:8092/${USERNAME}:${TOKEN}`, replace username and token
 - Have unread notifications be added to your rss feed
 - Run the server (`./gitnotify-rss`)
