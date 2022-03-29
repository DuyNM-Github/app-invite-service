package main

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestGenerateTokenSuccess(t *testing.T) {
	url := "http://localhost:8080/token"
	method := "GET"

	request, reqErr := http.NewRequest(method, url, nil)
	if reqErr != nil {
		t.Fatal(reqErr)
	}

	request.Header.Add("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")

	response, resErr := http.DefaultClient.Do(request)
	if resErr != nil {
		t.Fatal(resErr)
	}
	defer response.Body.Close()

	body, ioErr := ioutil.ReadAll(response.Body)
	if ioErr != nil {
		t.Fatal(ioErr)
	}
	t.Log(string(body))
}

func TestGenerateTokenFail(t *testing.T) {
	url := "http://localhost:8080/token"

	response, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	t.Logf("Server reponse with Status Code %v", response.StatusCode)
}

func TestGetAllTokensSuccess(t *testing.T) {
	url := "http://localhost:8080/tokens"
	method := "GET"

	request, reqErr := http.NewRequest(method, url, nil)
	if reqErr != nil {
		t.Fatal(reqErr)
	}

	request.Header.Add("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")

	response, resErr := http.DefaultClient.Do(request)
	if resErr != nil {
		t.Fatal(resErr)
	}

	body, ioErr := ioutil.ReadAll(response.Body)
	if ioErr != nil {
		t.Fatal(ioErr)
	}
	t.Log(string(body))
}

func TestGetAllTokensFail(t *testing.T) {
	url := "http://localhost:8080/tokens"

	response, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	t.Logf("Server reponse with Status Code %v", response.StatusCode)
}

func TestValidatingTokenSuccessExpired(t *testing.T) {
	url := "http://localhost:8080/validate"
	method := "POST"

	request, reqErr := http.NewRequest(method, url, nil)
	if reqErr != nil {
		t.Fatal(reqErr)
	}

	request.Header.Add("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	request.Header.Add("token", "abc123")

	response, resErr := http.DefaultClient.Do(request)
	if resErr != nil {
		t.Fatal(resErr)
	}
	defer response.Body.Close()

	body, ioErr := ioutil.ReadAll(response.Body)
	if ioErr != nil {
		t.Fatal(ioErr)
	}
	t.Log(string(body))
}

func TestValidatingTokenSuccessValid(t *testing.T) {
	url := "http://localhost:8080/validate"
	method := "POST"

	request, reqErr := http.NewRequest(method, url, nil)
	if reqErr != nil {
		t.Fatal(reqErr)
	}

	request.Header.Add("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	request.Header.Add("token", "abc456")

	response, resErr := http.DefaultClient.Do(request)
	if resErr != nil {
		t.Fatal(resErr)
	}
	defer response.Body.Close()

	body, ioErr := ioutil.ReadAll(response.Body)
	if ioErr != nil {
		t.Fatal(ioErr)
	}
	t.Log(string(body))
}

func TestValidatingTokenSuccessNotfound(t *testing.T) {
	url := "http://localhost:8080/validate?token"
	method := "POST"

	request, reqErr := http.NewRequest(method, url, nil)
	if reqErr != nil {
		t.Fatal(reqErr)
	}

	request.Header.Add("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	request.Header.Add("token", "asdefasd1230")

	response, resErr := http.DefaultClient.Do(request)
	if resErr != nil {
		t.Fatal(resErr)
	}
	defer response.Body.Close()

	body, ioErr := ioutil.ReadAll(response.Body)
	if ioErr != nil {
		t.Fatal(ioErr)
	}
	t.Log(string(body))
}

func TestValidatingTokenSuccessInvalid(t *testing.T) {
	url := "http://localhost:8080/validate?token"
	method := "POST"

	request, reqErr := http.NewRequest(method, url, nil)
	if reqErr != nil {
		t.Fatal(reqErr)
	}

	request.Header.Add("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	request.Header.Add("token", "abc")

	response, resErr := http.DefaultClient.Do(request)
	if resErr != nil {
		t.Fatal(resErr)
	}
	defer response.Body.Close()

	body, ioErr := ioutil.ReadAll(response.Body)
	if ioErr != nil {
		t.Fatal(ioErr)
	}
	t.Log(string(body))
}
