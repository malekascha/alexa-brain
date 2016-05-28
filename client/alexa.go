package main

import (
    "fmt"
    "net/http"
    // "net/url"
    "io/ioutil"
    "encoding/json"
)

type AuthResponse struct {
  Token_type string `json:"token_type"`
  Refresh_token string `json:"refresh_token"`
  Access_token string `json:"access_token"`
  Expires_in int `json:"expires_in"`
}


func fetchAccessToken() string {
  data, err := ioutil.ReadFile("token.json")
  var auth_response AuthResponse
  err = json.Unmarshal(data, &auth_response)
  if(err != nil){
    panic(err)
  }
  return auth_response.Access_token
}

func initDownchannel() {
  api_endpoint := "https://avs-alexa-na.amazon.com/v1/directives"
  access_token := fetchAccessToken() //retrieves token from local file
  req, err := http.NewRequest("GET", api_endpoint, nil)
  if(err != nil){
    panic(err)
  }
  req.Header.Add("authorization", fmt.Sprintf("Bearer %s", access_token))
  client := &http.Client{}
  res, err := client.Do(req)
  if(err != nil){
    panic(err)
  }
  contents, err := ioutil.ReadAll(res.Body)
  if(err != nil){
    panic(err)
  }
  \
}

func main() {
  initDownchannel()
}