package main

import "flag"
import "net/url"
import "strings"
import "net/http"
import "regexp"
import "time"
import "io/ioutil"
import "fmt"

import "github.com/reiver/go-porterstemmer"
import "code.google.com/p/go.net/html"
import "github.com/advancedlogic/GoOse"
import "github.com/foize/go.fifo"

var do_setup = flag.Bool("setup", false, "Run setup for GoSearch")
var spiders = 1
var in_progress = fifo.NewQueue()

func main() {
  flag.Parse()
  if *do_setup {
    Setup()
  }
  
  // process now
  // load urls into the in_progress queue
  b, err := ioutil.ReadFile("test_urls.txt")
  if err != nil { 
   println("Error is %v", err) 
  }
  for _, line := range strings.Split(string(b), "\n") {
    in_progress.Add(line)
  }    

  // index the in progress queue
  for {    
    index()
  }
}

// indexes a page
func index() {
  url := in_progress.Next().(string)
  scrubbed_url, err := scrub(url)
  if err != nil {
    fmt.Println("Error while indexing is %v", err)
  }
  page := Page{Url: scrubbed_url}
  DB.Where(page).FirstOrCreate(&page)
  println("Indexing page ", page.Url)
  if time.Since(page.UpdatedAt).Hours() < 24 {
    println("Indexing now")
    words := words_from(url)
    for i := 0; i < len(words); i++ {
      w := words[i]
      word := Word{Stem: w}
      DB.Where(word).FirstOrCreate(&word)
      loc := Location{Position: int64(i), WordId: word.Id, PageId: page.Id}
      DB.Where(loc).FirstOrCreate(&loc)
    }       
  } else {
    println(" - Page already indexed")
  }
  extracted_links := links_from(url)
  for _, link := range extracted_links {
    in_progress.Add(link)
  }  
  return
}

// Get words from a given link, returning an array of strings
// words are set to lower case, checked for stop words and stemmed
func words_from(link string) (words []string) {
  g := goose.New()
  article := g.ExtractFromUrl(link)
  var text string = article.CleanedText

  str := strings.TrimSpace(text)
  split_words := strings.Split(str, " ")
  
  r, _ := regexp.Compile("[^A-Za-z]")   
  for _, val := range split_words {
    w := r.ReplaceAllString(val, "")
    if !ignore(w) {
      w = porterstemmer.StemString(w)
      words = append(words, w)
    }    
  }  
  return
}

// Scrub a link for uniformity
func scrub(link string) (scrubbed string, err error){
  u, err := url.Parse(link)
  scrubbed = u.String()
  return
}

// Find links from a given URL
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

// Iterative function to find links, given a node
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

// Check if the word is a stopword and should therefore be ignored

func ignore(word string) (ignored bool) {
  ignored = Stopwords[strings.ToLower(word)]
  return
}