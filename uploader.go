package main

import (
  "net/http"
  "os"
  "io"
  "io/ioutil"
  "log"
  "bytes"
  "mime/multipart"
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

  var requestBody bytes.Buffer

  multiPartWriter := multipart.NewWriter(&requestBody)

  fileWriter, err := multiPartWriter.CreateFormFile("upload", filepath.Base(fileName))

  if err != nil {
    log.Fatalln(err)
  }

  _, err = io.Copy(fileWriter, file)
  if err != nil {
    log.Fatalln(err)
  }

  multiPartWriter.Close()

  req, err := http.NewRequest("POST", config.Url + "/" + filepath.Base(fileName) + "?type=UPLOAD_FILE", &requestBody)
  if err != nil {
    log.Fatalln(err)
  }
  req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())

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
