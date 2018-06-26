package irhabi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alfatih/irhabi/orm"
	"github.com/alfatih/irhabi/validation"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func fakeContext(method string, endpoint string, response string, rec *httptest.ResponseRecorder) (*Context, error) {
	e := New()

	req, err := http.NewRequest(method, endpoint, strings.NewReader(response))
	ctx := e.NewContext(req, rec)
	c := NewContext(ctx)

	return c, err
}

func TestNewResponse(t *testing.T) {
	r := NewResponse()

	assert.Equal(t, HTTPResponseSuccess, r.Status)
}

func TestContextData(t *testing.T) {
	type user struct {
		ID   int    `json:"id" xml:"id" form:"id"`
		Name string `json:"name" xml:"name" form:"name"`
	}

	js := `{"status":"success","data":{"id":1,"name":"Jon Snow"},"total":20}`

	rec := httptest.NewRecorder()
	ctx, err := fakeContext(echo.POST, "/", js, rec)
	ctx.Data(user{1, "Jon Snow"}, 20)

	err = ctx.Serve(err)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, echo.MIMEApplicationJSONCharsetUTF8, rec.Header().Get(echo.HeaderContentType))
		assert.Equal(t, js, rec.Body.String())
	}

	rec = httptest.NewRecorder()
	ctx, err = fakeContext(echo.OPTIONS, "/", "", rec)
	err = ctx.Serve(err)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
		assert.Equal(t, "", rec.Header().Get(echo.HeaderContentType))
		assert.Equal(t, "", rec.Body.String())
	}
}

func TestContextResponseErrorType(t *testing.T) {
	js := `{"status":"fail","message":"Unprocessable Entity","errors":{"name":"The name field is required."}}`

	// validation error
	rec := httptest.NewRecorder()
	ctx, _ := fakeContext(echo.POST, "/", js, rec)

	type user struct {
		Name string `json:"name" valid:"required"`
	}

	v := validation.New()
	o := v.Struct(user{""})

	err := ctx.Serve(o)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		assert.Equal(t, echo.MIMEApplicationJSONCharsetUTF8, rec.Header().Get(echo.HeaderContentType))
		assert.Equal(t, js, rec.Body.String())
	}

	// http error
	rec = httptest.NewRecorder()
	ctx, _ = fakeContext(echo.POST, "/", "", rec)
	err = ctx.Serve(echo.ErrNotFound)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, echo.MIMEApplicationJSONCharsetUTF8, rec.Header().Get(echo.HeaderContentType))
		assert.Equal(t, `{"status":"fail","message":"Not Found"}`, rec.Body.String())
	}

	// orm error
	rec = httptest.NewRecorder()
	ctx, _ = fakeContext(echo.POST, "/", "", rec)
	err = ctx.Serve(orm.NewOrmError("fatal"))
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, echo.MIMEApplicationJSONCharsetUTF8, rec.Header().Get(echo.HeaderContentType))
		assert.Equal(t, `{"status":"fail","message":"Bad Request"}`, rec.Body.String())
	}

	// orm error
	rec = httptest.NewRecorder()
	ctx, _ = fakeContext(echo.POST, "/", "", rec)
	err = ctx.Serve(ErrDataNotExists("name", "The name field is required."))
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		assert.Equal(t, echo.MIMEApplicationJSONCharsetUTF8, rec.Header().Get(echo.HeaderContentType))
		assert.Equal(t, `{"status":"fail","message":"Unprocessable Entity","errors":{"name":"The name field is required."}}`, rec.Body.String())
	}

	rec = httptest.NewRecorder()
	ctx, _ = fakeContext(echo.POST, "/", "", rec)
	err = ctx.Serve(ErrDataExists("name", "The name field is required."))
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		assert.Equal(t, echo.MIMEApplicationJSONCharsetUTF8, rec.Header().Get(echo.HeaderContentType))
		assert.Equal(t, `{"status":"fail","message":"Unprocessable Entity","errors":{"name":"The name field is required."}}`, rec.Body.String())
	}
}

func TestContextErrorShouldPresistences(t *testing.T) {
	type user struct {
		ID   int    `json:"id" xml:"id" form:"id"`
		Name string `json:"name" xml:"name" form:"name"`
	}

	js := `{"status":"fail","message":"Unprocessable Entity","errors":{"name":"The name field is required."}}`

	rec := httptest.NewRecorder()
	ctx, err := fakeContext(echo.POST, "/", js, rec)
	ctx.Data(user{1, "Jon Snow"}, 20)

	ctx.Failure("name", "The name field is required.")

	err = ctx.Serve(err)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		assert.Equal(t, echo.MIMEApplicationJSONCharsetUTF8, rec.Header().Get(echo.HeaderContentType))
		assert.Equal(t, js, rec.Body.String())
	}
}

func TestContextRequestQuery(t *testing.T) {
	e := New()
	req, _ := http.NewRequest(echo.GET, "/?perpage=100&page=5&fields=name,email&orderby=id,-user.name&embeds=users,entities,usergroups&conditions=Or.item.item_name.icontains%3Aopa%252COr.code.icontains%3Aopa%252COr.opt1.icontains%3Aopa%252COr.opt2.icontains%3Aopa%252COr.opt3.icontains%3Aopa%252COr.opt4.icontains%3Aopa%7Cis_deleted%3A0%7Cis_archived%3A0", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	c := NewContext(ctx)

	qs := c.RequestQuery()
	assert.Equal(t, 100, qs.Limit)
	assert.Equal(t, 400, qs.Offset)
	assert.Equal(t, 2, len(qs.Fields))
	assert.Equal(t, 2, len(qs.OrderBy))
	assert.Equal(t, 3, len(qs.Embeds))
	assert.Equal(t, 3, len(qs.Conditions))
}

func TestContextRequestQueryWithDecryption(t *testing.T) {
	var cases = []struct {
		query    string
		expected []map[string]string
	}{
		{"/?conditions=name:alif%252Cemail:alifamri@qasico.com", []map[string]string{{"name": "alif", "email": "alifamri@qasico.com"}}},
		{"/?conditions=id.e:65536%252Citem.item_category_id.id.e:65536%252Cis_archived:1%252Cidx:1", []map[string]string{{"id": "1", "item.item_category_id.id": "1", "is_archived": "1", "idx": "1"}}},
	}

	e := New()
	for _, i := range cases {
		req, _ := http.NewRequest(echo.GET, i.query, nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		c := NewContext(ctx)

		qs := c.RequestQuery()
		assert.Equal(t, i.expected, qs.Conditions)
	}
}
