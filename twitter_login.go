package main

/*
This snippet is written in GoLang and posts a comment on your twitter timeline. Includes: login, post comment, logout
fb_login_m.exe test@gmail.com password this text will posted
*/

import (
    "fmt"
    "net/http"
    "net/url"
    "os"
    "strings"
    "io/ioutil"
    "encoding/json"
    "net/http/cookiejar"
    //"html"
    "net/http/httputil"
    //"strconv"
    "flag"
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

type FBEngine struct {
    client *http.Client
}

func NewFBEngine() *FBEngine {
    jar, _ := cookiejar.New(nil)
    return &FBEngine{client: &http.Client{nil, nil, jar}}
}

func generate_post_data(email string, pass string, authenticity_token string) string {
    data, err := json.Marshal(map[string]interface{} {
        "session[username_or_email]": email,
        "session[password]": pass,
        "remember_me": "1",
        "return_to_ssl": "true",
        "scribe_log": "",
        "redirect_after_login": "/",
        "authenticity_token": authenticity_token,
        /*"username": email,
        "password:": pass,
        "commit": "Sign in",
        "authenticity_token": authenticity_token,
        */
    })
    if err != nil {
        fmt.Println(err)
    }
 
    return string(data)
}

func twitter_login(fbengine *FBEngine, email string, pass string) (string, error) {
    const use_new_http_request = 1

    ///////////////////
    /*fmt.Println("\nSend First request")
    resp, err := fbengine.client.Get("https://twitter.com/login") //https://twitter.com")
    defer resp.Body.Close()
    if err != nil {
        return "", "", "", fmt.Errorf("Get request failed: %s", err)
    }

    // should be: redirect_url := resp.Request.URL.String()
    redirect_url, _ := url.QueryUnescape(resp.Request.URL.String())

    // Write to file
    b, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", "", "", fmt.Errorf("Read HTML body failed: %s", err)
    }
    s := string(b)
    err = write_to_file("output1.html", s)
    if err != nil {
        return "", "", "", fmt.Errorf("Write file failed: %s", err)
    }

    // Print cookies
    for _, c := range resp.Cookies() { fmt.Println(c) }
   
    // extract some data for the post request
    //redirect_url = "https://twitter.com/sessions"
    redirect_url, _ = get_data(s, "<form action=\"https", "\"")
    redirect_url = "https" + redirect_url
    //authenticity_token, _ := get_data(s, "name=\"authenticity_token\" type=\"hidden\" value=\"", "\"")
    authenticity_token, _ := get_data(s, "name=\"authenticity_token\" value=\"", "\">")

    fmt.Println(email)
    fmt.Println(pass)
    fmt.Println(redirect_url)
    fmt.Println(authenticity_token)
    */
    req, err := http.NewRequest("GET", "https://twitter.com/login", nil)
    if err != nil {
        return "", fmt.Errorf("Post request failed: %s", err)
    }
    resp, err := fbengine.client.Do(req)
    
    if err != nil {
        return "", fmt.Errorf("Post request failed: %s", err)
    }
                
    defer resp.Body.Close()
            
    // should be: redirect_url := resp.Request.URL.String()
    redirect_url, _ := url.QueryUnescape(resp.Request.URL.String())
 
    // Write to file
    b, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("Read HTML body failed: %s", err)
    }
    s := string(b)
    err = write_to_file("output1.html", s)
    if err != nil {
        return "", fmt.Errorf("Write file failed: %s", err)
    }
    
    redirect_url, _ = get_data(s, "<form action=\"https", "\"")
    redirect_url = "https" + redirect_url
    authenticity_token, _ := get_data(s, "name=\"authenticity_token\" value=\"", "\">")
    
    // Print cookies
    for _, c := range resp.Cookies() { fmt.Println(c) }
    
    fmt.Println(redirect_url)
    fmt.Println(authenticity_token)
    
    ///////////////////
    fmt.Println("\nSend second request")
    if use_new_http_request == 0 {
        resp, err = fbengine.client.PostForm(redirect_url,
        url.Values{"session[username_or_email]": {email},
                   "session[password]": {pass},
                   "remember_me": {"1"},
                   "return_to_ssl": {"true"},
                   "scribe_log": {""},
                   "redirect_after_login": {"/"},
                   "authenticity_token": {authenticity_token}})
        /*url.Values{"username": {email},
                   "password:": {pass},
                   "commit": {"Sign in"},
                   "authenticity_token": {authenticity_token}})
                   */
    } else if use_new_http_request == 1 {
        // get post data
        //data := generate_post_data(email, pass, authenticity_token)

        data := url.Values{"session[username_or_email]": {email},
                   "session[password]": {pass},
                   "remember_me": {"1"},
                   "return_to_ssl": {"true"},
                   "scribe_log": {""},
                   "redirect_after_login": {"/"},
                   "authenticity_token": {authenticity_token}}
                   
        // prepare post request
        //bytes.NewBufferString(data.Encode()))
        req, err := http.NewRequest("POST", redirect_url, strings.NewReader(data.Encode())); //strings.NewReader(data()))
        if err != nil {
            return "", fmt.Errorf("Post request failed: %s", err)
        }

        // add cookies to the request which are in jar cookie
        //for _, cookie := range fbengine.client.Cookies() { //req.URL) {
        //    fmt.Println(cookie)
        //req.AddCookie(cookie)
        //}

        req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
        req.Header.Set("Accept-Encoding", "gzip,deflate,sdch")
        req.Header.Set("Accept-Language", "en-US,en;q=0.8")
        req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
        //req.Header.Set("Connection", "Connection:keep-alive")
        //req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
        //req.Header.Set("Referer", "https://twitter.com/")
        req.Header.Set("Host", "twitter.com")
        req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:16.0) Gecko/20100101 Firefox/16.0")
        req.Header.Set("Origin", "https://twitter.com")
        req.Header.Set("Referer", "https://twitter.com/")
        req.Header.Set("Cache-Control", "max-age=0")
        //req.Header.Set("Content-Length", strconv.Itoa(len(data)))
        //req.SetBasicAuth("u", "p")
        
        //dumpHead, _ := httputil.DumpRequest(req, false)
        dumpBody, _ := httputil.DumpRequest(req, true)

        //fmt.Println(string(dumpHead))
        fmt.Println(string(dumpBody))
        
        resp, err = fbengine.client.Do(req)
    }
       
    if err != nil {
        return "", fmt.Errorf("Post request failed: %s", err)
    }
                
    defer resp.Body.Close()
            
    // should be: redirect_url := resp.Request.URL.String()
    redirect_url, _ = url.QueryUnescape(resp.Request.URL.String())
 
    // Write to file
    b, err = ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("Read HTML body failed: %s", err)
    }
    s = string(b)
    err = write_to_file("output2.html", s)
    if err != nil {
        return "", fmt.Errorf("Write file failed: %s", err)
    }

    // print cookies
    fmt.Println("cookies:")
    for _, c := range resp.Cookies() { fmt.Println(c) }

    fmt.Println(resp.Header)
    
    // extract some data for the post request
    /*status_update_url, _ := get_data(s, "<form method=\"post\" class=\"composer_form\" id=\"composer_form\" action=\"", "\"")
    status_update_url = "https://m.facebook.com" + status_update_url
    fb_dtsg, _ := get_data(s, "name=\"fb_dtsg\" value=\"", "\"")
    privacy, _ := get_data(s, "name=\"privacy\" value=\"", ";\"")
    privacy = html.UnescapeString(privacy)
    */
    
    return authenticity_token, fmt.Errorf("err")
}

