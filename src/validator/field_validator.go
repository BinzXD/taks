package validator

type FieldCreateRequest struct {
	Name     string `json:"name" validate:"required"`
	Price    int    `json:"price" validate:"required"`
	Location string `json:"location" validate:"required"`
}

type FieldUpdateRequest struct {
	Name     string `json:"name" validate:"required,min=5,max=100"`
	Price    int    `json:"price" validate:"required,gte=0"`
	Location string `json:"location" validate:"required,min=5,max=200"`
}
