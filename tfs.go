package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"crypto/tls"

	ntlm "github.com/vadimi/go-http-ntlm"
)

// TFSServer serves as an interface to abstract away the differences between hosted and on-prem servers.
type TFSServer interface {
	BuildsURL() string
}

// TFSBuildDefinition contains the build definition to monitor on the server.
type TFSBuildDefinition struct {
	Collection string
	Project    string
	BuildID    string
}

// TFSCredentials contains the information required to authenticate to the server.
type TFSCredentials struct {
	Username string
	Password string
	Domain   string
}

// TFSOnPremServer contains the information required to access an on-premise server.
type TFSOnPremServer struct {
	Instance   string
	Definition TFSBuildDefinition
}

// TFSHostedServer contains the information required to access a VSO instance.
type TFSHostedServer struct {
	Account    string
	Definition TFSBuildDefinition
}

// TFSBuildResponse contains the results of the last n builds.
type TFSBuildResponse struct {
	Count int
	Value []TFSBuildStatus
}

// TFSBuildStatus contains the result for a single build.
type TFSBuildStatus struct {
	Status string
	Result string
}

// BuildsURL generates the URL for accessing an on-prem server.
func (srv TFSOnPremServer) BuildsURL() string {
	baseURL := fmt.Sprintf("https://%s/%s/%s/_apis/build/builds", srv.Instance, url.PathEscape(srv.Definition.Collection), url.PathEscape(srv.Definition.Project))

	query := url.Values{}
	query.Add("api-version", "2.0")
	query.Add("$top", "8")
	query.Add("definitions", srv.Definition.BuildID)
	query.Add("statusFilter", url.QueryEscape("inProgress,completed"))

	return baseURL + "?" + query.Encode()
}

// BuildsURL generates the URL for accessing a VSO instance.
func (srv TFSHostedServer) BuildsURL() string {
	baseURL := fmt.Sprintf("https://%s.visualstudio.com/DefaultCollection/%s/_apis/build/builds", srv.Account, url.PathEscape(srv.Definition.Project))

	query := url.Values{}
	query.Add("api-version", "2.0")
	query.Add("$top", "8")
	query.Add("definitions", srv.Definition.BuildID)
	query.Add("statusFilter", url.QueryEscape("inProgress,completed"))

	return baseURL + "?" + query.Encode()
}

// FetchBuild makes the http request to the specified server using the specified credentials.
func FetchBuild(srv TFSServer, cred TFSCredentials) ([]TFSBuildStatus, error) {
	if srv == nil {
		return make([]TFSBuildStatus, 0), nil
	}

	client := http.DefaultClient

	switch srv.(type) {
	case TFSOnPremServer:
		client = &http.Client{
			Transport: &ntlm.NtlmTransport{
				Domain:          cred.Domain,
				User:            cred.Username,
				Password:        cred.Password,
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
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
