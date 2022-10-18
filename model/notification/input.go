package notification

type InputNotification struct {
	Status string `json:"status" validate:"required"`
	Type string `json:"type" validate:"required"`
	Period string `json:"period" validate:"required"`
}
