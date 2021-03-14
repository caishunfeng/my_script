package basehttp

import (
	"bytes"
	log "code.google.com/p/log4go"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type BaseClient struct {
	c *http.Client
}

func NewBaseClient() *BaseClient {
	return &BaseClient{
		c: http.DefaultClient,
	}
}

func (client *BaseClient) HttpSend(method, _url string, data []byte, headerMap map[string]string) ([]byte, error) {
	var (
		req    *http.Request
		resp   *http.Response
		result []byte
		u      *url.URL
		err    error
	)

	for {
		req, err = http.NewRequest(method, _url, bytes.NewBuffer(data))
		if err != nil {
			break
		}

		if method == http.MethodPost {
			req.Header.Set("Content-Type", "application/json")
		}

		u, err = url.Parse(_url)
		if err != nil {
			log.Error("url.Parse %s err: %s", _url, err.Error())
		} else {
			req.Header.Set("Host", u.Host)
		}

		// 放在这个位置，可以覆盖上面的Set操作
		for key, value := range headerMap {
			req.Header.Set(key, value)
		}

		resp, err = client.c.Do(req)
		if err != nil {
			break
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("异常的http状态码: %d", resp.StatusCode)
			break
		}

		result, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			break
		}

		break
	}

	client.WriteLog(method, _url, req.Header, data, result, err)
	return result, err
}

func (client *BaseClient) HttpGet(_url string) ([]byte, error) {
	var (
		resp   *http.Response
		result []byte
		err    error
		req    *http.Request
	)

	for {
		req, err = http.NewRequest(http.MethodGet, _url, nil)
		if err != nil {
			break
		}

		resp, err = client.c.Do(req)
		if err != nil {
			break
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("异常的http状态码: %d", resp.StatusCode)
			break
		}

		result, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			break
		}

		break
	}

	client.WriteLog(http.MethodGet, _url, req.Header, []byte{}, result, err)
	return result, err

}

func (client *BaseClient) HttpPost(_url string, data []byte) ([]byte, error) {
	var (
		resp   *http.Response
		result []byte
		err    error
		req    *http.Request
	)

	for {
		req, err = http.NewRequest(http.MethodPost, _url, bytes.NewBuffer(data))
		if err != nil {
			break
		}

		req.Header.Set("Content-Type", "application/json")
		resp, err = client.c.Do(req)
		if err != nil {
			break
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("异常的http状态码: %d", resp.StatusCode)
			break
		}

		result, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			break
		}

		break
	}

	client.WriteLog(http.MethodPost, _url, req.Header, data, result, err)
	return result, err
}

// 形参存在拷贝的过程，考虑性能开销的话将所有参数封装在一个结构体内，传结构体指针
func (client *BaseClient) WriteLog(method, _url string, header http.Header, reqData, respData []byte, err error) {
	newRespData := bytes.Replace(respData, []byte("\n"), []byte(""), -1)
	log.Info(fmt.Sprintf("调用接口: method: %s, url: %s, headers: %+v, req: %s, resp: %s, err: %v\n", method, _url, header, string(reqData), string(newRespData), err))
}
