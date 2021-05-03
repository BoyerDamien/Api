package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

/****************************************************************************************
 *
 ***************************************************************************************/

type API interface {
	SetHeader(key string, value string)
	Get(uri string) (APIResponse, error)
	Post(uri string, data interface{}) (APIResponse, error)
	Put(uri string, data interface{}) (APIResponse, error)
	Delete(uri string) (APIResponse, error)
}

type apiImpl struct {
	url     string
	headers map[string]string
}

type APIResponse struct {
	response *http.Response
	data     []byte
}

/****************************************************************************************
 *
 ***************************************************************************************/

func (s *apiImpl) ReadHTTPResponse(response *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return []byte(""), err
	}
	return body, nil
}

func (s *apiImpl) doRequest(reqType, uri string, data interface{}) (APIResponse, error) {
	client := &http.Client{}
	endpoint := fmt.Sprintf("%s/%s", s.url, uri)
	marshData, err := json.Marshal(data)
	if err != nil {
		return APIResponse{}, err
	}
	req, err := http.NewRequest(reqType, endpoint, bytes.NewBuffer(marshData))
	if err != nil {
		return APIResponse{}, err
	}
	for key, elem := range s.headers {
		req.Header.Add(key, elem)
	}
	resp, err := client.Do(req)
	if err != nil {
		return APIResponse{}, err
	}
	respData, err := s.ReadHTTPResponse(resp)
	if err != nil {
		return APIResponse{}, err
	}
	return APIResponse{response: resp, data: respData}, nil
}

func New(url string) API {
	return &apiImpl{url: url, headers: make(map[string]string)}
}

func (s *apiImpl) SetHeader(key string, value string) {
	s.headers[key] = value
}
func (s *apiImpl) Get(uri string) (APIResponse, error) {
	return s.doRequest("GET", uri, "")
}

func (s *apiImpl) Post(uri string, data interface{}) (APIResponse, error) {
	return s.doRequest("POST", uri, data)
}

func (s *apiImpl) Put(uri string, data interface{}) (APIResponse, error) {
	return s.doRequest("PUT", uri, data)
}

func (s *apiImpl) Delete(uri string) (APIResponse, error) {
	return s.doRequest("Delete", uri, "")
}
