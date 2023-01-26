package main

import (
 "flag"
 "fmt"
 "path"
 "net/http"
 "os"
 "log"
 "io/ioutil"
 "net/url"
 "strings"
 "encoding/json"
 "bytes"
)

var (
 URL string
 API string
 format string

 vin string
 raw, yaml bool
 meta bool
 sparse bool
 fields string
)

// https://mholt.github.io/json-to-go/
type JSONOutput struct {
	Count          int    `json:"Count"`
	Message        string `json:"Message"`
	SearchCriteria string `json:"SearchCriteria"`
	Results        []struct {
		Value      string `json:"Value"`
		ValueID    string `json:"ValueId"`
		Variable   string `json:"Variable"`
		VariableID int    `json:"VariableId"`
	} `json:"Results"`
}


func removeEmpty(obj interface{}) JSONOutput {
    switch v := obj.(type) {
    case map[string]interface{}:
        for key := range v {
            if v[key] == nil {
                delete(v, key)
            } else if str, ok := v[key].(string); ok {
                if strings.TrimSpace(str) == "" {
                    delete(v, key)
                }
            } else if raw, ok := v[key].(json.RawMessage); ok {
                if bytes.Equal(raw, []byte("null")) {
                    delete(v, key)
                }
            }
        }
    }
    return obj.(JSONOutput)
}

func callAPI(URL string, API string, VIN string, format string, meta bool, sparse bool, fields string) {
 u, _ := url.Parse(URL)
 u.Path = path.Join(u.Path, API, VIN)

 client := &http.Client{}
 req, err := http.NewRequest(http.MethodGet, u.String(), nil)
 if err != nil {
  log.Fatal(err)
 }

 // appending to existing query args
 q := req.URL.Query()
 q.Add("format", format)

 // assign encoded query string to http request
 req.URL.RawQuery = q.Encode()

 resp, err := client.Do(req)
 if err != nil {
  fmt.Println("Errored when sending request to the server")
  os.Exit(1)
 }

 defer resp.Body.Close()
 responseBody, err := ioutil.ReadAll(resp.Body)
 //_, err = ioutil.ReadAll(resp.Body)
 if err != nil {
  log.Fatal(err)
 }

 if format == "json" {
  data := JSONOutput{}
  json.Unmarshal([]byte(string(responseBody)), &data)
  if sparse {
   data = removeEmpty(data)
  }
  if fields != "" {
  }
  if meta {
   fmt.Println("Count of fields and results:", data.Count)
  } else {
   fmt.Printf("%+v\n", data)
   //fmt.Println(data)
  }
 } else if format == "yaml" {
 }
 //fmt.Println(string(responseBody))
}

/*
fields - match name or regexp
sparse - non empty
meta - how many fields, results

service - over GRPC
help
*/

func init() {
 flag.StringVar(&vin, "vin", "3AKJHHDRXKSKX6687", "Query for given VIN")
 flag.BoolVar(&raw, "raw", true, "Output in JSON")
 flag.BoolVar(&yaml, "yaml", false, "Output in YAML")
 flag.BoolVar(&meta, "meta", false, "Get count of fields and results")
 flag.BoolVar(&sparse, "sparse", false, "Only output fields that have data in them")
 flag.StringVar(&fields, "fields", "", "Match name or regexp")
}

func main() {
 flag.Parse()
 URL = "https://vpic.nhtsa.dot.gov/api/"
 API = "/vehicles/DecodeVin/"

 if yaml && raw {
  args := strings.Join(os.Args, " ")
  if strings.Index(args, "yaml") > strings.Index(args, "raw") {
   raw = false
  } else {
   yaml = false
  }
 }
 if raw { format = "json" }
 if yaml { format = "yaml" }
 callAPI(URL, API, vin, format, meta, sparse, fields)
}
