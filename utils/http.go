package utils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/opentracing/opentracing-go"
)

// BuildUrl 构建 GET URL
func BuildUrl(u string, p map[string]string) (string, error) {
	params := url.Values{}
	rawUrl, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	for k, v := range p {
		params.Set(k, v)
	}
	rawUrl.RawQuery = params.Encode()
	return rawUrl.String(), nil
}

// parseTimeout 解析超时时间；默认是 1s
func parseTimeout(t ...int) time.Duration {
	var to int
	if len(t) > 0 {
		to = t[0]
	} else {
		to = 1
	}

	return time.Duration(to)
}

// HttpGet 发起 Get 请求
func HttpGet(ctx context.Context, u string, header map[string]string, timeout ...int) ([]byte, int, error) {
	client := &http.Client{
		Timeout: parseTimeout(timeout...) * time.Second,
	}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, 0, err
	}
	for k, v := range header {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	code := resp.StatusCode
	if code != 200 {
		return nil, code, errors.New(fmt.Sprintf("code status:%d", code))
	}
	ret, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return nil, code, errors.New(fmt.Sprintf("read body:%s", err.Error()))
	}

	// tracing
	span, _ := opentracing.StartSpanFromContext(ctx, "http-get")
	if span != nil {
		span.LogKV("url", u, "code", code, "response", string(ret))
		defer span.Finish()
	}

	return ret, code, nil
}

// HttpPostJson 发起 Post Json 请求
func HttpPostJson(ctx context.Context, u string, json []byte, header map[string]string, timeout ...int) ([]byte, int, error) {
	client := &http.Client{
		Timeout: parseTimeout(timeout...) * time.Second,
	}
	req, err := http.NewRequest("POST", u, bytes.NewReader(json))
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range header {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	code := resp.StatusCode
	if code != 200 {
		return nil, code, errors.New(fmt.Sprintf("code status:%d", code))
	}
	ret, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return nil, code, errors.New(fmt.Sprintf("read body:%s", err.Error()))
	}

	span, _ := opentracing.StartSpanFromContext(ctx, "http-post-json")
	if span != nil {
		span.LogKV("url", u, "header", req.Header, "request", json, "code", code, "response", string(ret))
		defer span.Finish()
	}

	return ret, code, nil
}

// HttpPost 发起 Post 请求
func HttpPost(ctx context.Context, u string, param []byte, header map[string]interface{}, timeout ...int) ([]byte, int, error) {
	client := &http.Client{
		Timeout: parseTimeout(timeout...) * time.Second,
	}
	req, err := http.NewRequest("POST", u, bytes.NewReader(param))
	if err != nil {
		return nil, 0, err
	}
	for k, v := range header {
		req.Header.Add(k, fmt.Sprintf("%v", v))
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	code := res.StatusCode
	if code != 200 {
		return nil, code, errors.New(fmt.Sprintf("code status:%d", code))
	}
	ret, err := ioutil.ReadAll(res.Body)
	_ = res.Body.Close()
	if err != nil {
		return nil, code, errors.New(fmt.Sprintf("read body:%s", err.Error()))
	}

	// tracing
	span, _ := opentracing.StartSpanFromContext(ctx, "http-post")
	if span != nil {
		span.LogKV("url", u, "header", req.Header, "request", string(param), "code", code, "response", string(ret))
		defer span.Finish()
	}
	return ret, code, nil
}

// HttpPostForm 发起 Post Form 请求
func HttpPostForm(ctx context.Context, u string, data map[string]string, timeout ...int) ([]byte, int, error) {
	values := url.Values{}
	for k, v := range data {
		values.Set(k, v)
	}
	params := []byte(values.Encode())
	return HttpPost(ctx, u, params, map[string]interface{}{
		"Content-Type": "application/x-www-form-urlencoded",
	}, timeout...)
}

// HttpPostFile 发起 Post file 请求
func HttpPostFile(ctx context.Context, url, formFileName, filePath string, params map[string]string, timeout ...int) ([]byte, code, error) {
	// 读取文件
	file, err := os.Open(filePath)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	// 创建 form 表单
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	// 设置上传文件的 form 选项
	part, err := writer.CreateFormFile(formFileName, formFileName)
	if err != nil {
		return nil, 0, err
	}
	// 将文件放入 form
	_, err = io.Copy(part, file)
	// 将其他参数放入 form
	for k, v := range params {
		_ = writer.WriteField(k, v)
	}
	// 关闭 writer 不能使用 defer
	err = writer.Close()
	if err != nil {
		return nil, 0, err
	}

	// 创建请求
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, 0, err
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())
	// 发起请求
	client := &http.Client{
		Timeout: parseTimeout(timeout...) * time.Second,
	}
	res, err := client.Do(request)
	if err != nil {
		return nil, 0, err
	}
	// 读取响应
	code := res.StatusCode
	if code != 200 {
		return nil, code, errors.New(fmt.Sprintf("code status:%d", code))
	}
	ret, err := ioutil.ReadAll(res.Body)
	_ = res.Body.Close()
	if err != nil {
		return nil, code, errors.New(fmt.Sprintf("read body:%s", err.Error()))
	}

	// tracing
	span, _ := opentracing.StartSpanFromContext(ctx, "http-post-file")
	if span != nil {
		span.LogKV("url", url, "request", fmt.Sprintf("%v", params), "code", code, "response", string(ret))
		defer span.Finish()
	}
	return ret, code, nil
}
