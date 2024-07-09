package models

type TranfserMessage struct {
	User   string `json:"user"`
	Type   string `json:"type" default:"transfer"`
	Target string `json:"target"`
	Amount int    `json:"amount"`
}

type UploadMessage struct {
	User       string `json:"user"`
	Type       string `json:"type" default:"upload"`
	Image_uuid string `json:"image_uuid"`
}

type BuyMessage struct {
	User       string `json:"user"`
	Type       string `json:"type" default:"buy"`
	Image_uuid string `json:"image_uuid"`
	Amount     int    `json:"amount"`
}
