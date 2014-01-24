package main

/* 
This snippet is written in GoLang and posts a comment on your Facebook pinwall. Includes: login, post comment, logout
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
	"html"
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

func generate_post_data(email string, m_ts string, li string, lsd string, charset_test string, signup_layout string, pass string) string {
    data, err := json.Marshal(map[string]interface{} {
        "email": email,
        "login": "Log in",  
        "m_ts": m_ts,
        "li": li,
        "lsd": lsd, 
        "charset_test": charset_test,
        "ajax": "0",
        "gps": "0",
        "pxr": "0",
        "width": "0",
        "version": "1",
        "signup_layout": signup_layout,
        "pass": pass,
    })
    if err != nil {
        fmt.Println(err)
    }
 
    return string(data)
}

func fb_login(fbengine *FBEngine, email string, pass string) (string, string, string, error) {
	const use_new_http_request = 0

	///////////////////
    fmt.Println("\nSend First request")
    resp, err := fbengine.client.Get("http://m.facebook.com/")
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
    redirect_url, _ = get_data(s, "action=\"", "\"")
    m_ts, _ := get_data(s, "name=\"mts\" value=\"", "\"")
    li, _ := get_data(s, "name=\"li\" value=\"", "\"")
    lsd, _ := get_data(s, "name=\"lsd\" value=\"", "\"")
    charset_test, _ := get_data(s, "name=\"charset_test\" value=\"", "\"")
    signup_layout, _ := get_data(s, "name=\"signup_layout\" value=\"", "\"")

	///////////////////
    fmt.Println("\nSend second request")
    if use_new_http_request == 0 {
		resp, err = fbengine.client.PostForm(redirect_url,
		url.Values{"email": {email},
				   "login": {"Log in"},  
				   "m_ts": {m_ts},
				   "li": {li},
				   "lsd": {lsd}, 
				   "charset_test": {charset_test},
				   "ajax": {"0"},
				   "gps": {"0"},
				   "pxr": {"0"},
				   "width": {"0"},
				   "version": {"1"},
				   "signup_layout": {signup_layout},
				   "pass": {pass}})
	} else if use_new_http_request == 1 {
		// get post data
		data := generate_post_data(email, m_ts, li, lsd, charset_test, signup_layout, pass)
		
		// prepare post request
		req, err := http.NewRequest("POST", redirect_url, strings.NewReader(data))
		if err != nil {
			return "", "", "", fmt.Errorf("Post request failed: %s", err)
		}

		// add cookies to the request which are in jar cookie
		//for _, cookie := range jar.Cookies(req.URL) {
		//	req.AddCookie(cookie)
		//}

		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:16.0) Gecko/20100101 Firefox/16.0")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
		req.Header.Set("Referer", "http://m.facebook.com/")
		req.Header.Set("Host", "m.facebook.com")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-gb,en;q=0.5")

		resp, err = fbengine.client.Do(req)
	}

    if err != nil {
		return "", "", "", fmt.Errorf("Post request failed: %s", err)
    }
	defer resp.Body.Close()

    // should be: redirect_url := resp.Request.URL.String()
    redirect_url, _ = url.QueryUnescape(resp.Request.URL.String())
 
    // Write to file
    b, err = ioutil.ReadAll(resp.Body)
    if err != nil {
		return "", "", "", fmt.Errorf("Read HTML body failed: %s", err)
    }
    s = string(b)
    err = write_to_file("output2.html", s)
    if err != nil {
		return "", "", "", fmt.Errorf("Write file failed: %s", err)
    }

	// print cookies
    for _, c := range resp.Cookies() { fmt.Println(c) }

	// extract some data for the post request
	status_update_url, _ := get_data(s, "<form method=\"post\" class=\"composer_form\" id=\"composer_form\" action=\"", "\"")
    status_update_url = "https://m.facebook.com" + status_update_url  
    fb_dtsg, _ := get_data(s, "name=\"fb_dtsg\" value=\"", "\"")
    privacy, _ := get_data(s, "name=\"privacy\" value=\"", ";\"")
	privacy = html.UnescapeString(privacy)

	return status_update_url, fb_dtsg, privacy, nil
}

func fb_post_comment(fbengine *FBEngine, comment []string, status_update_url string, fb_dtsg string, privacy string) (string, error) {
	var text_to_share string
	for _, c := range comment { text_to_share += c + " "}

	fmt.Println("\nSend third request")
	resp, err := fbengine.client.PostForm(status_update_url,
    url.Values{"fb_dtsg": {fb_dtsg},
               "update": {"Share"},
               "target": {""}, 
               "status": {text_to_share},
			   "privacy": {privacy}})
    if err != nil {
		return "", fmt.Errorf("Post request failed: %s", err)
    }
	defer resp.Body.Close()

    // Write to file
    b, err := ioutil.ReadAll(resp.Body)
    if err != nil {
		return "", fmt.Errorf("Post request failed: %s", err)
    }
    s := string(b)
    err = write_to_file("output3.html", s)
    if err != nil {
		return "", fmt.Errorf("Post request failed: %s", err)
    }

	// print cookies
    for _, c := range resp.Cookies() { fmt.Println(c) }

	// extract some data for the get request
	logout_url, _ := get_data(s, "<a href=\"/logout.php?", "\"")
	logout_url = "https://m.facebook.com/logout.php?" + logout_url

	return logout_url, nil
}

func fb_logout(fbengine *FBEngine, logout_url string) (error) {
	///////////////////
	fmt.Println("\nSend fourth request")
	resp, err := fbengine.client.Get(logout_url)
    if err != nil {
		return fmt.Errorf("Get request failed: %s", err)
    }
	defer resp.Body.Close()

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
    for _, c := range resp.Cookies() { fmt.Println(c) }

	return nil
}

func main() {
	flag.Parse()
    args := flag.Args()
    if len(args) < 3 {
        fmt.Println("Please pass arguments: Email, Password, Comment. e.g. fb_login_m.exe test@gmail.com password text to post")
		return
    }

	// here the login infos
	var email = args[0]
	var pass = args[1]
	var text = args[2:]

	var fbengine FBEngine = *NewFBEngine()
	status_update_url, fb_dtsg, privacy, _ := fb_login(&fbengine, email, pass)
	logout_url, _ := fb_post_comment(&fbengine, text, status_update_url, fb_dtsg, privacy)
	fb_logout(&fbengine, logout_url)
}		