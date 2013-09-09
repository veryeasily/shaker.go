package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
  "time"
)

// LOGGING bool
// constant that tells us if we're logging or not
var LOGGING bool = false
var requests chan Request
var NUMWORKERS int = 10

type Request struct {
	word           string
	index          int
	thesaurus_word chan Result
}

type Result struct {
    word string
    index int
}

func myFatal(err error) {
	if LOGGING {
		if err != nil {
			log.Fatal(err)
		}
	}
}

func myLog(str interface{}) {
	if LOGGING {
		fmt.Println(str)
	}
}

func getThesaurusWord(word string) string {
	myLog("Getting word")
  myLog(word)
	res, err := http.Get("http://words.bighugelabs.com/api/2/7c1a1031524ef2b6d72070ec9bcf5e5d/" + word + "/json")
  myLog("just got word")
	myFatal(err)
	contents, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	myFatal(err)
	if string(contents) != "" {
		var f map[string]interface{}
		err = json.Unmarshal(contents, &f)
		myFatal(err)
		if u := printJson(&f); u != "" {
			return u
		}
	}
	return word
}

func GetWords(words []string) []string {
	l := len(words)
	newWords := make([]string, l)
	results := make(chan Result, NUMWORKERS)

	messagesRecieved := 0
  // fmt.Printf("Number of words is %v\n", l)
  i := 0
  for i < l {
    request := Request{word: words[i], index: i, thesaurus_word: results}
    requests <- request
    // fmt.Println("gave a request")
    // fmt.Printf("i = %v\n", i)
    // fmt.Printf("l = %v\n", l)
    i++
  }
  // fmt.Printf("i == l is %v", i==l)
  // fmt.Printf("Left the loop?\n")

	for messagesRecieved < l {
    // fmt.Println("inside the loop")
    result := <-results
    // fmt.Println("inside of range")
    newWords[result.index] = result.word
    messagesRecieved += 1
    // fmt.Println(messagesRecieved)
	}
	return newWords
}

func main() {
  startingTime := time.Now()
	NUMWORKERS := 10
	var words = []string{"friendly", "bad"}
	if len(os.Args) > 1 {
		words = os.Args[1:]
	}

	requests = make(chan Request, NUMWORKERS)
	for i := 0; i < NUMWORKERS; i++ {
		go func(requests chan Request) {
			for request := range requests {
        // fmt.Println("got a request")
        // fmt.Println(request.word)
        newWord := getThesaurusWord(request.word)
        // fmt.Println("got word and pushing into results")
        // fmt.Println(newWord)
        request.thesaurus_word <- Result{word: newWord, index: request.index}
			}
		}(requests)
	}
  fmt.Printf("New words: %v\n", GetWords(words))
  endingTime := time.Now()
  fmt.Printf("Starting time was %v\n", startingTime.UnixNano())
  fmt.Printf("Ending time was %v\n", endingTime.UnixNano())
  fmt.Printf("Total nanoseconds elapsed is %v\n", endingTime.UnixNano() - startingTime.UnixNano())
}

func printJson(f *map[string]interface{}) (str string) {
	for k, v := range *f {
		switch vv := v.(type) {
		case []interface{}:
			// fmt.Println(k, "is an array", v)
			if k != "syn" && k != "sim" {
				break
			}
			str = vv[0].(string)
		case map[string]interface{}:
			// fmt.Println(k, "is an object")
			str = printJson(&vv)
		default:
			myLog(k)
			myLog("is weird")
		}
		if str != "" {
			break
		}
	}
	return str
}