func twitter_geo_locate(fbengine *FBEngine, authenticity_token string) (error) {
      
    data := url.Values{"authenticity_token": {authenticity_token}}
                   
    // prepare post request
    redirect_url := "https://twitter.com/account/geo_locate"
    req, err := http.NewRequest("POST", redirect_url, strings.NewReader(data.Encode()));
    if err != nil {
        return fmt.Errorf("Post request failed: %s", err)
    }

    req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
    req.Header.Set("Accept-Encoding", "gzip,deflate,sdch")
    req.Header.Set("Accept-Language", "en-US,en;q=0.8")
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
    req.Header.Set("Host", "twitter.com")
    req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:16.0) Gecko/20100101 Firefox/16.0")
    req.Header.Set("Origin", "https://twitter.com")
    req.Header.Set("Referer", "https://twitter.com/")
    req.Header.Set("Cache-Control", "max-age=0")
    req.Header.Set("X-Requested-With", "XMLHttpRequest")
    //req.Header.Set("Content-Length", strconv.Itoa(len(data)))
    //req.SetBasicAuth("u", "p")
        
    //dumpHead, _ := httputil.DumpRequest(req, false)
    dumpBody, _ := httputil.DumpRequest(req, true)

    //fmt.Println(string(dumpHead))
    fmt.Println(string(dumpBody))
        
    resp, err := fbengine.client.Do(req)
        
    if err != nil {
        return fmt.Errorf("Post request failed: %s", err)
    }
                
    defer resp.Body.Close()
            
    // should be: redirect_url := resp.Request.URL.String()
    redirect_url, _ = url.QueryUnescape(resp.Request.URL.String())
 
    // Write to file
    b, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("Read HTML body failed: %s", err)
    }
    s := string(b)
    err = write_to_file("output3.html", s)
    if err != nil {
        return fmt.Errorf("Write file failed: %s", err)
    }

    // print cookies
    fmt.Println("cookies:")
    for _, c := range resp.Cookies() { fmt.Println(c) }

    fmt.Println(resp.Header)
    
    return nil
}

