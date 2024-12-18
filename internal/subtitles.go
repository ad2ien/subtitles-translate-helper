package internal

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type LineType int

const (
	SUB_NUMBER LineType = iota
	TIMESTAMP
	SONG
	DIALOGUE
	SEPARATOR
	WTF
)

func getLineType(line string) LineType {
	if line == "" {
		return SEPARATOR
	}
	if regexp.MustCompile(`^[0-9]+$`).MatchString(line) {
		return SUB_NUMBER
	}
	if regexp.MustCompile(`^[0-9]{2}:[0-9]{2}:[0-9]{2}`).MatchString(line) {
		return TIMESTAMP
	}
	if strings.HasPrefix(line, config.IgnoreSubStartingWithChar) {
		return SONG
	}
	return DIALOGUE
}

func splitLine(translation string) string {
	nbOfChar := len(translation)
	tokens := strings.Split(translation, " ")
	count := 0
	result := []string{}
	done := false
	for i, token := range tokens {
		count = count + len(token) + 1
		result = append(result, token)
		if count > nbOfChar/2 && !done && i != len(tokens)-1 {
			result = append(result, "\n")
			done = true
		}
	}
	return strings.ReplaceAll(strings.Join(result, " "), " \n ", "\n")
}

func translateSentence(subtitle string) string {
	logger.Printf("Translate %s", subtitle)
	jsonArg := map[string]interface{}{
		"q":      subtitle,
		"source": config.Input.Lang,
		"target": config.Output.Lang,
	}
	// logger.Println("jsonArg")
	// logger.Println(jsonArg)
	jsonValue, err := json.Marshal(jsonArg)
	if err != nil {
		logger.Printf("%s", err)
		return "## Translation error json marshall"
	}
	url := "http://localhost:" + config.LibreTranslateServicePort + "/translate"
	// logger.Println("jsonValue : ", jsonValue)
	retries := 3
	var response *http.Response
	for i := range retries {
		response, err = http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			logger.Printf("%s", err)
			if i == retries-1 {
				return "## Translation error calling service"
			}
			time.Sleep(time.Millisecond * 200 * time.Duration(i+1))
		}
		break
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		logger.Printf("%d %s", response.StatusCode)
		body, err := io.ReadAll(response.Body)
		if err != nil {
			logger.Println(err)
		} else {
			logger.Println(string(body))
		}
		return "## Translation error service error"
	}

	var jsonResponse map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&jsonResponse)
	if err != nil {
		logger.Printf("Error %s", err)
		return "## Translation error service decode response"
	}

	if translatedText, ok := jsonResponse["translatedText"].(string); ok {
		return translatedText
	} else {
		return "## Translation error service response format error"
	}
}
