package request

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

const (
	defaultTimeout = 10 * time.Second
)

// GetPage executes a http GET and attempts retries.
func GetPage(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", errors.WithStack(err)
	}

	resp, err := (&http.Client{Timeout: defaultTimeout}).Do(setHeader(req))
	if err != nil {
		return "", errors.WithStack(err)
	}
	defer closeResponseBody(resp)

	buf := new(strings.Builder)
	if _, err := io.Copy(buf, resp.Body); err != nil {
		return "", errors.WithStack(err)
	}

	return buf.String(), nil
}

func setHeader(req *http.Request) *http.Request {
	req.Header.Set("Authority", "www.sec.gov")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Sec-Ch-Ua", "\"Not?A_Brand\";v=\"8\", \"Chromium\";v=\"108\", \"Google Chrome\";v=\"108\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")

	return req
}

func closeResponseBody(resp *http.Response) {
	if _, err := io.Copy(io.Discard, resp.Body); err != nil {
		log.Error(err)
	}
	if err := resp.Body.Close(); err != nil {
		log.Error(err)
	}
}
