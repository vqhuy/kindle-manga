# kindle-manga

`kindle-manga` periodically fetches your favorite manga(s) from
[truyentranhtuan.com](http://truyentranhtuan.com) and
[truyentranh.net](http://truyentranh.net) for latest chapters, converts them to
MOBI using [kcc](https://github.com/ciromattia/kcc), and then sends them to your
Kindle.

kcc uses **Kindle Oasis 2** profile by default. I was too lazy to put it into a
config file, so please forgive me and change [this
line](https://github.com/c633/kindle-manga/blob/master/kcc.go#L9) to your
device's profile.

## Installation

### Dependencies

- [kcc](https://github.com/ciromattia/kcc) and its dependencies (install them
  with `pip`).
- [KindleGen](http://www.amazon.com/gp/feature.html?ie=UTF8&docId=1000765211)
  v2.9+ in a directory reachable by your _PATH_.

You need to have Golang version 1.8 or higher installed. If Golang is set up
correctly, you can simply run:

```
go get github.com/c633/kindle-manga...
```

## Usage

1. You need to edit `config.toml` according to your manga list. This sample file
   should be pretty self-explanatory. Finally, copy it to
   `~/.config/kindle-manga` directory.

2. At the first run, `kindle-manga` will ask for your:

- approved Send-To-Kindle email & password
- Kindle's email

(Note that all these information are stored _in the clear_. Use at your own
risk.)

3. Set up a cron job

```
crontab -e
# 00 21 * * * /path/to/your/kindle-manga
```

## TODO

- automatically split file when file size is larger than 25MB (Gmail's message
  size limits)