package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func buildURI(path string) string {
	return "https://api.noopschallenge.com" + path
}

func postAnswer(answer string, questionURL string) string {
	var data map[string]interface{}

	requestBody, err := json.Marshal(map[string]string{
		"answer": answer,
	})
	if err != nil {
		log.Fatalln(err)
	}

	postResp, postErr := http.Post(questionURL, "application/json", bytes.NewBuffer(requestBody))
	if postErr != nil {
		log.Fatalln(postErr)
	}
	defer postResp.Body.Close()
	json.NewDecoder(postResp.Body).Decode(&data)

	nextQuestionURL := data["nextQuestion"].(string)
	return nextQuestionURL
}

func solve(questionURL string) string {
	reader := bufio.NewReader(os.Stdin)

	resp, err := http.Get(questionURL)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	var data map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&data)

	fmt.Println(data["message"])
	if data["numbers"] != nil {
		fmt.Println(data["numbers"])
	}
	answer, _ := reader.ReadString('\n')
	answer = answer[:len(answer)-1]

	nextQuestionURL := postAnswer(answer, questionURL)

	return nextQuestionURL
}

func main() {
	fmt.Println("--------------------------------------------------------------------")
	questionURL := "/fizzbot/questions/1"

	for len(questionURL) > 0 {
		questionURL = solve(buildURI(questionURL))
		fmt.Println()
	}
}
