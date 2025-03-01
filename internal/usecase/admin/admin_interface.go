package admin

type AdminUseCase interface {
	CreateAdmin(req *CreateAdminRequest) error
}
