package entity

type CreateMaterialsRequest struct {
	UserId      string `validate:"required"`
	ClassId		int 	`json:"class_id"`
	Title       string `json:"title" validate:"required"`
}

type CreateMaterialsResponse struct {
	Id        int    `json:"id" db:"id"`
	CreatorId string `json:"creator_materials_id" db:"creator_materials_id"`
	ClassId   int    `json:"class_id" db:"class_id"`
	Title     string `json:"title" db:"title"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}
