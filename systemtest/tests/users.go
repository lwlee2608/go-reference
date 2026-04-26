package tests

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwlee2608/go-reference/internal/api/http/dto"
)

const usersPath = "/api/v1/users"

func TestUserCRUD(t *testing.T, router *gin.Engine) {
	t.Run("create", func(t *testing.T) { testCreateUser(t, router) })
	t.Run("create_duplicate_returns_409", func(t *testing.T) { testCreateDuplicate(t, router) })
	t.Run("get_not_found", func(t *testing.T) { testGetNotFound(t, router) })
	t.Run("update_partial", func(t *testing.T) { testUpdatePartial(t, router) })
	t.Run("list", func(t *testing.T) { testListUsers(t, router) })
	t.Run("delete", func(t *testing.T) { testDeleteUser(t, router) })
}

func testCreateUser(t *testing.T, router *gin.Engine) {
	rr := DoJSON(t, router, http.MethodPost, usersPath, dto.CreateUserRequest{
		Username: "alice",
		Password: "s3cret!!",
		FullName: "Alice Example",
	})
	AssertStatus(t, rr, http.StatusCreated)

	var resp dto.UserResponse
	DecodeJSON(t, rr, &resp)
	assert.Equal(t, "alice", resp.Username)
	assert.Equal(t, "Alice Example", resp.FullName)
	assert.NotEmpty(t, resp.ID)
}

func testCreateDuplicate(t *testing.T, router *gin.Engine) {
	body := dto.CreateUserRequest{Username: "dupe", Password: "p"}
	rr := DoJSON(t, router, http.MethodPost, usersPath, body)
	AssertStatus(t, rr, http.StatusCreated)

	rr = DoJSON(t, router, http.MethodPost, usersPath, body)
	AssertStatus(t, rr, http.StatusConflict)
}

func testGetNotFound(t *testing.T, router *gin.Engine) {
	rr := DoJSON(t, router, http.MethodGet, usersPath+"/00000000-0000-0000-0000-000000000000", nil)
	AssertStatus(t, rr, http.StatusNotFound)
}

func testUpdatePartial(t *testing.T, router *gin.Engine) {
	created := mustCreateUser(t, router, "bob", "pw", "Bob Original")

	rr := DoJSON(t, router, http.MethodPatch, usersPath+"/"+created.ID, map[string]any{})
	AssertStatus(t, rr, http.StatusOK)
	var unchanged dto.UserResponse
	DecodeJSON(t, rr, &unchanged)
	assert.Equal(t, "Bob Original", unchanged.FullName)

	newName := "Bob Updated"
	rr = DoJSON(t, router, http.MethodPatch, usersPath+"/"+created.ID, map[string]any{
		"full_name": newName,
	})
	AssertStatus(t, rr, http.StatusOK)
	var updated dto.UserResponse
	DecodeJSON(t, rr, &updated)
	assert.Equal(t, newName, updated.FullName)
}

func testListUsers(t *testing.T, router *gin.Engine) {
	rr := DoJSON(t, router, http.MethodGet, usersPath, nil)
	AssertStatus(t, rr, http.StatusOK)

	var users []dto.UserResponse
	DecodeJSON(t, rr, &users)
	assert.NotEmpty(t, users)
}

func testDeleteUser(t *testing.T, router *gin.Engine) {
	created := mustCreateUser(t, router, "carol", "pw", "")

	rr := DoJSON(t, router, http.MethodDelete, usersPath+"/"+created.ID, nil)
	AssertStatus(t, rr, http.StatusNoContent)

	rr = DoJSON(t, router, http.MethodGet, usersPath+"/"+created.ID, nil)
	AssertStatus(t, rr, http.StatusNotFound)

	rr = DoJSON(t, router, http.MethodDelete, usersPath+"/"+created.ID, nil)
	AssertStatus(t, rr, http.StatusNotFound)
}

func mustCreateUser(t *testing.T, router *gin.Engine, username, password, fullName string) dto.UserResponse {
	t.Helper()
	rr := DoJSON(t, router, http.MethodPost, usersPath, dto.CreateUserRequest{
		Username: username,
		Password: password,
		FullName: fullName,
	})
	require.Equal(t, http.StatusCreated, rr.Code, rr.Body.String())
	var resp dto.UserResponse
	DecodeJSON(t, rr, &resp)
	return resp
}
