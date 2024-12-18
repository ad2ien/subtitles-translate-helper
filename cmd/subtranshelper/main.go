package main

import (
	"flag"
	"subtitles-translate-helper/internal"
)

var logger = internal.GetLogger()

func main() {
	configPath := flag.String("config-path", "config.yml", "config file path")

	flag.Parse()

	logger.Println("Subtitle translate helper")
	logger.Println(" config file : ", *configPath)

	internal.InitConfig(*configPath)

	internal.StartTranslatorService()

	internal.TranslateFiles()

	internal.StopTranslatorService()

	logger.Println("Done. Thanks LibreTranslate : https://libretranslate.com/ ❤️")

}
