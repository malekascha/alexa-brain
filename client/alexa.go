package main

import (
    "fmt"
    "net/http"
    // "net/url"
    "io/ioutil"
    "encoding/json"
    "os"
    "os/exec"
)

type AuthResponse struct {
  Token_type string `json:"token_type"`
  Refresh_token string `json:"refresh_token"`
  Access_token string `json:"access_token"`
  Expires_in int `json:"expires_in"`
}

type Context struct {
  context []interface{}
  event map[string]interface{}
}

type ContextHeader struct {
  namespace string
  name string
}

type AudioPlayerPayload struct {
  token string
  offsetInMilliseconds string
  playerActivity string
}

type Alert struct {
  token string
  type string
  scheduledTime string
}

type AlertsPayload struct {
  allAlerts []Alert
  activeAlerts []Alert
}

type SpeakerPayload struct {
  volume int
  muted bool
}

type SpeechSynthesizerPayload struct {
  token string
  offsetInMilliseconds string
  playerActivity string
}

func generateUUID() []uint8 {
  id, err := exec.Command("uuidgen").Output()
  if(err != nil) {
    panic(err)
  }
  return id
}

func setAuthEnvVars() {
  data, err := ioutil.ReadFile("../token.json")
  if(err != nil){
    panic(err)
  }
  str := string(data[:])
  err = os.Setenv("auth", str)
  if(err != nil){
    panic(err)
  }
}

func fetchAccessToken() string {
  auth_info_string := os.Getenv("auth")
  auth_info := []byte(auth_info_string)
  var auth_response AuthResponse
  err := json.Unmarshal(auth_info, &auth_response)
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
  fmt.Println(contents)
  if(err != nil){
    panic(err)
  }
  
}

// func initDeviceContext() string {
//   context := {
//     ""
//   }
// }

func synchronizeInitialState() {
  id := generateUUID()
  context := Context{
    {
     ContextHeader{
      "Audioplayer",
      "PlaybackState"
     } ,
     interface{}
    },
    {
      ContextHeader{
        "Alerts",
        "AlertsState"
      },
      AlertsPayload{
        {},
        {}
      }
    },
    {
      ContextHeader{
        "Speaker",
        "VolumeState"
      },
      
    }
  }
}

func main() {
  setAuthEnvVars()
  initDownchannel()
}