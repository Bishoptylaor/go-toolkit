package request

/*
 *  ┏┓      ┏┓
 *┏━┛┻━━━━━━┛┻┓
 *┃　　　━　　  ┃
 *┃   ┳┛ ┗┳   ┃
 *┃           ┃
 *┃     ┻     ┃
 *┗━━━┓     ┏━┛
 *　　 ┃　　　┃神兽保佑
 *　　 ┃　　　┃代码无BUG！
 *　　 ┃　　　┗━━━┓
 *　　 ┃         ┣┓
 *　　 ┃         ┏┛
 *　　 ┗━┓┓┏━━┳┓┏┛
 *　　   ┃┫┫  ┃┫┫
 *      ┗┻┛　 ┗┻┛
 @Time    : 2024/7/13 -- 14:03
 @Author  : bishop ❤️ MONEY
 @Software: GoLand
 @Description: httpwrapper.go
*/

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"
)

type PackPayload struct {
	Url     string            `json:"url"`
	Payload map[string]string `json:"payload"`
	Method  string            `json:"method"`
	Header  http.Header
}

func HttpRequest(ctx context.Context, ppl *PackPayload) ([]byte, error) {
	fun := "HttpRequest"
	fmt.Printf("%s HttpPack %s", fun, ppl)
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	for k, v := range ppl.Payload {
		_ = writer.WriteField(k, v)
	}
	err := writer.Close()
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	req, err := http.NewRequest(ppl.Method, ppl.Url, payload)

	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, errors.New("request server error")
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
