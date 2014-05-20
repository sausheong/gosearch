package main

import "flag"
import "net/url"
// import "fmt"
// import "strings"

var do_setup = flag.Bool("setup", false, "Run setup for GoSearch")
var spiders = 1

func main() {
  flag.Parse()
  if *do_setup {
    Setup()
  }
}

func process() {
  
}

func index(url string) (words []string, title string, err error) {
  
  return
}

func add_to_index() {
  
}

func scrub(link string) (scrubbed string, err error){
  u, err := url.Parse(link)
  scrubbed = u.String()
  return
}