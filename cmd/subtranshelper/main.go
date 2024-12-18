package main

import (
	"flag"
	"subtitles-translate-helper/internal"
)

var logger = internal.GetLogger()

func main() {
	configPath := flag.String("config-path", "config.yml", "config file path")
	cleanup := flag.Bool("cleanup", false, "Only cleanup images and volumes")

	flag.Parse()
	logger.Println("Subtitle translate helper")

	internal.InitConfig(*configPath)

	if *cleanup {
		internal.CleanUp()
		return
	} else {
		runScript(configPath)
	}

}

func runScript(configPath *string) {
	logger.Println(" config file : ", *configPath)

	internal.StartTranslatorService()

	internal.TranslateFiles()

	internal.StopTranslatorService()

	logger.Println("Done. Thanks LibreTranslate : https://libretranslate.com/ ❤️")

}
