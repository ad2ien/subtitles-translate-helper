package internal

import (
	"bufio"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func TranslateFiles() {

	logger.Println("Translate files from ", config.Input.Path, " to ", config.Output.Path)

	os.Mkdir(config.Output.Path, fs.ModePerm)

	files, err := os.ReadDir(config.Input.Path)
	if err != nil {
		StopTranslatorService()
		logger.Panic(err)
	}
	for _, file := range files {

		if err != nil {
			StopTranslatorService()
			logger.Panic(err)
		}
		translateFile(config.Input.Path + "/" + file.Name())
	}
}

func WriteLine(line string, file *os.File) {
	file.WriteString(line + "\n")
}

func translateFile(file string) {
	logger.Println("Translate file ", file)
	inFile, err := os.Open(file)
	if err != nil {
		log.Printf("Error opening file: %s", err)
		return
	}
	defer inFile.Close()

	outFile, errOpenFile := os.OpenFile(config.Output.Path+"/"+filepath.Base(file), os.O_WRONLY|os.O_CREATE, fs.ModePerm)
	if errOpenFile != nil {
		log.Printf("Error creating file: %s", errOpenFile)
		return
	}
	defer outFile.Close()

	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	var buffer string
	subLines := 0
	song := false
	// i := 0
	for scanner.Scan() {
		line := scanner.Text()
		stripedLine := strings.TrimSpace(line)
		lineType := getLineType(stripedLine)

		if lineType == SONG {
			WriteLine(stripedLine, outFile)
			song = true
		} else if lineType == DIALOGUE && song {
			WriteLine(stripedLine, outFile)
			song = false
		} else if lineType == DIALOGUE {
			buffer = buffer + " " + stripedLine
			subLines = subLines + 1
		} else if lineType == SEPARATOR {
			if buffer != "" {
				translation := translateSentence(buffer)
				if subLines > 1 {
					translation = splitLine(translation)
				}
				WriteLine(translation, outFile)
				buffer = ""
				subLines = 0
			}
			WriteLine("", outFile)
		} else {
			// SUB_NUMBER and TIMESTAMP
			WriteLine(stripedLine, outFile)
		}
		// i++
		// if i > 65 {
		// 	os.Exit(0)
		// }
	}
}
