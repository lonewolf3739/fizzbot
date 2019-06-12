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

func main() {

	reader := bufio.NewReader(os.Stdin)

	startURL := "/fizzbot/questions/1"

	resp, err := http.Get(buildURI(startURL))
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	var data map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&data)

	fmt.Println(data["message"])

	var answer string

	answer, _ = reader.ReadString('\n')

	requestBody, mErr := json.Marshal(map[string]string{
		"answer": answer,
	})

	if mErr != nil {
		log.Fatalln(mErr)
	}

	postResp, postErr := http.Post(buildURI(startURL), "application/json", bytes.NewBuffer(requestBody))
	if postErr != nil {
		log.Fatalln(postErr)
	}
	defer postResp.Body.Close()
	json.NewDecoder(postResp.Body).Decode(&data)

	for data["nextQuestion"] != nil {

		tmpResp, tmpErr := http.Get(buildURI(data["nextQuestion"].(string)))
		if tmpErr != nil {
			log.Fatalln(tmpErr)
		}
		defer tmpResp.Body.Close()

		json.NewDecoder(tmpResp.Body).Decode(&data)

		fmt.Println(data)

		var answer string //= "1 2 Fizz 4 Buzz Fizz 7 8 Fizz Buzz 11 Fizz 13 14 FizzBuzz"

		answer, _ = reader.ReadString('\n')
		answer = answer[:len(answer)-1]

		requestBody, mErr := json.Marshal(map[string]string{
			"answer": answer,
		})

		if mErr != nil {
			log.Fatalln(mErr)
		}

		fmt.Println(answer)

		postResp, postErr := http.Post(buildURI(data["nextQuestion"].(string)), "application/json", bytes.NewBuffer(requestBody))
		if postErr != nil {
			log.Fatalln(postErr)
		}
		defer postResp.Body.Close()
		json.NewDecoder(postResp.Body).Decode(&data)

		fmt.Println(data)
	}

}
