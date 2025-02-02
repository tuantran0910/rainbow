package dtos

type RegisterUserRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6,max=255"`
	FirstName   string `json:"first_name" binding:"required,min=1,max=255"`
	LastName    string `json:"last_name" binding:"required,min=1,max=255"`
	PhoneNumber string `json:"phone_number" binding:"required,len=10"`
	Role        string `json:"role" binding:"omitempty,oneof=ADMIN USER"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=255"`
}

type GetUserResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	LastLogin   string `json:"last_login"`
	IsActive    bool   `json:"is_active"`
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	DeletedAt   string `json:"deleted_at,omitempty"`
}

type ListUsersResponse struct {
	Users []*GetUserResponse `json:"users"`
}
