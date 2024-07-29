package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var URL = getenvOrDefault("TT_API_GATEWAY_URL", "https://tikwm.com")
var tikwmMutex = &sync.Mutex{}
var tikwmTimeout = time.Second + time.Millisecond*100

func main() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Any("/*UrlPath", func(c *gin.Context) {
		tikwmMutex.Lock()
		code, body, err := raw(c.Request.URL.Path, c.Request.URL.Query())
		go time.AfterFunc(tikwmTimeout, tikwmMutex.Unlock)

		switch {
		case err != nil:
			c.JSON(502, gin.H{
				"ok":     false,
				"msg":    "server error",
				"status": code,
				"body":   body,
			})

		case !json.Valid(body):
			c.JSON(404, gin.H{
				"ok":     false,
				"msg":    "url path not found or other service's unexpected return",
				"status": 404,
				"body":   body,
			})

		default:
			c.String(code, string(body))
		}
	})

	router.Run(":80")
}

func raw(method string, query url.Values) (int, []byte, error) {
	url := fmt.Sprintf("%s/%s", URL, method)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 502, nil, err
	}

	q := req.URL.Query()
	for key, values := range query {
		for _, val := range values {
			q.Add(key, val)
		}
	}
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if resp == nil {
			return 502, nil, err
		}
		return resp.StatusCode, nil, err
	}
	defer resp.Body.Close()

	buffer, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, err
	}

	return resp.StatusCode, buffer, nil
}

func getenvOrDefault(env string, dflt string) string {
	ret, ok := os.LookupEnv(env)
	if ok {
		return ret
	}
	return dflt
}
