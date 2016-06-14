package main

import (
    "fmt"
    "net/http"
    // "net/url"
    "io/ioutil"
    "encoding/json"
    "os"
    "os/exec"
    "strconv"
)

//STRUCTS//////////////////////////////////////////////////////////

type AuthResponse struct {
  Token_type string `json:"token_type"`
  Refresh_token string `json:"refresh_token"`
  Access_token string `json:"access_token"`
  Expires_in int `json:"expires_in"`
}

type EventRequestBody struct {
  context []interface{} `json:"context"`
  event interface{} `json:"event"`
}

type Event struct {
  header interface{} `json:"header"`
  payload interface{} `json:"payload"`
}

type EventHeaders struct {
  namespace string `json:"namespace"`
  name string `json:"name"`
  messageId string `json:"messageId"`
}

type ContextHeader struct {
  namespace string `json:"namespace"`
  name string `json:"name"`
}

type AudioPlayerPayload struct {
  token string `json:"token"`
  offsetInMilliseconds string `json:"offsetInMilliseconds"`
  playerActivity string `json:"playerActivity"`
}

type Alert struct {
  token string `json:"token"`
  // type string
  scheduledTime string `json:"scheduledTime"`
}

type AlertsPayload struct {
  allAlerts []Alert `json:"allAlerts"`
  activeAlerts []Alert `json:"activeAlerts"`
}

type SpeakerPayload struct {
  volume int `json:"volume"`
  muted bool `json:"muted"`
}

type SpeechSynthesizerPayload struct {
  token string `json:"token"`
  offsetInMilliseconds string `json:"offsetInMilliseconds"`
  playerActivity string `json:"playerActivity"`
}


//AUTH FUNCTIONS/////////////////////////////////////////////

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

//UTILITY FUNCTIONS/////////////////////////////////////////////////

func generateUUID() string {
  id, err := exec.Command("uuidgen").Output()
  if(err != nil) {
    panic(err)
  }
  return string(id)
}

func getVolume() int {
  out, err := exec.Command("/bin/sh", "./scripts/getVolume.sh").Output()
  if(err != nil) {
    panic(err)
  }
  vol := string(out)[0:2]
  percent, err := strconv.Atoi(vol)
  if(err != nil) {
    panic(err)
  }
  return percent
}

//ALEXA API CALLS////////////////////////////////////////////////////

func initDownchannel() {
  api_endpoint := "https://avs-alexa-na.amazon.com/v20160207/directives"
  access_token := fetchAccessToken()
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
  fmt.Println(string(contents))
  if(err != nil){
    panic(err)
  }
  
}

func synchronizeInitialState() {
  id := generateUUID()
  volume := getVolume()
  audioHeader := ContextHeader{"Audioplayer", "PlaybackState"}
  alertsHeader := ContextHeader{"Alerts","AlertsState"}
  speakerHeader := ContextHeader{"Speaker","VolumeState"}
  speechSynthesizerHeader := ContextHeader{"SpeechSynthesizer","SpeechState"}
  audioContext := Event{audioHeader, AudioPlayerPayload{}}
  alertsContext := Event{alertsHeader, AlertsPayload{}}
  speakerContext := Event{speakerHeader, SpeakerPayload{volume,false}} //TODO: actually check mute status of system
  speechSynthesizerContext := Event{speechSynthesizerHeader, SpeechSynthesizerPayload{}}
  context := []interface{}{audioContext,alertsContext,speakerContext,speechSynthesizerContext}
  eventHeaders := EventHeaders{"System","SynchronizeState",id}
  event := Event{eventHeaders, nil}
  body := EventRequestBody{context,event}
  fmt.Println(body.context[1:])
  encoded_body, err := json.Marshal(body)
  if(err != nil){
    panic(err)
  }
  fmt.Println(string(encoded_body))
  // var m interface{}
  // _ = json.Unmarshal(encoded_body, m)
  // fmt.Println(m)
}

//MAIN///////////////////////////////////////////////

func main() {
  setAuthEnvVars()
  // initDownchannel()
  synchronizeInitialState()
}