package main

import (
  "bytes"
  "fmt"
  "net/http"
  "net/url"
  "io/ioutil"
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

func retrieveTokenQuery() string{
  auth_info := retrieveAuthEnvVars()
  client_id := auth_info.Client_id
  client_secret := auth_info.Client_secret
  code := auth_info.Auth_code
  grant_type := "authorization_code"
  redirect_uri := "https://localhost:3000/authresponse"
  access_url :="https://api.amazon.com/auth/o2/token/"
  post_body := url.Values{}
  post_body.Add("grant_type", grant_type)
  post_body.Add("code", code)
  post_body.Add("client_id", client_id)
  post_body.Add("client_secret", client_secret)
  post_body.Add("redirect_uri", redirect_uri)
  req, err := http.NewRequest("POST", access_url, bytes.NewBufferString(post_body.Encode()))
  // fmt.Sprintf("grant_type=%s&code=%s&client_id=%s&client_secret=%s&redirect_uri=%s", grant_type, code, client_id, client_secret, url.QueryEscape(redirect_uri))
  client := &http.Client{}
  req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
  res, err := client.Do(req)
  if(err != nil){
      fmt.Println("Error getting auth code: ", err)
  }
  defer res.Body.Close()
  body, err := ioutil.ReadAll(res.Body)
  return string(body)
}

func tokenResponseHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "body: %s", retrieveTokenQuery())
}

func main() {
  http.HandleFunc("/", tokenResponseHandler)
  fmt.Println("listening")
  http.ListenAndServe(":3000", nil)
}
