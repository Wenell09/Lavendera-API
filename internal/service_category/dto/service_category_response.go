package dto

type ServiceCategoryResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type ServiceCategoryListResponse struct {
	Data []ServiceCategoryResponse `json:"data"`
}
