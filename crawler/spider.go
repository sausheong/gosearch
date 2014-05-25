package main

import "os"
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
var do_force = flag.Bool("force", false, "Force to revisit all indexed pages from seed file")
var spiders = 1
var in_progress = fifo.NewQueue()

func main() {
  flag.Parse()
  if *do_setup {
    Setup()
  }
  
  // process now
  // load urls into the in_progress queue
  b, err := ioutil.ReadFile("seed.txt")
  if err != nil { 
   println("Error is %v", err) 
  }
  for _, line := range strings.Split(string(b), "\n") {
    in_progress.Add(line)
  }    

  // index the urls in progress queue
  go index()
  
  var input string
  fmt.Scanln(&input)
}

// indexes a page
func index() {
  for {
    item := in_progress.Next()
  
    if item != nil {
      u := item.(string)
      page := Page{Url: u}
      DB.Where(page).FirstOrCreate(&page)
      println("Indexing page ", page.Url)
      if *do_force || (time.Since(page.CreatedAt).Seconds() < 1) || (time.Since(page.UpdatedAt).Hours() > 24) {
        page.UpdatedAt = time.Now()
        DB.Save(&page)
    
        println("- OK")
        words := words_from(u)
        for i := 0; i < len(words); i++ {
          w := words[i]
          word := Word{Stem: w}
          DB.Where(word).FirstOrCreate(&word)
          loc := Location{Position: int64(i), WordId: word.Id, PageId: page.Id}
          DB.Where(loc).FirstOrCreate(&loc)
        }       
        extracted_links := links_from(u)
        println("- no of links found -> ", len(extracted_links))
        for _, link := range extracted_links {
          in_progress.Add(link)
        }  
      } else {
        println(" - Already indexed")
      }
    } else {
      println("No more items - exiting crawler")
      os.Exit(0)
    }
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

// Find links from a given URL
func links_from(url string) (links []string) { 
  resp, err := http.Get(url)
  if err != nil {
    fmt.Println("Error in getting URL - ", err)
    return
  }
  defer resp.Body.Close()
  
  doc, err := html.Parse(resp.Body)
  if err != nil {
    fmt.Println("Error in parsing html - ", err)
    return
  }  
  find_links(url, doc, &links)
  return  
}

// Iterative function to find links, given a node
func find_links(parent string, n *html.Node, links *[]string) {
  parent_url, _ := url.Parse(parent)
  if n.Type == html.ElementNode && n.Data == "a" {
    for _, a := range n.Attr {
      if a.Key == "href" {
        link, err := url.Parse(a.Val)
        if err != nil {
          fmt.Println("Error while finding link - ", err)
        } else {
          link = parent_url.ResolveReference(link)    
          if !ignored_link(link.Path) {
            *links = append(*links, link.String())
          }
        }
      }
    }
  }
  for c := n.FirstChild; c != nil; c = c.NextSibling {
    find_links(parent, c, links)
  }
}

// Check if the word is a stopword and should therefore be ignored

func ignore(word string) (ignored bool) {
  ignored = Stopwords[strings.ToLower(word)]
  return
}

func ignored_link(link string) (ignored bool) {
  ignored = strings.HasSuffix(link, "jpg") || 
  strings.HasSuffix(link, "gif") || 
  strings.HasSuffix(link, "png") ||   
  strings.HasSuffix(link, "pdf") 
  return
}