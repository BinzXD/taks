package validator

type BookingCreateRequest struct {
	FieldID   uint   `json:"field_id" validate:"required"`
	StartTime string `json:"start_time" validate:"required"`
	EndTime   string `json:"end_time" validate:"required"`
}

type BookingUpdateRequest struct {
	FieldID   uint   `json:"field_id" validate:"required"`
	Status    string `json:"status" validate:"required,oneof=draft pending done"`
	StartTime string `json:"start_time" validate:"required"`
	EndTime   string `json:"end_time" validate:"required"`
}
