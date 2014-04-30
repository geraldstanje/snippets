package main

/*
This snippet is written in GoLang and posts a comment on your Twitter timeline. Includes: login, post comment, logout
./twitter test@gmail.com password this text will be tweeted
*/

import (
  "fmt"
  "net/http"
  "net/url"
  "os"
  "strings"
  "io/ioutil"
  "net/http/cookiejar"
  "flag"
  //"encoding/json"
  //"html"
  //"net/http/httputil"
  //"strconv"
)

func get_data(s string, start_str string, end_str string) (string, error) {
  var data string
 
  i_start := strings.Index(s, start_str)
  if i_start == -1 {
    return "", fmt.Errorf("start string not found")
  }
 
  s_new := s[i_start + len(start_str):]
     
  i_end := strings.Index(s_new, end_str)
  if i_end == -1 {
    return "", fmt.Errorf("end string not found")
  }
 
  data = s[i_start + len(start_str) : i_start + len(start_str) + i_end]
    
  return data, nil
}

func write_to_file(file_str string, s string) error {
  file, err := os.Create(file_str)

  if err == nil {
    file.WriteString(s)
    defer file.Close()
    return nil
  }

  return err
}

type TwitterEngine struct {
  client *http.Client
}

func NewTwitterEngine() *TwitterEngine {
  jar, _ := cookiejar.New(nil)
  return &TwitterEngine{client: &http.Client{nil, nil, jar}}
}

func send_http_request(TwitterEngine *TwitterEngine, urlstr string, send_post_data bool, post_data url.Values) (string, string, error) {
  var req *http.Request
  var err error

  if send_post_data == false {
    req, err = http.NewRequest("GET", urlstr, nil)
    if err != nil {
        return "", "", fmt.Errorf("Get request failed: %s", err)
    }
  } else {
    req, err = http.NewRequest("POST", urlstr, strings.NewReader(post_data.Encode()))
    if err != nil {
       return "", "", fmt.Errorf("Post request failed: %s", err)
    }

    req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
    req.Header.Set("Accept-Encoding", "gzip,deflate,sdch")
    req.Header.Set("Accept-Language", "en-US,en;q=0.8")
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("Host", "twitter.com")
    req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:16.0) Gecko/20100101 Firefox/16.0")
    req.Header.Set("Origin", "https://twitter.com")
    req.Header.Set("Referer", "https://twitter.com/")
    req.Header.Set("Cache-Control", "max-age=0")
  }

  resp, err := TwitterEngine.client.Do(req)
  if err != nil {
    return "", "", fmt.Errorf("Http request failed: %s", err)
  }
                
  defer resp.Body.Close()
            
  // should be: redirect_url := resp.Request.URL.String()
  redirect_url, _ := url.QueryUnescape(resp.Request.URL.String())

  // Read HTML body
  b, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return "", "", fmt.Errorf("Read HTML body failed: %s", err)
  }
  str := string(b)
    
  // print cookies
  fmt.Println("cookies:")
  for _, c := range resp.Cookies() { fmt.Println(c) }

  return str, redirect_url, nil
}

func twitter_login(TwitterEngine *TwitterEngine, email string, pass string) (string, error) {
  var err error
  var s string
  var redirect_url string

  s, redirect_url, _ = send_http_request(TwitterEngine, "https://twitter.com/login", false, nil)

  err = write_to_file("output1.html", s)
  if err != nil {
    return "", fmt.Errorf("Write file failed: %s", err)
  }
    
  redirect_url, _ = get_data(s, "<form action=\"https", "\"")
  redirect_url = "https" + redirect_url
  authenticity_token, _ := get_data(s, "name=\"authenticity_token\" value=\"", "\">")
    
  data := url.Values{"session[username_or_email]": {email},
                     "session[password]": {pass},
                     "remember_me": {"1"},
                     "return_to_ssl": {"true"},
                     "scribe_log": {""},
                     "redirect_after_login": {"/"},
                     "authenticity_token": {authenticity_token}}
    
  s, redirect_url, _ = send_http_request(TwitterEngine, redirect_url, true, data)

  err = write_to_file("output2.html", s)
  if err != nil {
    return "", fmt.Errorf("Write file failed: %s", err)
  }
    
  return authenticity_token, nil
}

func twitter_geo_locate(TwitterEngine *TwitterEngine, authenticity_token string) (error) {
  var err error
  var s string

  data := url.Values{"authenticity_token": {authenticity_token}}
  s, _, _ = send_http_request(TwitterEngine, "https://twitter.com/account/geo_locate", true, data)

  err = write_to_file("output3.html", s)
  if err != nil {
    return fmt.Errorf("Write file failed: %s", err)
  }
    
  return nil
}

func twitter_post_comment(TwitterEngine *TwitterEngine, authenticity_token string, comment []string) (string, error) {   
  var err error
  var s string
  var tweet string
  for _, c := range comment { tweet += c + " "}
    
  // place id from:
  // http://nominatim.openstreetmap.org/search?q=800%20Ocean%20Drive,%20Miami%20Beach,%20USA&format=xml  
  data := url.Values{"status": {tweet},
                     "place_id": {"df51dec6f4ee2b2c"}, //4b58830723ec6371"}, //445781941346779136
                     "authenticity_token": {authenticity_token}}
        
  s, _, _ = send_http_request(TwitterEngine, "https://twitter.com/i/tweet/create", true, data)

  err = write_to_file("output4.html", s)
  if err != nil {
    return "", fmt.Errorf("Write file failed: %s", err)
  }
    
  // extract some data for the get request
  logout_url, _ := get_data(s, "id=\"signout-form\" action=\"/", "\"")
  logout_url = "https://twitter.com/" + logout_url

  return logout_url, nil
}

func twitter_logout(TwitterEngine *TwitterEngine, authenticity_token string, logout_url string) (error) { 
  var err error
  var s string

  data := url.Values{"reliability_event": {},
                     "scribe_log": {},
                     "authenticity_token": {authenticity_token}}
                   
  s, _, _ = send_http_request(TwitterEngine, logout_url, true, data)
  
  err = write_to_file("output4.html", s)
  if err != nil {
    return fmt.Errorf("Write file failed: %s", err)
  }

  return nil
}

func main() {
  flag.Parse()

  args := flag.Args()
  
  if len(args) < 3 {
    fmt.Println("Please pass arguments: Email, Password, Comment. e.g. ./twitter test@gmail.com password this text will be tweeted")
    return
  }

  // here the login infos
  var email = args[0]
  var pass = args[1]
  var text = args[2:]
        
  var TwitterEngine TwitterEngine = *NewTwitterEngine()
 
  authenticity_token, _ := twitter_login(&TwitterEngine, email, pass)
  twitter_geo_locate(&TwitterEngine, authenticity_token)
  logout_url, _ := twitter_post_comment(&TwitterEngine, authenticity_token, text)
  twitter_logout(&TwitterEngine, authenticity_token, logout_url)
}