package http

import (
	"io/ioutil"
	"net/http"
)

func Get(url string, header http.Header) []byte {
	// 转发
	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	for k, v := range header {
		request.Header.Add(k, v[0])
	}
	response, _:= client.Do(request)

	responseBody, _:= ioutil.ReadAll(response.Body)
	return responseBody
}
