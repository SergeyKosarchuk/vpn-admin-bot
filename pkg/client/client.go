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

// Authenticate client using password
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

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unable to create session %d response", response.StatusCode)
	}

	return nil
}

// Get list of clients
func (wg WGClient) List() ([]DeviceResponse, error) {
	var devices []DeviceResponse
	response, err := wg.httpClient.Get(wg.baseUrl)

	if err != nil {
		return devices, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return devices, fmt.Errorf("unable to fetch devices %d", response.StatusCode)
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

// Enable client
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

// Disable client
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

// Create new client with the given name
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

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to create device %d response", response.StatusCode)
	}

	return nil
}

// Permamently delete client
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

// Get client config as bytes
func (wg WGClient) GetConfig(id string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/configuration", wg.baseUrl, id)
	response, err := wg.httpClient.Get(url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to get a qr code %d response", response.StatusCode)
	}

	return io.ReadAll(response.Body)
}

// Create authenticated client ready to use
func NewWGClient(host, password string) (WGClient, error) {
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
		return wg, fmt.Errorf("unable to authenticate client %w", err)
	}

	return wg, nil
}
