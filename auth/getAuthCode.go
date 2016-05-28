package main

import (
    "fmt"
    "net/http"
    "net/url"
    "io/ioutil"
    "encoding/json"
)

type AuthEnvVars struct {
  Client_id string `json:"client_id"`
  Client_secret string `json:"client_secret"`
  Device_type_id string `json:"device_type_id"`
  Device_serial_number int `json:device_serial_number`
  Auth_code string `json:"auth_code"`
}

func retrieveAuthEnvVars() AuthEnvVars {
  data, err := ioutil.ReadFile("../authInfo.json")
  var auth_info AuthEnvVars

  err = json.Unmarshal(data, &auth_info)
  if(err != nil){
    panic(err)
  }
  return auth_info
}

func retrieveAuthCode() string {
  auth_info := retrieveAuthEnvVars()
  client_id := auth_info.Client_id
  device_type_id := auth_info.Device_type_id
  device_serial_number := auth_info.Device_serial_number
  redirect_uri := "https://localhost:3000/authresponse"
  response_type := "code"
  scope := "alexa:all"
  scope_data := fmt.Sprintf("{\"alexa:all\": {\"productID\": \"%s\", \"productInstanceAttributes\": {\"deviceSerialNumber\": \"%d\"}}}", device_type_id, device_serial_number)
  auth_url := fmt.Sprintf("https://www.amazon.com/ap/oa?client_id=%s&scope=%s&scope_data=%s&response_type=%s&redirect_uri=%s", client_id, url.QueryEscape(scope), url.QueryEscape(scope_data), response_type, url.QueryEscape(redirect_uri))
  res, err := http.Get(auth_url)
  if(err != nil){
      fmt.Println("Error getting auth code: ", err)
  }
  defer res.Body.Close()
  body, err := ioutil.ReadAll(res.Body)
  return string(body)
}

func authCodeRequestHandler(w http.ResponseWriter, r *http.Request) {
    code := retrieveAuthCode();  
    fmt.Fprintf(w, code)
}

func authCodeResponseHandler(w http.ResponseWriter, r *http.Request){
  fmt.Fprintf(w, "Query: %s",r.URL.Path)
}

func main() {
    http.HandleFunc("/", authCodeRequestHandler)
    http.HandleFunc("/authresponse", authCodeResponseHandler)
    fmt.Println("listening")
    http.ListenAndServe(":3000", nil)
}