package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	ntlm "github.com/vadimi/go-http-ntlm"
)

type TFSServer interface {
	BuildsURL() string
}

type TFSBuildDefinition struct {
	Collection string
	Project    string
	BuildID    string
}

type TFSCredentials struct {
	Username string
	Password string
	Domain   string
}

type TFSOnPremServer struct {
	Instance   string
	Definition TFSBuildDefinition
}

type TFSHostedServer struct {
	Account    string
	Definition TFSBuildDefinition
}

type TFSBuildResponse struct {
	Count int
	Value []TFSBuildStatus
}

type TFSBuildStatus struct {
	Status string
	Result string
}

func (srv TFSOnPremServer) BuildsURL() string {
	baseURL := fmt.Sprintf("https://%s/%s/%s/_apis/build/builds", srv.Instance, url.PathEscape(srv.Definition.Collection), url.PathEscape(srv.Definition.Project))

	query := url.Values{}
	query.Add("api-version", "2.0")
	query.Add("$top", "8")
	query.Add("definitions", srv.Definition.BuildID)
	query.Add("statusFilter", url.QueryEscape("inProgress,completed"))

	return baseURL + "?" + query.Encode()
}

func (srv TFSHostedServer) BuildsURL() string {
	baseURL := fmt.Sprintf("https://%s.visualstudio.com/DefaultCollection/%s/_apis/build/builds", srv.Account, url.PathEscape(srv.Definition.Project))

	query := url.Values{}
	query.Add("api-version", "2.0")
	query.Add("$top", "8")
	query.Add("definitions", srv.Definition.BuildID)
	query.Add("statusFilter", url.QueryEscape("inProgress,completed"))

	return baseURL + "?" + query.Encode()
}

func FetchBuild(srv TFSServer, cred TFSCredentials) ([]TFSBuildStatus, error) {
	if srv == nil {
		return make([]TFSBuildStatus, 0), nil
	}

	client := http.DefaultClient

	switch srv.(type) {
	case TFSOnPremServer:
		client = &http.Client{
			Transport: &ntlm.NtlmTransport{
				Domain:   cred.Domain,
				User:     cred.Username,
				Password: cred.Password,
			},
		}
		break
	case TFSHostedServer:
		client = http.DefaultClient
		break
	}

	fmt.Println("URL:", srv.BuildsURL())

	req, _ := http.NewRequest("GET", srv.BuildsURL(), nil)
	req.SetBasicAuth(cred.Username, cred.Password)

	res, err := client.Do(req)

	if err != nil {
		return make([]TFSBuildStatus, 0), err
	} else if res.StatusCode != http.StatusOK {
		return make([]TFSBuildStatus, 0), nil
	}

	var response TFSBuildResponse
	dec := json.NewDecoder(res.Body)

	dec.Decode(&response)

	return response.Value, nil
}
