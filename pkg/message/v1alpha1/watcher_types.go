package message

type NotificationManagerResponse struct {
	Status  int    `json:"Status"`
	Message string `json:"Message"`
}

type NotificationManagerRequest struct {
	Receiver *Receiver `json:"receiver"`
	Alert    *struct {
		Alerts Alerts `json:"alerts"`
	} `json:"alert"`
}
