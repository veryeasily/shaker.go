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

type JSON struct {
    nodes map[string]*JSON
}

func printJson(f *map[string]interface{}) (str string) {
    for k, v := range *f {
        switch vv := v.(type) {
            case string:
                fmt.Println(k, "is string", vv)
                str = v.(string)
            case []interface{}:
                fmt.Println(k, "is an array")
                for i, u := range vv {
                    fmt.Println(i, u)
                    str = u.(string)
                    break
                }
            case map[string]interface{}:
                fmt.Println(k, "is an object")
                temp := v.(map[string]interface{})
                str = printJson(&temp)
        }
        if str != "" { break }
    }
    return str
}

type Word struct {
    word string
    i int
}

func main() {
    var words = []string{"friendly", "bad"}
    if len(os.Args) > 1 {
        words = os.Args[1:]
    }
    str, length := "", len(words)
    c := make(chan Word, length)
    for i, word := range words {
        go getWord(word, i, c)
    }
    fmt.Println("made it here")
    x := 0
    for thing := range c {
        words[thing.i] = thing.word
        x++
        if x == 10 { break }
    }
    str = strings.Join(words, " ")
    fmt.Println("")
    fmt.Println("")
    fmt.Println(str)
}

func getWord(word string, i int, c chan Word) {
    fmt.Println("Started a goroutine")
    res, err := http.Get("http://words.bighugelabs.com/api/2/7c1a1031524ef2b6d72070ec9bcf5e5d/" + word + "/json")
    if err != nil { log.Fatal(err) }
    contents, err := ioutil.ReadAll(res.Body)
    defer res.Body.Close()
    if err != nil { log.Fatal(err) }
    if string(contents) != "" {
        var f map[string]interface{}
        err = json.Unmarshal(contents, &f)
        if err != nil { log.Fatal(err) }
        word = printJson(&f)
    }
    fmt.Println(word)
    c <- Word{word, i}
}
