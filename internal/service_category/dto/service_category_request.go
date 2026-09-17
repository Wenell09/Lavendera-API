package dto

type CreateServiceCategoryRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}

type UpdateServiceCategoryRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}
