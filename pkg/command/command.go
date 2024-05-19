// Package command implements Wrapper around `net/http` to make HTTP requests to the admin REST API.
package command

import (
	api "github.com/SergeyKosarchuk/vpn-admin-bot/pkg/client"
)

type APIClient interface {
	List() ([]api.DeviceResponse, error)
	Enable(id string) error
	Disable(id string) error
	Create(name string) error
	Delete(id string) error
	GetConfig(id string) ([]byte, error)
}
