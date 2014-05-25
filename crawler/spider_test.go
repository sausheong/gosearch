package main

import "testing"
// import "io/ioutil"
// import "strings"


func TestMain(t *testing.T) {
  main()
}

// func Test_Scrub(t *testing.T) {
//   b, err := ioutil.ReadFile("test_urls.txt")
//   if err != nil { 
//    t.Errorf("Error is %v", err) 
//   }
//   for _, line := range strings.Split(string(b), "\n") {
//     _, err := scrub(line)
//     if err != nil {
//       t.Errorf("Error is %v", err) 
//     }
//   }  
// }
// 
// func Test_Links(t *testing.T) {
//   link := "https://sg.news.yahoo.com/red-faces-french-trains-too-wide-stations-144918613.html"
//   extracted := links_from(link)
//   println(len(extracted))
//   for _, val := range extracted {
//     println(val)
//   }
// }
// 
// func Test_Index(t *testing.T) {
//   link := "https://sg.news.yahoo.com/red-faces-french-trains-too-wide-stations-144918613.html"
//   err := index(link)
//   if err != nil { 
//    t.Errorf("Error is %v", err) 
//   }
//   
// }
// 
// func Test_WordsFrom(t *testing.T) {
//   words := words_from("https://sg.news.yahoo.com/red-faces-french-trains-too-wide-stations-144918613.html")
//   for _, val := range words {
//    println(val) 
//   }
//   
// }

func Test_IgnoreStopWord(t *testing.T) {
  if ignore("a") == false {
    t.Errorf("a is not ignored") 
  }
  if ignore("nothing") == false {
    t.Errorf("nothing is not ignored") 
  }
}
