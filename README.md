# Subtitle translate helper

A Go tool to translate subtiles using [LibreTranslate](https://github.com/LibreTranslate/LibreTranslate)

## Prerequisites

- Linux station
- Docker is installed

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

`path` are input and output folders

### Downloading tool

On a Linux host:

```sh
curl -O https://github.com/ad2ien/subtitles-translate-helper/release/download/subtitles-translate-helper
chmod +x subtitles-translate-helper
./subtitles-translate-helper
```

### with Go

```sh
go run cmd/subtranshelper/main.go 
```

## Clean volumes

...

## todo

- [ ] check without the image in the first place
- [ ] specify subtitle format
- [ ] Version
- [ ] CI:
  - downloadable last artefact
  - lint?
- [ ] Badgess