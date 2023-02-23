package serviceclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type ServiceClient struct {
	ServiceUrl string
}

func New(serviceUrl string) *ServiceClient {
	return &ServiceClient{
		ServiceUrl: serviceUrl,
	}
}

func MakeRequest[Req any, Res any](ctx context.Context, client *ServiceClient, endpointPath string, reqPayload Req, resPayload *Res) error {
	reqUrl, err := url.JoinPath(client.ServiceUrl, endpointPath)
	if err != nil {
		return err
	}

	reqBody, err := json.Marshal(reqPayload)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, reqUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	httpReq.Header.Add("Content-Type", "application/json")

	httpRes, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpRes.Body.Close()

	log.Println(reqUrl, httpRes.StatusCode)

	ok := httpRes.StatusCode >= http.StatusOK && httpRes.StatusCode < http.StatusMultipleChoices
	if !ok {
		return fmt.Errorf("Got unsuccessful status code: %d", httpRes.StatusCode)
	}

	err = json.NewDecoder(httpRes.Body).Decode(resPayload)
	if err != nil {
		return err
	}

	return nil
}
