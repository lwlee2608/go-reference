package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lwlee2608/go-reference/internal/api/http/dto"
	"github.com/lwlee2608/go-reference/internal/db/sqlc"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	queries *sqlc.Queries
}

func NewUserHandler(queries *sqlc.Queries) *UserHandler {
	return &UserHandler{queries: queries}
}

func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		_ = c.Error(err)
		return
	}

	user, err := h.queries.CreateUser(c.Request.Context(), sqlc.CreateUserParams{
		Username:     req.Username,
		PasswordHash: string(hash),
		FullName:     pgText(req.FullName),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
			return
		}
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, dto.NewUserResponse(user))
}

func (h *UserHandler) Get(c *gin.Context) {
	id, err := parseUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	user, err := h.queries.GetUser(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.NewUserResponse(user))
}

func (h *UserHandler) List(c *gin.Context) {
	users, err := h.queries.ListUsers(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}

	resp := make([]dto.UserResponse, 0, len(users))
	for _, u := range users {
		resp = append(resp, dto.NewUserResponse(u))
	}
	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) Update(c *gin.Context) {
	id, err := parseUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	current, err := h.queries.GetUser(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		_ = c.Error(err)
		return
	}

	fullName := current.FullName
	if req.FullName != nil {
		fullName = pgText(*req.FullName)
	}

	user, err := h.queries.UpdateUser(c.Request.Context(), sqlc.UpdateUserParams{
		ID:       id,
		FullName: fullName,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.NewUserResponse(user))
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, err := parseUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.queries.DeleteUser(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

func parseUUID(s string) (pgtype.UUID, error) {
	var id pgtype.UUID
	err := id.Scan(s)
	return id, err
}

func pgText(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: s, Valid: true}
}
