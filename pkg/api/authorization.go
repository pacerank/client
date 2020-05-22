package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type InitializeAuthorizationFlowReply struct {
	AuthorizationId  string `json:"authorization_id"`
	AuthorizationUrl string `json:"authorization_url"`
}

func (api *Api) InitializeAuthorizationFlow() (*DefaultReplyStructure, *InitializeAuthorizationFlowReply, error) {
	type payload struct {
		ClientName string `json:"client_name"`
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, nil, err
	}

	b, err := json.Marshal(payload{
		ClientName: hostname,
	})
	if err != nil {
		return nil, nil, err
	}

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/digest/v1/authorizations", api.url), bytes.NewBuffer(b))
	if err != nil {
		return nil, nil, err
	}

	response, err := api.client.Do(request)
	if err != nil {
		return nil, nil, err
	}

	defer response.Body.Close()

	var reply DefaultReplyStructure
	err = json.NewDecoder(response.Body).Decode(&reply)
	if err != nil {
		return nil, nil, err
	}

	var endpointReply InitializeAuthorizationFlowReply
	err = json.Unmarshal(reply.Content, &endpointReply)
	if err != nil {
		return nil, nil, err
	}

	return &reply, &endpointReply, nil
}

type ConfirmAuthorizationFlowReply struct {
	AuthorizationId   string `json:"authorization_id"`
	Authorized        bool   `json:"authorized"`
	ClientName        string `json:"client_name"`
	ApiKey            string `json:"api_key"`
	UserSignatureName string `json:"user_signature_name"`
}

func (api *Api) ConfirmAuthorizationFlow(authorizationId string) (*DefaultReplyStructure, *ConfirmAuthorizationFlowReply, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/digest/v1/authorizations/%s", api.url, authorizationId), nil)
	if err != nil {
		return nil, nil, err
	}

	response, err := api.client.Do(request)
	if err != nil {
		return nil, nil, err
	}

	defer response.Body.Close()

	var reply DefaultReplyStructure
	err = json.NewDecoder(response.Body).Decode(&reply)
	if err != nil {
		return nil, nil, err
	}

	var endpointReply ConfirmAuthorizationFlowReply
	err = json.Unmarshal(reply.Content, &endpointReply)
	if err != nil {
		return nil, nil, err
	}

	return &reply, &endpointReply, nil
}
