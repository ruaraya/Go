package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
)

var wordList []string
var isLetter = regexp.MustCompile(`^[a-z]+$`).MatchString

/** Applies validation in the input parameters
 */
func validateInput(word string) string {

	var errorMessage string

	if (word == "") || (len(word) < 2) {
		errorMessage = "word should be minimun 2 characthers long"
		return errorMessage
	} else if isLetter(word) == false {
		errorMessage = "only characters allowed"
		return errorMessage
	} else {
		return "OK"
	}

}

/** Validates if two words share the same characters
 */
func hasPattern(word1 string, word2 string) bool {

	match1 := true
	match2 := true
	final := true

	for i := 0; i < len(word1); i++ {
		char := string(word1[i])
		match1 = strings.Contains(word2, string(char))

		if match1 == false {
			return false
		}
	}

	for i := 0; i < len(word2); i++ {
		char := string(word2[i])
		match2 = strings.Contains(word1, string(char))

		if match2 == false {
			return false
		}
	}

	if match1 == true && match2 == true {
		final = true
	}
	return final
}

/**
 * @api {get} /find Find Anagrams
 * @apiDescription This endpoint will find all anagrams in the english dictionary based on the string sent
 * @apiParam (query) {String} word
 * @apiExample {curl} Example usage:
 *   curl -X GET -H "Content-Type: application/json" http://localhost:3001/find?word=test
 *
 * @apiSuccessExample {json} Success-Response:
 *   HTTP/1.1 200 OK
 *   [
 *      "word1",
 *      "word2",
 *      "word3"
 *   ]
 */
func getFindEndpoint(w http.ResponseWriter, req *http.Request) {

	var word string
	var errorMessage string
	var match1 bool
	var match2 bool
	var anagrams []string

	params := mux.Vars(req)
	word1 := strings.ToLower(params["word"])
	errorMessage = validateInput(word1)

	if errorMessage != "OK" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

	for i := 0; i < len(wordList); i++ {
		word = wordList[i]

		if len(word1) == len(word) {

			if word1 == word {
				match1 = true
				continue
			}

			if hasPattern(word1, word) == true {
				anagrams = append(anagrams, word)
				match2 = true
			}

		}

	}

	if match1 == false || match2 == false {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(false)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(anagrams)
	}

}

/**
 * @api {get} /compare Compare Anagrams
 * @apiDescription This endpoint will receive two words, and compare them to see if they are anagrams
 * @apiGroup Anagram
 * @apiParam (query) {String} word1
 * @apiParam (query) {String} word2
 * @apiExample {curl} Example usage:
 *   curl -X GET -H "Content-Type: application/json" http://localhost:3001/compare?word1=test&word2=tset
 *
 * @apiSuccessExample {json} Success-Response:
 *   HTTP/1.1 200 OK
 *   false
 */
func getCompareEndpoint(w http.ResponseWriter, req *http.Request) {

	var word string
	var errorMessage string
	var match1 bool
	var match2 bool

	params := mux.Vars(req)

	word1 := strings.ToLower(params["word1"])
	word2 := strings.ToLower(params["word2"])

	errorMessage = validateInput(word1)

	if errorMessage == "OK" {
		errorMessage = validateInput(word2)
	}

	if errorMessage != "OK" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

	if word1 == word2 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("words should not be equal")
		return
	}

	if (len(word1) == len(word2)) && (word1 != word2) {

		if hasPattern(word1, word2) == false {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(false)
			return
		} else {

			for i := 0; i < len(wordList); i++ {
				word = wordList[i]

				if word == word1 {
					match1 = true
					continue
				}

				if word == word2 {
					match2 = true
					continue
				}

				if match1 == true && match2 == true {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(true)
					return
				}

			}

			if match1 == false || match2 == false {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(false)
				return
			}

		}

	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(false)
	}

}

func main() {

	file, err := os.Open("words.txt")

	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wordList = append(wordList, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return
	}

	router := mux.NewRouter()
	router.Path("/find").Queries("word", "{word}").HandlerFunc(getFindEndpoint).Methods("GET")
	router.Path("/compare").Queries("word1", "{word1}", "word2", "{word2}").HandlerFunc(getCompareEndpoint).Methods("GET")
	log.Fatal(http.ListenAndServe(":3001", router))
}
