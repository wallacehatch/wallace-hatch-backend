package easypost

import "github.com/gregm711/easypost-go/fake"

type RequestControllerFake struct{}

func (rif RequestControllerFake) do(method string, objectType string, objectUrl string, payload string) ([]byte, error) {
	//Address
	if objectType == "address" {
		//Address.Create()
		if objectUrl == "" {
			return []byte(fake.EasypostFakeAddress), nil
		}
		//Address.Get()
		if objectUrl != "" {
			return []byte(fake.EasypostFakeAddress), nil
		}
	}

	//Shipment
	if objectType == "shipment" {
		//Shipment.Create()
		if objectUrl == "" {
			return []byte(fake.ShipmentCreate), nil
		}
		//Shipment.Buy()
		if objectUrl != "" {
			return []byte(fake.ShipmentBuy), nil
		}
	}

	if objectType == "order" {
		//Order.Create()
		if objectUrl == "" {
			return []byte(fake.EasypostFakeCreateOrder), nil
		}
		//Order.Buy()
		if objectUrl != "" {
			return []byte(fake.EasypostFakeBuyOrder), nil
		}
	}

	if objectType == "customs_info" {
		if objectUrl == "" {
			return []byte(fake.EasypostFakeCustoms), nil
		}
	}
	return nil, nil
}
