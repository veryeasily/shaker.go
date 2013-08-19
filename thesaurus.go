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
    str = ""
    for k, v := range *f {
        switch vv := v.(type) {
        case string:
            fmt.Println(k, "is string", vv)
            str = v.(string)
            break
        case int:
            fmt.Println(k, "is int", vv)
        case []interface{}:
            fmt.Println(k, "is an array")
            for i, u := range vv {
                fmt.Println(i, u)
                str = u.(string)
                break
            }
            break
        case map[string]interface{}:
            fmt.Println(k, "is an object")
            temp := v.(map[string]interface{})
            str = printJson(&temp)
        default:
            fmt.Println(vv, "is probably an object")
        }
        if str != "" {
            break
        }
    }
    fmt.Println(str)
    return str
}


func main() {
    var words = []string{"friendly", "bad"}
    if len(os.Args) > 1 {
        words = os.Args[1:]
    }
    str := ""
    var replacements []string
    length := len(words)
    replacements = make([]string, length, length)
    for i, word := range words {
        res, err := http.Get("http://words.bighugelabs.com/api/2/7c1a1031524ef2b6d72070ec9bcf5e5d/" + word + "/json")
        if err != nil {
            log.Fatal(err)
        }
        contents, err := ioutil.ReadAll(res.Body)
        defer res.Body.Close()
        if err != nil {
            log.Fatal(err)
        }
        var f map[string]interface{}
        err = json.Unmarshal(contents, &f)
        if err != nil {
            log.Fatal(err)
        }
        replacements[i] = printJson(&f)
    }
    str = strings.Join(replacements, " ")
    fmt.Println("")
    fmt.Println("")
    fmt.Println(str)
}
