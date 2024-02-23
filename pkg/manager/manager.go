package manager

import (
	"github.com/SergeyKosarchuk/vpn-admin-bot/pkg/client"
)

type DeviceManager struct {
	Client client.APIClient
	Devices []client.Device
}

func (m *DeviceManager) Fetch() error {
	devices, err := m.Client.List()
	if err != nil {
		return err
	}

	m.Devices = devices
	return nil
}

func (m *DeviceManager) Enable(idx int) error {
	err := m.Client.Enable(m.Devices[idx].Id)

	if err != nil {
		return err
	}

	m.Devices[idx].Enabled = true
	return nil
}

func (m *DeviceManager) Disable(idx int) error {
	err := m.Client.Disable(m.Devices[idx].Id)

	if err != nil {
		return err
	}

	m.Devices[idx].Enabled = false
	return nil
}

func (m *DeviceManager) Create(name string) error {
	device, err := m.Client.Create(name)
	if err != nil {
		return nil
	}
	m.Devices = append(m.Devices, device)
	return nil
}

func (m *DeviceManager) Delete(idx int) error {
	err := m.Client.Delete(m.Devices[idx].Id)

	if err != nil {
		return err
	}

	m.Devices = append(m.Devices[:idx], m.Devices[idx+1:]...)
	return nil
}
