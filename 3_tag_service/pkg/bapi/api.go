package bapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"io"
	"net/http"
)

type API struct {
	URL    string
	Client *http.Client
}

type AccessToken struct {
	Token string `json:"token"`
}

const (
	APP_KEY    = "eddycjy"
	APP_SECRET = "go-programming-tour-book"
)

func NewAPI(url string) *API {
	return &API{URL: url, Client: http.DefaultClient}
}

func (a *API) GetTagList(ctx context.Context, name string) ([]byte, error) {
	token, err := a.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	body, err := a.httpGet(ctx, fmt.Sprintf("%s?token=%s&name=%s", "api/v1/tags", token, name))
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (a *API) httpGet(ctx context.Context, path string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", a.URL, path)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	span, newCtx := opentracing.StartSpanFromContext(
		ctx, "HTTP GET: "+a.URL,
		opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
	)
	span.SetTag("url", url)
	_ = opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)

	req = req.WithContext(newCtx)
	client := http.Client{} //Timeout: time.Second * 60
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	defer span.Finish()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (a *API) httpPost(ctx context.Context, path string, jsonData []byte) ([]byte, error) {
	// 拼接完整的 URL
	url := fmt.Sprintf("%s/%s", a.URL, path)

	// 创建 POST 请求，数据通过 URL 编码（application/json）
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}

	// 添加 OpenTracing 标记
	span, _ := opentracing.StartSpanFromContext(
		ctx, "HTTP POST: "+a.URL,
		opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
	)
	span.SetTag("url", url)
	_ = opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)

	// 设置 Content-Type 为 x-www-form-urlencoded
	req.Header.Set("Content-Type", "application/json")

	// 使用新的 Context 替换原来的 Context（可以加入取消机制）
	req = req.WithContext(ctx)

	// 创建 HTTP 客户端并发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	defer span.Finish()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (a *API) getAccessToken(ctx context.Context) (string, error) {

	// 创建表单数据

	// 创建 map 来构造 JSON 数据
	data := map[string]interface{}{
		"app_key":    APP_KEY,
		"app_secret": APP_SECRET,
	}

	// 将 map 序列化为 JSON 数据
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	// 发送 POST 请求
	body, err := a.httpPost(ctx, "auth", jsonData)
	if err != nil {
		return "", err
	}

	var accessToken AccessToken
	err = json.Unmarshal(body, &accessToken)
	if err != nil {
		return "", err
	}

	return accessToken.Token, nil
}
