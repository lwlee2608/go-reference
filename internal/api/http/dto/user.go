package dto

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lwlee2608/go-reference/internal/db/sqlc"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	FullName string `json:"full_name"`
}

type UpdateUserRequest struct {
	FullName string `json:"full_name"`
}

type UserResponse struct {
	ID        string        `json:"id"`
	Username  string        `json:"username"`
	FullName  string        `json:"full_name,omitempty"`
	Role      sqlc.UserRole `json:"role"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

func NewUserResponse(u sqlc.User) UserResponse {
	resp := UserResponse{
		ID:        pgUUIDString(u.ID),
		Username:  u.Username,
		Role:      u.Role,
		CreatedAt: u.CreatedAt.Time,
		UpdatedAt: u.UpdatedAt.Time,
	}
	if u.FullName.Valid {
		resp.FullName = u.FullName.String
	}
	return resp
}

func pgUUIDString(id pgtype.UUID) string {
	if !id.Valid {
		return ""
	}
	b, _ := id.MarshalJSON()
	if len(b) >= 2 {
		return string(b[1 : len(b)-1])
	}
	return ""
}
