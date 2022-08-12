package baidu

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

const Host = "pan.baidu.com"

var Api = NewApi()

func (v *api) SetCookies(_ *url.URL, cookies []*http.Cookie) {
	v.cookies = cookies
}

func (v *api) Cookies(_ *url.URL) []*http.Cookie {
	return v.cookies
}

type api struct {
	Client  *http.Client
	cookies []*http.Cookie
}
type apiTransport struct {
	r http.RoundTripper
}
type resp struct {
	Errno     int    `json:"errno"`
	ErrMsg    string `json:"err_msg"`
	RequestId string `json:"request_id"`
}

// File 百度通用文件类型
type File struct {
	Category       string `json:"category"`
	FsId           string `json:"fs_id"`
	Isdir          string `json:"isdir"`
	LocalCtime     string `json:"local_ctime"`
	LocalMtime     string `json:"local_mtime"`
	Md5            string `json:"md5"`
	Path           string `json:"path"`
	ServerCtime    string `json:"server_ctime"`
	ServerFilename string `json:"server_filename"`
	ServerMtime    string `json:"server_mtime"`
	Size           string `json:"size"`
}

func (m apiTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	request.Header.Add("Referer", "pan.baidu.com")
	request.Header.Add("User-Agent", "pan.baidu.com")
	return m.r.RoundTrip(request)
	//return
}

func NewApi() *api {
	tmp := &api{}
	tmp.Client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar:     tmp,
		Timeout: 20 * time.Second,
		Transport: apiTransport{
			r: http.DefaultTransport,
		},
	}
	return tmp
}

func GetSurl(u string) string {
	res, err := Api.Client.Get(u)
	if err != nil {
		return ""
	}
	defer res.Body.Close()
	if res.StatusCode != 302 {
		return ""
	}
	str := res.Header.Get("Location")
	sub := strings.Split(str, "&")
	sub = strings.Split(sub[0], "=")
	if len(sub) == 2 {
		return sub[1]
	}
	return ""
}
