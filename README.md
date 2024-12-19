# Subtitle translate helper

[![CI status](https://img.shields.io/github/actions/workflow/status/ad2ien/subtitles-translate-helper/ci.yml?label=CI&logo=github)](https://github.com/ad2ien/subtitles-translate-helper/actions)
![GoLand](https://img.shields.io/badge/GoLand-0f0f0f?&logo=goland&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?logo=docker&logoColor=white)
[![Gitmoji](https://img.shields.io/badge/gitmoji-%20😜%20😍-FFDD67.svg)](https://gitmoji.dev)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

A Go tool to translate SRT subtiles using [LibreTranslate](https://github.com/LibreTranslate/LibreTranslate)

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

## Clean volumes and image

```sh
./subtitles-translate-helper --cleanup
```
