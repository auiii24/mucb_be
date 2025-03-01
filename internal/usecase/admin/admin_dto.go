package admin

import "github.com/go-playground/validator/v10"

var validate = validator.New()

type CreateAdminRequest struct {
	Name     string `json:"name" binding:"required,max=64"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=16"`
	Role     string `json:"role" binding:"required,oneof=SUPER_ADMIN ADMIN"`
}

func (req *CreateAdminRequest) Validate() error {
	return validate.Struct(req)
}
