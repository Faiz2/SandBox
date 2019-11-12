package http

import (
	"io"
	"io/ioutil"
	"net/http"
)

func Post(url string, header http.Header, body io.Reader) []byte {
	// 转发
	client := &http.Client{}
	request, _ := http.NewRequest("POST", url, body)
	for k, v := range header {
		request.Header.Add(k, v[0])
	}
	response, _:= client.Do(request)

	responseBody, _:= ioutil.ReadAll(response.Body)
	return responseBody
}
