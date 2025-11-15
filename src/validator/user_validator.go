package validator

type UserCreateRequest struct {
	RoleID   uint   `json:"role_id" validate:"required"`
	Fullname string `json:"fullname" validate:"required,min=3,max=100"`
	Username string `json:"username" validate:"required,min=3,max=25"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserUpdateRequest struct {
	RoleID   uint   `json:"role_id" validate:"required"`
	Fullname string `json:"fullname" validate:"required,min=3"`
}
