package main

import (
  "net/http"
  "os"
  "io/ioutil"
  "log"
  "path/filepath"
  "gopkg.in/yaml.v2"
  "os/exec"
  "time"
)

type uploadConfig struct {
  Url string `yaml:"url"`
  Launch string `yaml:"launch"`
}

func main() {
  fileName := os.Args[1]

  file, err := os.Open(fileName)

  if err != nil {
    log.Fatalln(err)
  }
  defer file.Close()

  config := &uploadConfig{}

  curPath, err := filepath.Abs(filepath.Dir(os.Args[0]))

  if err != nil {
    log.Fatalln(err)
  }

  configFile, err := ioutil.ReadFile(filepath.Join(curPath, "config.yml"))

  if err != nil {
    log.Fatalln(err)
  }

  err = yaml.Unmarshal(configFile, &config)
  if err != nil {
    log.Fatalln(err)
  }

  if config.Url == "" {
    log.Fatalln("Url was not configured")
  }

  req, err := http.NewRequest("POST", config.Url + "/" + filepath.Base(fileName) + "?type=UPLOAD_FILE", file)
  if err != nil {
    log.Fatalln(err)
  }

  client := &http.Client{}
  _, err = client.Do(req)
  if err != nil {
    log.Fatalln(err)
  }

  if config.Launch != "" {
    exec.Command(config.Launch, os.Args[1]).Start()
  }
  // PRs welcome :)
  time.Sleep(5 * time.Second)
}
