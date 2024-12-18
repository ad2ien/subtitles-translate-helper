# Subtitle translate helper

A Go tool to translate subtiles using [LibreTranslate](https://github.com/LibreTranslate/LibreTranslate)

## Prerequisites

- Docker

## todo

- [ ] check without the image in the first place
- [ ] specify subtitle format
- [ ] Tidy mode
- [ ] CI

## Build

```sh
go build ./cmd/subtranshelper
```

## Start

## Config

Fill the config file named `config.yml`. Ex:

```yaml
input:
  path: "english"
  lang: "en"

output:
  path: "french"
  lang: "fr"

ignoreSubStartingWithChar: "♪"
libreTranslateServicePort: 5021
libreTranslateImageVersion: "latest"
```

### Downloading tool

On a Linux host:

```sh
curl https://github.com/ad2ien/subtitles-translate-helper/release ...
chmod +x subtitles-translate-helper
./subtitles-translate-helper
```

### with Go

```sh
go run cmd/subtranshelper/main.go 
```

## Clean volumes

...
