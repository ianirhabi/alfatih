package irhabi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestJwtToken(t *testing.T) {
	token := JwtToken("id", 1)
	validAuth := "Bearer " + token

	assert.NotEmpty(t, token, "JWT is not working as expexted.")

	req, _ := http.NewRequest(echo.GET, "/", nil)
	req.Header.Set(echo.HeaderAuthorization, validAuth)
	res := httptest.NewRecorder()

	e := New()
	c := e.NewContext(req, res)
	ctx := NewContext(c)

	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}

	h := Authorized()(handler)
	if assert.NoError(t, h(c), "jwt token invalid middleware handle error") {
		x := &testJwtUser{}
		i := ctx.JwtUsers(x).(*userJwt)

		assert.Equal(t, 1, i.ID)
		assert.Equal(t, "Demo", i.Name)
	}
}

type userJwt struct {
	ID   int    `json:"id" xml:"id" form:"id"`
	Name string `json:"name" xml:"name" form:"name"`
}

type testJwtUser struct{}

func (t *testJwtUser) GetUser(id int64) (interface{}, error) {
	u := &userJwt{ID: 1, Name: "Demo"}
	return u, nil
}
