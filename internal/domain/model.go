package domain

import "time"

type Plan struct {
	DefaultHoldMinutes int        `json:"default_hold_minutes"`
	Routes             []Route    `json:"routes"`
	Shipments          []Shipment `json:"shipments"`
}

type Route struct {
	ID       string `json:"id"`
	Capacity int    `json:"capacity"`
}

type Shipment struct {
	ID        string     `json:"id"`
	RouteID   string     `json:"route_id"`
	Labels    []string   `json:"labels"`
	HoldUntil *time.Time `json:"hold_until,omitempty"`
}

type Assignment struct {
	ShipmentID string    `json:"shipment_id"`
	RouteID    string    `json:"route_id"`
	State      string    `json:"state"`
	ReadyAt    time.Time `json:"ready_at"`
	Labels     []string  `json:"labels"`
}

type Rejection struct {
	ShipmentID string `json:"shipment_id"`
	Reason     string `json:"reason"`
}

type RouteLoad struct {
	RouteID  string `json:"route_id"`
	Used     int    `json:"used"`
	Capacity int    `json:"capacity"`
}

type Report struct {
	Assignments []Assignment `json:"assignments"`
	Rejections  []Rejection  `json:"rejections"`
	RouteLoads  []RouteLoad  `json:"route_loads"`
}
