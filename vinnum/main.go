package main

import (
"net"
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
 "google.golang.org/grpc"
)

var (
	vinFlag     = flag.String("vin", "", "the VIN to decode")
	rawFlag     = flag.Bool("raw", false, "output the raw JSON response")
	fieldsFlag  = flag.String("fields", "", "output the fields that match the provided name or regular expression")
	sparseFlag  = flag.Bool("sparse", false, "only output fields that have data")
	metaFlag    = flag.Bool("meta", false, "output metadata about the response")
	yamlFlag    = flag.Bool("yaml", false, "output the results in YAML format")
	serviceFlag = flag.Bool("service", false, "expose the utility as a gRPC service")
)

func jsonToYAML(jsonStr string) string {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		panic(err)
	}
	yamlStr, err := yaml.Marshal(data)
	if err != nil {
		panic(err)
	}
	return string(yamlStr)
}

func apiCall(URL string) string {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		log.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("format", "json")
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Problem sending request to server.")
		os.Exit(1)
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(responseBody)
}

func main() {
	flag.Parse()

if *serviceFlag {
 listen, err :=	net.Listen("tcp", ":9000")
 if err != nil {
  log.Fatalf("Failed  to listen to port 9000: %v", err)
 }
 grpcServer := grpc.NewServer()
 if err := grpcServer.Serve(listen); err != nil {
  log.Fatalf("Failed to serve at grpc server of port 9000: %v", err)
 }
}

	if *vinFlag == "" {
		fmt.Println("Enter the VIN to decode: ")
		fmt.Scanln(vinFlag)
	}
	u, _ := url.Parse("https://vpic.nhtsa.dot.gov/")
	u.Path = path.Join(u.Path, "api/vehicles/DecodeVin/", *vinFlag)

	if *yamlFlag && *rawFlag {
		args := strings.Join(os.Args, " ")
		if strings.Index(args, "yaml") > strings.Index(args, "raw") {
			*rawFlag = false
		} else {
			*yamlFlag = false
		}
	} else if !*yamlFlag && !*rawFlag && !*metaFlag {
		*rawFlag = true
	}

	data := apiCall(u.String())
	if *sparseFlag {
		re := regexp.MustCompile(`\{[^{}]*?(?:null)[^{}]*?\}`)
		data = re.ReplaceAllString(data, "")
		re = regexp.MustCompile(`\{[^{}]*?(?:"")[^{}]*?\}`)
		data = re.ReplaceAllString(data, "")
		re = regexp.MustCompile(`,+`)
		data = re.ReplaceAllString(data, ",")
		data = strings.Replace(data, "[,", "[", -1)
		data = strings.Replace(data, "{,", "{", -1)
		data = strings.Replace(data, ",]", "]", -1)
		data = strings.Replace(data, ",}", "}", -1)
	}
	if *fieldsFlag != "" {
		re := regexp.MustCompile(*fieldsFlag)
		matches := re.FindAllString(data, -1)
		for _, match := range matches {
			fmt.Println(match)
		}
	}
	if *metaFlag {
		re := regexp.MustCompile("{[^}]*}")
		matches := re.FindAllString(data, -1)
		fmt.Printf("{'Count': %d}", len(matches))
	}
	if *rawFlag {
		fmt.Printf("%+v\n", data)
	} else if *yamlFlag {
		fmt.Println(jsonToYAML(data))
	}
}
