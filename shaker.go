package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// LOGGING bool
// constant that tells us if we're logging or not
var LOGGING bool = true
var requests chan Request

type Request struct {
	word           string
	index          uint
	thesaurus_word chan Result
}

type Result struct {
    word string
    index uint
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
	return ""
}

func GetWords(words []string) []string {
	l := uint(len(words))
	newWords := make([]string, l)
	results := make(chan Result, l)

	i := uint(0)
	messagesRecieved := uint(0)
	for {
		if messagesRecieved == l {
			break
		}
		select {
		case result := <-results:
			newWords[result.index] = result.word
			messagesRecieved += 1
		default:
      fmt.Printf("messagesRecieved: %v\n", messagesRecieved)
			if i != l {
        request := Request{word: words[i], index: i, thesaurus_word: results}
				requests <- request
        fmt.Println("gave a request")
				i += 1
			}
		}
	}
	return newWords
}

func main() {
	NumWorkers := 10
	var words = []string{"friendly", "bad"}
	if len(os.Args) > 1 {
		words = os.Args[1:]
	}

	requests = make(chan Request, NumWorkers)
	for i := 0; i < NumWorkers; i++ {
		go func(requests chan Request) {
			for request := range requests {
        fmt.Println("got a request")
        fmt.Println(request.word)
        newWord := getThesaurusWord(request.word)
        fmt.Println(newWord)
        request.thesaurus_word <- Result{word: newWord, index: request.index}
			}
		}(requests)
	}
  fmt.Printf("New words: %v\n", GetWords(words))
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
