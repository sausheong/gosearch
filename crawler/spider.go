package main

import "flag"
import "net/url"
import "fmt"
import "strings"
import "net/http"
// import "io/ioutil"
import "code.google.com/p/go.net/html"
import "github.com/advancedlogic/GoOse"
import "regexp"

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
  resp, err := http.Get(url)
  if err != nil {
    fmt.Println("Error is %v", err)
  }
  defer resp.Body.Close()
  
  doc, err := html.Parse(resp.Body)
  if err != nil {
    fmt.Println("Error is %v", err)
  }  

  fmt.Println("doc", doc)
  return
}

func add_to_index() {

}

func words_from(link string) (words []string) {
  g := goose.New()
  article := g.ExtractFromUrl(link)
  var text string = article.CleanedText
  str := strings.TrimSpace(text)
  split_words := strings.Split(str, " ")
  
  r, _ := regexp.Compile("[^A-Za-z]")   
  for _, val := range split_words {
    w := r.ReplaceAllString(val, "")
    words = append(words, w)
  }  
  return
}

func scrub(link string) (scrubbed string, err error){
  u, err := url.Parse(link)
  scrubbed = u.String()
  return
}


func links_from(url string) (links []string) {
  resp, err := http.Get(url)
  if err != nil {
    fmt.Println("Error is %v", err)
  }
  defer resp.Body.Close()
  
  doc, err := html.Parse(resp.Body)
  if err != nil {
    fmt.Println("Error is %v", err)
  }  
  find_links(doc, &links)
  return  
}

func find_links(n *html.Node, links *[]string) {
  if n.Type == html.ElementNode && n.Data == "a" {
    for _, a := range n.Attr {
      if a.Key == "href" {
        *links = append(*links, a.Val)
        break
      }
    }
  }
  for c := n.FirstChild; c != nil; c = c.NextSibling {
    find_links(c, links)
  }
}

