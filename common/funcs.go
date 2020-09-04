package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func HttpGet(url string, params map[string]interface{}) (retData map[string]interface{}) {
	paramStr := "?"
	i := 0
	for index, param := range params {
		i++
		if i < len(params) {
			paramStr += index + "=" + param.(string) + "&"
		} else {
			paramStr += index + "=" + param.(string)
		}

	}

	resp, err := http.Get(url + paramStr)

	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
	}

	if err := json.Unmarshal(body, &retData); err == nil {
		fmt.Println(err)
	}
	return retData
}

func HttpPost(url string, params string) (stringBody string) {
	resp, err := http.Post(
		url,
		"application/x-www-form-urlencoded",
		strings.NewReader(params))

	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
	}

	return string(body)
}
