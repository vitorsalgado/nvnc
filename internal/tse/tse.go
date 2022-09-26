package tse

import (
	"archive/zip"
	"bytes"
	"io"
	"net/http"
)

func FetchCandidates(client *http.Client, u, file string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	// simulate a browser request headers
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:104.0) Gecko/20100101 Firefox/104.0")
	req.Header.Add("accept-encoding", "gzip, deflate, br")
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Add("accept-language", "en-US,en;q=0.5")
	req.Header.Add("connection", "keep-alive")
	req.Header.Add("upgrade-insecure-requests", "1")
	req.Header.Add("sec-fetch-dest", "document")
	req.Header.Add("sec-fetch-mode", "navigate")
	req.Header.Add("sec-fetch-site", "none")
	req.Header.Add("sec-fetch-user", "?1")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	zz, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return nil, err
	}

	var data *zip.File

	for _, f := range zz.File {
		if file == f.FileInfo().Name() {
			data = f
		}
	}

	if data == nil {
		return nil, nil
	}

	opened, err := data.Open()
	if err != nil {
		return nil, err
	}

	defer opened.Close()

	return io.ReadAll(opened)
}
