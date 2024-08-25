// Package client implements Wrapper around `net/http` to make HTTP requests to the admin REST API.
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
)

const APPLICATION_JSON = "application/json"

type WGClient struct {
	httpClient *http.Client
	host       string
	baseUrl    string
}

type DeviceResponse struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

func (wg WGClient) authenticate(password string) error {
	payload := map[string]interface{}{"password": password}
	body, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	url := wg.host + "/api/session"
	response, err := wg.httpClient.Post(url, APPLICATION_JSON, bytes.NewReader(body))

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d", response.StatusCode)
	}

	return nil
}

// List - returns list of all clients.
func (wg WGClient) List() ([]DeviceResponse, error) {
	var devices []DeviceResponse
	response, err := wg.httpClient.Get(wg.baseUrl)

	if err != nil {
		return devices, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return devices, fmt.Errorf("unexpected status code %d", response.StatusCode)
	}

	data, err := io.ReadAll(response.Body)

	if err != nil {
		return devices, err
	}

	err = json.Unmarshal(data, &devices)
	if err != nil {
		return devices, err
	}

	return devices, nil
}

// Enable client by id.
func (wg WGClient) Enable(id string) error {
	url := fmt.Sprintf("%s/%s/enable", wg.baseUrl, id)

	response, err := wg.httpClient.Post(url, APPLICATION_JSON, nil)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code %d", response.StatusCode)
	}

	return nil
}

// Disable client by id.
func (wg WGClient) Disable(id string) error {
	url := fmt.Sprintf("%s/%s/disable", wg.baseUrl, id)

	response, err := wg.httpClient.Post(url, APPLICATION_JSON, nil)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code %d", response.StatusCode)
	}

	return nil
}

// Create new client with the given name.
func (wg WGClient) Create(name string) error {
	payload := map[string]interface{}{"name": name}
	body, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	response, err := wg.httpClient.Post(wg.baseUrl, APPLICATION_JSON, bytes.NewReader(body))

	if err != nil {
		return err
	}

	// There are no `id` field in response so we are unable to return `DeviceResponse`.
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d", response.StatusCode)
	}

	return nil
}

// Delete client by id.
func (wg WGClient) Delete(id string) error {
	url := fmt.Sprintf("%s/%s", wg.baseUrl, id)
	req, err := http.NewRequest(http.MethodDelete, url, nil)

	if err != nil {
		return err
	}

	response, err := wg.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code %d", response.StatusCode)
	}

	return nil
}

// GetConfig returns encoded client config by id.
func (wg WGClient) GetConfig(id string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/configuration", wg.baseUrl, id)
	response, err := wg.httpClient.Get(url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", response.StatusCode)
	}

	return io.ReadAll(response.Body)
}

// Create authenticated `WGClient` ready to use.
func NewWGClient(host, password string) (WGClient, error) {
	// Use non-empty cookiejar to save sessionid cookie from authenticate request.
	jar, err := cookiejar.New(nil)
	wg := WGClient{}

	if err != nil {
		return wg, err
	}

	wg.httpClient = &http.Client{Jar: jar}
	wg.host = host
	wg.baseUrl = host + "/api/wireguard/client"
	err = wg.authenticate(password)

	if err != nil {
		return wg, fmt.Errorf("unable to authenticate a client %w", err)
	}

	return wg, nil
}
