package client

import (
	"encoding/json"
	"fmt"
	"io"

	fastshot "github.com/opus-domini/fast-shot"
	"github.com/opus-domini/fast-shot/constant/mime"
)

type DeviceResponse struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Enabled bool `json:"enabled"`
}


type WgClient struct {
	httpClient fastshot.ClientHttpMethods
	cookie string
}


type APIClient interface {
	List() ([]DeviceResponse, error)
	Enable(id string) error
	Disable(id string) error
	Create(name string) error
	Delete(id string) error
}


func (wg *WgClient) authenticate(password string) error {
	payload := map[string]interface{}{"password": password}
	response, err := wg.httpClient.POST("/api/session").Body().AsJSON(payload).Send()

	if err != nil {
		return err
	}

	if !response.Is2xxSuccessful() {
		return fmt.Errorf("unable to create session %d response", response.StatusCode())
	}

	wg.cookie = "connect.sid=" + response.RawResponse.Cookies()[0].Value
	return nil
}


func (wg WgClient) List() ([]DeviceResponse, error) {
	var devices []DeviceResponse
	response, err := wg.httpClient.GET("/api/wireguard/client").Header().Add("Cookie", wg.cookie).Send()

	if err != nil {
		return devices, err
	}

	if !response.Is2xxSuccessful() {
		return devices, fmt.Errorf("unable to fetch devices %d", response.StatusCode())
	}

	data, err := io.ReadAll(response.RawBody())

	if err != nil {
		return devices, err
	}

	err = json.Unmarshal(data, &devices)

	if err != nil {
		return devices, err
	}

	return devices, nil
}

func (wg WgClient) Enable(id string) error {
	url := fmt.Sprintf("/api/wireguard/client/%s/enable", id)
	response, err := wg.httpClient.POST(url).Header().Add("Cookie", wg.cookie).Send()

	if err != nil {
		return err
	}

	if response.StatusCode() != 204 {
		return fmt.Errorf("unexpected status code %d", response.StatusCode())
	}

	return nil
}


func (wg WgClient) Disable(id string) error {
	url := fmt.Sprintf("/api/wireguard/client/%s/disable", id)
	response, err := wg.httpClient.POST(url).Header().Add("Cookie", wg.cookie).Send()

	if err != nil {
		return err
	}

	if response.StatusCode() != 204 {
		return fmt.Errorf("unexpected status code %d", response.StatusCode())
	}

	return nil
}

func (wg WgClient) Create(name string) error {
	payload := map[string]interface{}{"name": name}
	builder := wg.httpClient.POST("/api/wireguard/client").Header().Add("Cookie", wg.cookie)
	response, err := builder.Body().AsJSON(payload).Send()

	if err != nil {
		return err
	}

	if !response.Is2xxSuccessful() {
		return fmt.Errorf("unable to create device %d response", response.StatusCode())
	}

	return nil
}


func (wg WgClient) Delete(id string) error {
	url := fmt.Sprintf("/api/wireguard/client/%s", id)
	response, err := wg.httpClient.DELETE(url).Header().Add("Cookie", wg.cookie).Send()

	if err != nil {
		return err
	}

	if !response.Is2xxSuccessful() {
		return fmt.Errorf("unable to delete device %d response", response.StatusCode())
	}

	return nil
}


func NewWGClient(url, password string) (APIClient, error) {
	builder := fastshot.NewClient(url)
	builder = builder.Header().AddContentType(mime.JSON)
	client := &WgClient{httpClient: builder.Build()}
	err := client.authenticate(password)

	if err != nil {
		return nil, err
	}

	return client, nil
}
