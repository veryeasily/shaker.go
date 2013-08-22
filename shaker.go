package main

import (
    "encoding/json"
	"fmt"
    "io/ioutil"
	"log"
	"net/http"
    "os"
    "strings"
)

// LOGGING bool
// constant that tells us if we're logging or not
var LOGGING bool = false

func myFatal(err error) {
    if LOGGING {
        if (err != nil) {
            log.Fatal(err)
        }
    }
}

func myLog(str interface{}) {
    if LOGGING {
        fmt.Println(str)
    }
}

func printJson(f *map[string]interface{}) (str string) {
    for k, v := range *f {
        switch vv := v.(type) {
            case []interface{}:
                // fmt.Println(k, "is an array", v)
                if k != "syn" && k != "sim" { break }
                str = vv[0].(string)
            case map[string]interface{}:
                // fmt.Println(k, "is an object")
                str = printJson(&vv)
            default:
                myLog(k)
                myLog("is weird")
        }
        if str != "" { break }
    }
    return str
}

func main() {
    var words = []string{"friendly", "bad"}
    if len(os.Args) > 1 {
        words = os.Args[1:]
    }
    length := len(words)
    done := make(chan bool)
    x := 0

    for i, word := range words {
        go func(word string, i int) {
            myLog("started go routine")
            res, err := http.Get("http://words.bighugelabs.com/api/2/7c1a1031524ef2b6d72070ec9bcf5e5d/" + word + "/json")
            // fmt.Println(word)
            myLog("made it here")
            myFatal(err)
            contents, err := ioutil.ReadAll(res.Body)
            defer res.Body.Close()
            myFatal(err)
            if string(contents) != "" {
                var f map[string]interface{}
                err = json.Unmarshal(contents, &f)
                myFatal(err)
                if u := printJson(&f); u != "" {
                    word = u
                }
            }
            words[i] = word
            x++
            myLog(x)
            if x == length {
                done <- true
            }
        }(word, i)
    }

    <-done
    str := strings.Join(words, " ")
    fmt.Println("")
    fmt.Println("")
    fmt.Println(str)
}
