package main

import (
    "encoding/json"
	"fmt"
	"log"
	"net/http"
    "io/ioutil"
)

type JSON struct {
    nodes map[string]*JSON
}

func printJson(f *map[string]interface{}) {
    for k, v := range *f {
        switch vv := v.(type) {
        case string:
            fmt.Println(k, "is string", vv)
        case int:
            fmt.Println(k, "is int", vv)
        case []interface{}:
            fmt.Println(k, "is an array")
            for i, u := range vv {
                fmt.Println(i, u)
            }
        case map[string]interface{}:
            fmt.Println(k, "is an object")
            temp := v.(map[string]interface{})
            printJson(&temp)
        default:
            fmt.Println(vv, "is probably an object")
        }
    }
}


func main() {
    var words = []string{"friendly", "bad"}
    for _, word := range words {
        res, err := http.Get("http://words.bighugelabs.com/api/2/7c1a1031524ef2b6d72070ec9bcf5e5d/" + word + "/json")
        if err != nil {
            log.Fatal(err)
        }
        contents, err := ioutil.ReadAll(res.Body)
        if err != nil {
            log.Fatal(err)
        }
        var f map[string]interface{}
        err = json.Unmarshal(contents, &f)
        if err != nil {
            log.Fatal(err)
        }
        defer res.Body.Close()
        printJson(&f)
    }
}
