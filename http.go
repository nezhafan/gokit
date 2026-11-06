package toolkit

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	postJsonContentType = "application/json"
	postFormContentType = "application/x-www-form-urlencoded"
)

// 结构体支持
type Request struct {
	url    string
	method string
	body   []byte
	header map[string]string
}

func HttpGet(api string, params map[string]string) *Request {
	if len(params) > 0 {
		values := url.Values{}
		for k, v := range params {
			values.Add(k, v)
		}
		api += "?" + values.Encode()
	}

	r := &Request{
		url:    api,
		method: http.MethodGet,
	}

	return r
}

func HttpPostJson(api string, data any) *Request {
	r := &Request{
		url:    api,
		method: http.MethodPost,
		header: map[string]string{"Content-Type": postJsonContentType},
	}

	if data != nil {
		if bs, err := json.Marshal(data); err == nil {
			r.body = bs
		}
	}

	return r
}

func HttpPostForm(api string, data map[string]any) *Request {
	r := &Request{
		url:    api,
		method: http.MethodPost,
		header: map[string]string{"Content-Type": postFormContentType},
	}

	if len(data) > 0 {
		values := url.Values{}
		for k, v := range data {
			values.Add(k, fmt.Sprintf("%v", v))
		}
		r.body = []byte(values.Encode())
	}

	return r
}

func (r *Request) SetHeader(key, value string) *Request {
	if r.header == nil {
		r.header = make(map[string]string)
	}
	r.header[key] = value
	return r
}

func (r *Request) Do(ctx context.Context, data any) (err error) {
	start := time.Now()
	var b []byte
	defer r.log(start, &b, &err)

	var req *http.Request
	req, err = http.NewRequestWithContext(ctx, r.method, r.url, bytes.NewReader(r.body))
	if err != nil {
		return
	}

	for k, v := range r.header {
		req.Header.Set(k, v)
	}

	var res *http.Response
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	// 解析body
	b, err = io.ReadAll(res.Body)
	if err != nil {
		return
	}

	// 判断状态码，放在读取 body 之后，避免被 ReadAll/json.Unmarshal 覆盖
	if res.StatusCode >= http.StatusMultipleChoices {
		err = errors.New("HTTP状态码异常: " + res.Status)
		return
	}

	if data != nil && len(b) > 0 {
		err = json.Unmarshal(b, data)
	}

	return
}

type requestLog struct {
	Time     string `json:"time"`
	Url      string `json:"url"`
	Method   string `json:"method"`
	Cost     string `json:"cost"`
	Err      string `json:"error,omitempty"`
	Body     any    `json:"body,omitempty"`
	Response any    `json:"response,omitempty"`
}

func (r *Request) log(start time.Time, res *[]byte, err *error) {
	cost := int(time.Since(start).Milliseconds())
	data := requestLog{
		Time:   start.Format(time.DateTime),
		Url:    r.url,
		Method: r.method,
		Cost:   strconv.Itoa(cost) + "ms",
	}
	if len(r.body) > 0 {
		if r.header["Content-Type"] == postJsonContentType {
			data.Body = json.RawMessage(r.body)
		} else {
			data.Body = string(r.body)
		}
	}

	if res != nil {
		b := *res
		// 防止超大日志
		if len(b) >= 4096 {
			b = b[:4096]
		}
		// 如果返回的是json格式则防止二次编码
		if json.Valid(b) {
			data.Response = json.RawMessage(b)
		} else {
			data.Response = string(b)
		}
	}
	if err != nil && *err != nil {
		data.Err = (*err).Error()
		LogError("HTTP请求", "data", data)
	} else if cost >= 1000 {
		LogWarn("HTTP请求", "data", data)
	} else {
		LogInfo("HTTP请求", "data", data)
	}
}
