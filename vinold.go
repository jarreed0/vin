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
)

var (
 URL string
 API string
)

func callAPI(URL string, API string, VIN string, format string) {
 u, _ := url.Parse(URL)
 u.Path = path.Join(u.Path, API, VIN + "?format=" + format)
 fmt.Println(u)
 response, err := http.Get(u.String())
 if err != nil {
  fmt.Print(err.Error())
  os.Exit(1)
 }
 responseData, err := ioutil.ReadAll(response.Body)
 //_, err = ioutil.ReadAll(response.Body)
 if err != nil {
  log.Fatal(err)
 }
 fmt.Println(string(responseData))
}

/*
raw - JSON
fields - match name or regexp
sparse - non empty
meta - how many fields, results
yaml - choose last from raw and yaml
service - over GRPC
help
*/

func main() {
 flag.Parse()

 URL = "https://vpic.nhtsa.dot.gov/api/"
 API = "/vehicles/DecodeVin/"

 VIN := "3AKJHHDR3KSKX6689"

 format := "json"

 callAPI(URL, API, VIN, format)
}