func twitter_post_comment(fbengine *FBEngine, authenticity_token string, comment []string) (string, error) {   
    var tweet string
    for _, c := range comment { tweet += c + " "}
    
    // place id from:
    // http://nominatim.openstreetmap.org/search?q=800%20Ocean%20Drive,%20Miami%20Beach,%20USA&format=xml  
    data := url.Values{"status": {tweet},
                       "place_id": {"df51dec6f4ee2b2c"}, //4b58830723ec6371"}, //445781941346779136
                       "authenticity_token": {authenticity_token}}
                   
    // prepare post request
    redirect_url := "https://twitter.com/i/tweet/create"
    req, err := http.NewRequest("POST", redirect_url, strings.NewReader(data.Encode())); //strings.NewReader(data()))
    if err != nil {
        return "", fmt.Errorf("Post request failed: %s", err)
    }

    req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
    req.Header.Set("Accept-Encoding", "gzip,deflate,sdch")
    req.Header.Set("Accept-Language", "en-US,en;q=0.8")
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
    req.Header.Set("Host", "twitter.com")
    req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:16.0) Gecko/20100101 Firefox/16.0")
    req.Header.Set("Origin", "https://twitter.com")
    req.Header.Set("Referer", "https://twitter.com/")
    req.Header.Set("Cache-Control", "max-age=0")
    req.Header.Set("X-Requested-With", "XMLHttpRequest")
    //req.Header.Set("Content-Length", strconv.Itoa(len(data)))
    //req.SetBasicAuth("u", "p")
        
    //dumpHead, _ := httputil.DumpRequest(req, false)
    dumpBody, _ := httputil.DumpRequest(req, true)

    //fmt.Println(string(dumpHead))
    fmt.Println(string(dumpBody))
        
    resp, err := fbengine.client.Do(req)
        
    if err != nil {
        return "", fmt.Errorf("Post request failed: %s", err)
    }
                
    defer resp.Body.Close()
            
    // should be: redirect_url := resp.Request.URL.String()
    redirect_url, _ = url.QueryUnescape(resp.Request.URL.String())
 
    // Write to file
    b, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("Read HTML body failed: %s", err)
    }
    s := string(b)
    err = write_to_file("output4.html", s)
    if err != nil {
        return "", fmt.Errorf("Write file failed: %s", err)
    }

    // print cookies
    fmt.Println("cookies:")
    for _, c := range resp.Cookies() { fmt.Println(c) }

    fmt.Println(resp.Header)
    
    // extract some data for the get request
    logout_url, _ := get_data(s, "id=\"signout-form\" action=\"/", "\"")
    logout_url = "https://twitter.com/" + logout_url

    return logout_url, nil
}

func twitter_logout(fbengine *FBEngine, authenticity_token string, logout_url string) (error) {     
    data := url.Values{"reliability_event": {},
                       "scribe_log": {},
                       "authenticity_token": {authenticity_token}}
                   
    // prepare post request
    req, err := http.NewRequest("POST", logout_url, strings.NewReader(data.Encode()));
    if err != nil {
        return fmt.Errorf("Post request failed: %s", err)
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
    //req.Header.Set("Content-Length", strconv.Itoa(len(data)))
    //req.SetBasicAuth("u", "p")
        
    //dumpHead, _ := httputil.DumpRequest(req, false)
    dumpBody, _ := httputil.DumpRequest(req, true)

    //fmt.Println(string(dumpHead))
    fmt.Println(string(dumpBody))
        
    resp, err := fbengine.client.Do(req)
        
    if err != nil {
        return fmt.Errorf("Post request failed: %s", err)
    }
                
    defer resp.Body.Close()
            
    // should be: redirect_url := resp.Request.URL.String()
    //redirect_url, _ := url.QueryUnescape(resp.Request.URL.String())
 
    // Write to file
    b, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("Read HTML body failed: %s", err)
    }
    s := string(b)
    err = write_to_file("output4.html", s)
    if err != nil {
        return fmt.Errorf("Write file failed: %s", err)
    }

    // print cookies
    fmt.Println("cookies:")
    for _, c := range resp.Cookies() { fmt.Println(c) }

    fmt.Println(resp.Header)
    
    return fmt.Errorf("err")
}

func main() {
    flag.Parse()
    args := flag.Args()
    if len(args) < 3 {
        fmt.Println("Please pass arguments: Email, Password, Comment. e.g. twitter.exe test@gmail.com password text to post")
        return
    }

    // here the login infos
    var email = args[0]
    var pass = args[1]
    var text = args[2:]
        
    var fbengine FBEngine = *NewFBEngine()
 
    authenticity_token, _ := twitter_login(&fbengine, email, pass)
    twitter_geo_locate(&fbengine, authenticity_token)
    logout_url, _ := twitter_post_comment(&fbengine, authenticity_token, text)
    twitter_logout(&fbengine, authenticity_token, logout_url)
}	
