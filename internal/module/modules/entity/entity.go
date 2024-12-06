package entity

type CreateModulesRequest struct {
	UserId      string   `validate:"required"`
	MaterialsId int      `json:"materials_id"`
	Title       string   `json:"title" validate:"required"`
	Content     string   `json:"content"`
	Attachments []string `json:"attachments"`
	Videos      []string `json:"videos"`
}

type CreateModulesResponse struct {
	Id          int      `json:"id" db:"id"`
	Title       string   `json:"title" db:"title"`
	Content     string   `json:"content" db:"content"`
	Attachments []string `json:"attachments" db:"attachments"`
	Videos      []string `json:"videos" db:"videos"`
	CreatedAt   string   `json:"created_at" db:"created_at"`
	UpdatedAt   string   `json:"updated_at" db:"updated_at"`
}

type UpdateModulesRequest struct {
	UserId      string   `validate:"required"`
	ModulesId   int      `json:"modules_id"`
	Title       string   `json:"title" validate:"required"`
	Content     string   `json:"content"`
	Attachments []string `json:"attachments"`
	Videos      []string `json:"videos"`
}

type UpdateModulesResponse struct {
	Id          int      `json:"id" db:"id"`
	Title       string   `json:"title" db:"title"`
	Content     string   `json:"content" db:"content"`
	Attachments []string `json:"attachments" db:"attachments"`
	Videos      []string `json:"videos" db:"videos"`
	CreatedAt   string   `json:"created_at" db:"created_at"`
	UpdatedAt   string   `json:"updated_at" db:"updated_at"`
}

type DeleteModulesRequest struct {
	UserId      string   `validate:"required"`
	ModulesId   int      `json:"modules_id"`
}