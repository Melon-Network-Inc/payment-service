package processor

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetUsername(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"test_getUsername", args{
			c: CreateContextWithUsername("testaccount"),
		}, "testaccount"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetUsername(tt.args.c); got != tt.want {
				t.Errorf("GetUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestGetUserID(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"test_getUserID", args{
			c: CreateContextWithUserID("123"),
		}, "123"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetUserID(tt.args.c); got != tt.want {
				t.Errorf("GetUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestGetRole(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"test_getRole", args{
			c: CreateContextWithRole("1"),
		}, "1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRole(tt.args.c); got != tt.want {
				t.Errorf("GetRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func CreateContextWithRole(role string) *gin.Context {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	ctx.Request.Header.Set(ContextRoleKey, role)
	return ctx
}

func CreateContextWithUserID(userID string) *gin.Context {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	ctx.Request.Header.Set(ContextUserIdKey, userID)
	return ctx
}

func CreateContextWithUsername(username string) *gin.Context {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	ctx.Request.Header.Set(ContextUsernameKey, username)
	return ctx
}