package manager_test

import (
	"testing"

	"github.com/SergeyKosarchuk/vpn-admin-bot/pkg/client"
	"github.com/SergeyKosarchuk/vpn-admin-bot/pkg/manager"
)


type APIClientStub struct {
	list func() ([]client.Device, error)
	create func(name string) (client.Device, error)
	enable func (id string) error
	disable func (id string) error
	delete func(id string) error
}

func (a APIClientStub) List() ([]client.Device, error) {
	return a.list()
}

func (a APIClientStub) Create(name string) (client.Device, error) {
	return a.create(name)
}

func (a APIClientStub) Enable(id string) error {
	return a.enable(id)
}

func (a APIClientStub) Disable(id string) error {
	return a.disable(id)
}

func (a APIClientStub) Delete(id string) error {
	return a.delete(id)
}



func TestFetch(t *testing.T) {
	devices := []client.Device{
		{Id: "id1", Name: "Test1", Enabled: true},
		{Id: "id2", Name: "Test2", Enabled: false},
	}
	stub := APIClientStub{}
	stub.list = func() ([]client.Device, error) {return devices, nil}
	man := manager.DeviceManager{Client: stub}

	err := man.Fetch()
	if err != nil {
		t.Error("assert nil error, got", err)
	}

	if len(man.Devices) != len(devices) {
		t.Error("assert 2 devices, got", len(man.Devices))
	}
}

func TestEnable(t *testing.T) {
	devices := []client.Device{
		{Id: "id1", Name: "Test1", Enabled: true},
		{Id: "id2", Name: "Test2", Enabled: false},
	}
	stub := APIClientStub{}
	stub.enable = func(id string) error {return nil}
	man := manager.DeviceManager{Client: stub, Devices: devices}

	err := man.Enable(0)
	if err != nil {
		t.Error("assert nil error, got", err)
	}

	if !man.Devices[0].Enabled{
		t.Error("assert device 0 enabled")
	}

	err = man.Enable(1)
	if err != nil {
		t.Error("assert nil error, got", err)
	}

	if !man.Devices[1].Enabled{
		t.Error("assert device 1 enabled")
	}
}

func TestDisable(t *testing.T) {
	devices := []client.Device{
		{Id: "id1", Name: "Test1", Enabled: true},
		{Id: "id2", Name: "Test2", Enabled: false},
	}
	stub := APIClientStub{}
	stub.disable = func(id string) error {return nil}
	man := manager.DeviceManager{Client: stub, Devices: devices}

	err := man.Disable(0)
	if err != nil {
		t.Error("assert nil error, got", err)
	}

	if man.Devices[0].Enabled{
		t.Error("assert device 0 disabled")
	}

	err = man.Disable(1)
	if err != nil {
		t.Error("assert nil error, got", err)
	}

	if man.Devices[1].Enabled{
		t.Error("assert device 1 disabled")
	}
}

func TestCreate(t *testing.T) {
	devices := []client.Device{
		{Id: "id1", Name: "Test1", Enabled: true},
		{Id: "id2", Name: "Test2", Enabled: false},
	}
	newDevice := client.Device{Id: "id2", Name: "Test3", Enabled: true}
	stub := APIClientStub{}
	stub.create = func(name string) (client.Device, error) {return newDevice, nil}
	man := manager.DeviceManager{Client: stub, Devices: devices}

	err := man.Create("Test3")
	if err != nil {
		t.Error("assert nil error, got", err)
	}

	if len(man.Devices) != 3 {
		t.Error("assert 3 devices, got", len(man.Devices))
	}

	if man.Devices[2] != newDevice {
		t.Error("assert newDevice is last, got", man.Devices[2])
	}
}

func TestDelete(t *testing.T) {
	devices := []client.Device{
		{Id: "id1", Name: "Test1", Enabled: true},
		{Id: "id2", Name: "Test2", Enabled: false},
	}
	stub := APIClientStub{}
	stub.delete = func(id string) error {return nil}
	man := manager.DeviceManager{Client: stub, Devices: devices}

	err := man.Delete(0)
	if err != nil {
		t.Error("assert nil error, got", err)
	}

	err = man.Delete(0)
	if err != nil {
		t.Error("assert nil error, got", err)
	}

	if len(man.Devices) != 0 {
		t.Error("assert devices empty, got", len(man.Devices))
	}
}
