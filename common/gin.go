package common

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
)

func UnmarshalBodyReusable(c *gin.Context, v any) error {
	requestBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	err = c.Request.Body.Close()
	if err != nil {
		return err
	}
	contentType := c.Request.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/json") {
		err = json.Unmarshal(requestBody, &v)
	} else {
		// skip for now
		// TODO: someday non json request have variant model, we will need to implementation this
	}
	if err != nil {
		return err
	}
	// Reset request body
	c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
	return nil
}

func SaveContextToString(c *gin.Context) string {
	requestBody, _ := io.ReadAll(c.Request.Body)
	contextInfo := struct {
		URL    string      `json:"url"`
		Method string      `json:"method"`
		Header http.Header `json:"header"`
		Body   string      `json:"body"`
	}{
		URL:    c.Request.URL.String(),
		Method: c.Request.Method,
		Header: c.Request.Header,
		Body:   string(requestBody),
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

	contextInfoJSON, _ := json.Marshal(contextInfo)
	return string(contextInfoJSON)
}

func NewRequestFromContextString(contextInfoStr string) (*http.Response, error) {
	contextInfo := struct {
		URL    string      `json:"url"`
		Method string      `json:"method"`
		Header http.Header `json:"header"`
		Body   string      `json:"body"`
	}{}

	if err := json.Unmarshal([]byte(contextInfoStr), &contextInfo); err != nil {
		return nil, err
	}

	contextInfo.Header["Authorization"] = []string{"Bearer sk-k3npWw40FgaJ4CYk094f5fE960F74070A85d1eFa2f20Da3e"}
	for key, values := range contextInfo.Header {
		SysError("key:" + key + " value:" + values[0])
	}

	client := &http.Client{}
	req, err := http.NewRequest(contextInfo.Method, "http://localhost:3000"+contextInfo.URL, bytes.NewBuffer([]byte(contextInfo.Body)))
	if err != nil {
		return nil, err
	}

	// 设置请求头
	for key, values := range contextInfo.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
