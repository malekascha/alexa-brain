package main

import (
    "fmt"
    "net/http"
    "net/url"
    "io/ioutil"
)

func retrieveAuthCode() string {
  client_id := "amzn1.application-oa2-client.abe7a7cf539d4602bc92f11fbc646042"
  device_type_id := "my_device"
  device_serial_number := 123
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

func handler(w http.ResponseWriter, r *http.Request) {
    code := retrieveAuthCode();
    
    fmt.Fprintf(w, code)
}

func main() {
    // os.Setenv("FOO", "1")    
    http.HandleFunc("/", handler)
    http.ListenAndServe(":3000", nil)
}