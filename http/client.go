package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

func main() {
	request, err := http.NewRequest(http.MethodGet, "https://www.baidu.com", nil)
	request.Header.Add("User-Agent", "Mozilla/7.0 (iPhone;)")

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}

	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	string, err := httputil.DumpResponse(response, true)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", string)

}
