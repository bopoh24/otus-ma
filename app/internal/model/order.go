package model

type ShipmentType string

const (
	ShipmentTypePickup  ShipmentType = "pickup"
	ShipmentTypeCourier ShipmentType = "courier"
	ShipmentTypePost    ShipmentType = "post"
)

type Order struct {
	ID            int64        `json:"id,omitempty"`
	UserID        int64        `json:"user_id,omitempty"`
	Products      []int64      `json:"products"`
	PaymentMethod string       `json:"payment_method"`
	ShipmentType  ShipmentType `json:"shipment_type"`
	ShipmentAddr  string       `json:"shipment_addr"`
}
