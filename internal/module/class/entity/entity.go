package entity

type CreateClassRequest struct {
    UserId      string `json:"user_id" validate:"required"` 
    Title       string `json:"title" validate:"required"`   
    Description string `json:"description"`                 
    Image       string `json:"image"`                      
    Video       string `json:"video"`                      
    Status      string `json:"status" validate:"required,oneof=public draf"` 
}

type CreateClassResponse struct {
    Id          int    `json:"id"`           
    CreatorId   string `json:"creator_id"`   
    Title       string `json:"title"`       
    Status      string `json:"status"`      
    CreatedAt   string `json:"created_at"`   
    UpdatedAt   string `json:"updated_at"`   
}
