package model

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"strings"
)

type HttpClient struct {
	Client              *http.Client
	Cookie, RequestHead map[string]string
}

type ResponseModel struct {
	Code       string
	Data       string
	SourceData interface{}
}

func New(o *cookiejar.Options) *HttpClient {
	client := http.Client{}
	jar, err := cookiejar.New(o)
	if err != nil {
		log.Print("构造cookiejar时,发生错误->", err)
	}
	client.Jar = jar
	//登录一下
	hc := &HttpClient{}
	hc.Client = &client
	return hc
}

func (hc *HttpClient) GET(url, body string) (*ResponseModel, error) {

	return hc.SendRequest(url, "GET", body, hc.RequestHead, hc.Cookie)
}

func (hc *HttpClient) SendRequest(url, method, body string, requestHead, cookie map[string]string) (*ResponseModel, error) {

	request, err := http.NewRequest(method, url, strings.NewReader(body))

	if err != nil {
		log.Print("构造请求的时候发生错误,-->url,", url, ".请求体,", body, "..err-->", err)
		return nil, err
	}
	setHead(request, requestHead)
	setCookie(request, cookie)
	response, err := hc.Client.Do(request)
	if err != nil {
		log.Print("发生请求时,发生错误,-->url,", url, ".请求体,", body, "..err-->", err)
		return nil, err
	}

	all, err := ioutil.ReadAll(response.Body)
	if err != nil {

	}
	return &ResponseModel{response.Status, string(all), response}, nil
}

//
//
////发送请求, 返回响应提的 字符串数据.
//func (rm *RequestModel) SendRequestForResponseBody(hc *http.Client, cookie map[string]string) string {
//	request, err := http.NewRequest(rm.Method, rm.Url, strings.NewReader(rm.RequestBody))
//	if err != nil {
//		log.Print("构造请求的时候发生错误,-->url,", rm.Url, ".请求体,", rm.RequestBody, "..err-->", err)
//		return ""
//	}
//	//填充请求头...
//	setHead(request, rm.RequestHead)
//	//填充cookie
//	setCookie(request, cookie)
//	response, err := hc.Do(request)
//	if err != nil {
//		log.Print("发生请求时,发生错误,-->url,", rm.Url, ".请求体,", rm.RequestBody, "..err-->", err)
//		return ""
//	}
//	list := setCookieByCookieList(response.Header["Set-Cookie"])
//	hc.Jar.SetCookies(request.URL, list)
//	defer func() {
//		if response != nil {
//			err := response.Body.Close()
//			if err != nil {
//				log.Print("释放响应体流的时候发生错误.", err)
//			}
//		}
//	}()
//	body, err := ioutil.ReadAll(response.Body)
//	if err != nil {
//		log.Print("读取响应体流时,发生错误,-->url,", rm.Url, ".请求体,", rm.RequestBody, "..err-->", err)
//		return ""
//	}
//	return string(body)
//}

func setCookie(r *http.Request, cookie map[string]string) {
	for k, v := range cookie {
		c := &http.Cookie{Name: k, Value: v}
		r.AddCookie(c)
	}

}

func setCookieByCookieList(cookies []string) []*http.Cookie {
	var result []*http.Cookie
	for _, v := range cookies {
		split := strings.Split(v, ";")
		cookie := http.Cookie{}
		for _, v1 := range split {
			split1 := strings.Split(v1, "=")
			if split1[0] == " path" {
				cookie.Path = split1[1]
				break
			}
			if split1[0] == " Max-Age" {
				i, err := strconv.Atoi(split1[1])
				if err != nil {
					log.Print("字符串转整形时,发生错误.->", err)
					break
				}
				cookie.MaxAge = i
				break
			}
			if split1[0] == " HttpOnly" {
				cookie.HttpOnly = true
				break
			}
			cookie.Name = split1[0]
			cookie.Value = split1[1]

		}
		result = append(result, &cookie)
	}

	return result

}

func setHead(r *http.Request, head map[string]string) {
	for k, v := range head {
		r.Header.Set(k, v)
	}
}

//func MakeRequest(Url string, Method string, RequestHead map[string]string, RequestBody string) *RequestModel {
//	m := httpTool.MakeHeader()
//	if RequestHead != nil {
//		m = RequestHead
//	}
//	model := RequestModel{Url, Method, m, RequestBody}
//	return &model
//
//}
