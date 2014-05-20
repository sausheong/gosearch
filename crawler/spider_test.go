package main

import "testing"
import "io/ioutil"
import "strings"
import "fmt"

func Test_Scrub(t *testing.T) {
  b, err := ioutil.ReadFile("test_urls.txt")
  if err != nil { 
   t.Errorf("Error is %v", err) 
  }
  for _, line := range strings.Split(string(b), "\n") {
    _, err := scrub(line)
    if err != nil {
      t.Errorf("Error is %v", err) 
    }
  }  
}

func Test_Links(t *testing.T) {
  link := "https://sg.yahoo.com/?p=us"
  extracted := links_from(link)
  fmt.Println(extracted)
}

func Test_Index(t *testing.T) {
  link := "https://sg.yahoo.com/?p=us"
  words, title, err := index(link)
  if err != nil { 
   t.Errorf("Error is %v", err) 
  }
  fmt.Println(words)
  fmt.Println(title)
  
}