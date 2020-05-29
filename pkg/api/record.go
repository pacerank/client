package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type SendRecordReply struct{}

func (api *Api) SendRecord(payload []byte) (*DefaultReplyStructure, *SendRecordReply, error) {
	request, err := http.NewRequest("POST", fmt.Sprintf("%s/digest/v1/records", api.url), bytes.NewBuffer(payload))
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

	var endpointReply SendRecordReply
	_ = json.Unmarshal(reply.Content, &endpointReply)
	return &reply, &endpointReply, nil
}
