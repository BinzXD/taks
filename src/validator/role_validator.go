package validator

type RoleCreateRequest struct {
	Name   string `json:"name" validate:"required"`
	Status bool   `json:"status" validate:"required"`
}

type RoleUpdateRequest struct {
	Name string `json:"name" validate:"required"`
}
