package controller

import (
	"bluebell/models"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSignUpRequestBinding(t *testing.T) {
	ensureTestTranslator(t)
	for _, tc := range []struct {
		name string
		body string
		want ResCode
	}{
		{"matching passwords", `{"username":"signup-test","password":"same-password","re_password":"same-password"}`, CodeSucess},
		{"mismatched passwords", `{"username":"signup-test","password":"same-password","re_password":"different-password"}`, CodeInvalidParam},
		{"wrong confirmation field", `{"username":"signup-test","password":"same-password","confirm_password":"same-password"}`, CodeInvalidParam},
		{"missing confirmation", `{"username":"signup-test","password":"same-password"}`, CodeInvalidParam},
	} {
		t.Run(tc.name, func(t *testing.T) {
			response := performRequest(func(c *gin.Context) {
				var params models.ParamSignUp
				if err := c.ShouldBindJSON(&params); err != nil {
					ResponseError(c, CodeInvalidParam)
					return
				}
				ResponseSucess(c, nil)
			}, http.MethodPost, "/api/v1/signup", tc.body)
			if got := decodeResponse(t, response).Code; got != tc.want {
				t.Fatalf("expected code %d, got %d", tc.want, got)
			}
		})
	}
}
