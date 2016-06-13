package main

import (
    "bytes"
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

type AuthResponse struct {
  Token_type string `json:"token_type"`
  Refresh_token string `json:"refresh_token"`
  Access_token string `json:"access_token"`
  Expires_in int `json:"expires_in"`
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

func fetchRefreshToken() string {
  data, err := ioutil.ReadFile("../token.json")
  var auth_response AuthResponse

  err = json.Unmarshal(data, &auth_response)
  if(err != nil){
    panic(err)
  }
  return auth_response.Refresh_token
}

func retrieveAuthToken() []byte {
  auth_info := retrieveAuthEnvVars()
  auth_url := "https://api.amazon.com/auth/o2/token"
  client_id := auth_info.Client_id
  client_secret := auth_info.Client_secret
  redirect_uri := "https://localhost:3000/authresponse"
  grant_type := "refresh_token"
  refresh_token := fetchRefreshToken()
  post_body := url.Values{}
  post_body.Add("client_id", client_id)
  post_body.Add("client_secret", client_secret)
  post_body.Add("redirect_uri", url.QueryEscape(redirect_uri))
  post_body.Add("grant_type", grant_type)
  post_body.Add("refresh_token", refresh_token)
  // auth_url := fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=%s&redirect_uri=%s&refresh_token=%s", client_id, client_secret, grant_type, url.QueryEscape(redirect_uri), refresh_token)
  // res, err := http.Get(auth_url)
  client := &http.Client{}
  r, err := http.NewRequest("POST", auth_url, bytes.NewBufferString(post_body.Encode()))
  r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
  res, err := client.Do(r)
  if(err != nil){
      fmt.Println("Error getting new token: ", err)
  }
  defer res.Body.Close()
  body, err := ioutil.ReadAll(res.Body)
  return body
}
// func authCodeRequestHandler(w http.ResponseWriter, r *http.Request) {
//     code := retrieveAuthCode();  
//     fmt.Fprintf(w, code)
// }

// func authCodeResponseHandler(w http.ResponseWriter, r *http.Request){
//   fmt.Fprintf(w, "Query: %s",r.URL.Path)
// }

func main() {
    auth := retrieveAuthToken()
    fmt.Println(string(auth))
    err := ioutil.WriteFile("../token.json", auth, 0644)
    if(err != nil){
      panic(err)
    }
    // http.ListenAndServe(":3000", nil)
}