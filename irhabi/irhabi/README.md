# cuxs/cuxs

[![build status](https://git.qasico.com/cuxs/cuxs/badges/master/build.svg)](https://git.qasico.com/cuxs/cuxs/commits/master) [![coverage report](https://git.qasico.com/cuxs/cuxs/badges/master/coverage.svg)](https://git.qasico.com/cuxs/cuxs/commits/master)

## Feature Overview
- [ ] Auto TLS with letsencrypt
- [X] Read query string for ORM
- [X] Binding with automatic validation 
- [X] Config automatic read from dot env
- [X] Costum context function
- [X] Standartized response formated
- [X] JWT read authorized from context
- [X] Engine using labstack/echo

## Contribute

**Use issues for everything**

- Report issues
- Discuss on chat before sending a pull request
- Suggest new features or enhancements
- Improve/fix documentation

## Installation
```bash
- go get git.qasico.com/cuxs/cuxs
```

## Package mw
- **func HTTPLogger() echo.MiddlewareFunc**<br />
  HTTPLogger returns a middleware that logs HTTP requests.
- **func logRequest(hand echo.HandlerFunc, c echo.Context) (err error)**<br />
  logRequest print all http request on consoles.

## Binder

- **func (v *binderValidator) validate(obj interface{})**<br />
  ValidateStruct evaluate an object,will run validation request if the object is implementing validatonRequests.
- **func (v *binderValidator) lazyinit()**<br />
  lazyinit initialing validator instances for one of time only.
- **func (b binder) Bind(i interface{}, ctx echo.Context)**<br />
  Bind is decode request body and injecting into interfaces, only accept json data type. other type will return error bad requests.
  Also automaticly validate data with interfaces.

## Config

### Configuration variable
```go
type config struct {
	DebugMode  bool   // Switch debug mode for production or development
	JwtSecret  string // Secret key for Json web token algorithm
	GzipEnable bool   // State of gzip compression
	Host       string // IP Application will run, default is 0.0.0.0:8080
	DbEngine   string // Database engines
	DbHost     string // IP Database server, default is 0.0.0.0:3306
	DbName     string // Database name will be used
	DbUser     string // Database username
	DbPassword string // Database password
}
```

### Loadconfig
```go
// loadConfig set config value from environment variable.If not exists, it will have a default values.
func loadConfig() *config {
	c := new(config)
	c.DebugMode = env.GetBool("APP_DEBUGMODE", true)
	c.JwtSecret = env.GetString("APP_JWT_SECRET", "V3ryF*ck1ngS3cur3")
	c.GzipEnable = env.GetBool("APP_GZIP", false)
	c.Host = env.GetString("SERVER_HOST", "0.0.0.0:8080")
	c.DbEngine = env.GetString("DB_ENGINE", "mysql")
	c.DbHost = env.GetString("DB_HOST", "0.0.0.0:3306")
	c.DbName = env.GetString("DB_NAME", "konektifa_app")
	c.DbUser = env.GetString("DB_USER", "root")
	c.DbPassword = env.GetString("DB_PASS", "")/ Context is custom echo.Context
// has defined as middleware.
type Context struct {
	echo.Context
	responseFormat *ResponseFormat
}
	return c
}
// IsDebug returns true if the framework is running in debug mode. set environtment variable to release for disable debug.
func IsDebug() bool {
	return Config.DebugMode
}
```

## Context

### Struct Context
```go
// Context is custom echo.Context
// has defined as middleware.
type Context struct {
	echo.Context
	responseFormat *ResponseFormat
}
// jwtUser model user jwt token interface
// to check is the id given valid as users.
type jwtUser interface {
	GetUser(int64) (interface{}, error)
}
```

### Contex func
- **func NewContext(c echo.Context)**<br />
  NewContext for new instances of context
- **func (c *Context) Data(data interface{}, total ...int64)**<br />
  to set data and total into response format
- **func (c *Context) Failure(fail ...string)**<br />
  to set response format errors
- **func (c *Context) Serve(e error)**<br />
  to response with json data with data that already collected.
- **func (c *Context) RequestQuery()**<br />
  to set query param into orm in the repository
- **func (c *Context) JwtUsers(model jwtUser)**<br />
  to get a user sessions that having jwt token. will request header.

### RequestQuery()
```go
func (c *Context) RequestQuery() *orm.RequestQuery {
	rq := &orm.RequestQuery{
		Offset:     0,
		// if limit <0, there is no limit
		Limit:      -1,
		Conditions: make([]map[string]string, 0),
	}
    //QueryParams will return the query parameters as `url.Values`
    // and RequestQuery will return func ReadFromContext with param url.values
	return rq.ReadFromContext(c.QueryParams())
}
```
RequestQuery is used to set orm parameter that being used in repository.<br />
RequestQuery will take data query string from context and change it to RequestQuery form<br />
Below is Field Struct RequestQuery in ORM:<br />
```go
type RequestQuery struct {
	Conditions []map[string]string
	Fields     []string
	OrderBy    []string
	Embeds     []string //related table field in database
	Offset     int
	Limit      int
}
```

There is some function in cuxs ORM that being used to get a RequestQuery:<br />
- **func (rq *RequestQuery) ReadFromContext(params url.Values) *RequestQuery**<br />
  to read context value in url.values. <br />
- **func (rq *RequestQuery) GetCondition() *Condition** <br />
  to get a query condition based on context query string and condition.<br /> 
  a param value that being read is:"perpage","page","fields","orderby","embeds" and "conditions"<br />
- **func (rq *RequestQuery) GetJoin() interface{}** <br />
  will return related table.<br />

### RequestQuery example test with custom query string in url
```go
//creates an instance of Echo.
e := New()
    
    // create a request with methot get, url string and body nil
	req, _ := http.NewRequest(echo.GET, "/?perpage=100&page=5&fields=name,email&orderby=id,-user.name&embeds=users,entities,usergroups&conditions=name:alif,email:alifamri@qasico.com|AndNot.id:1", nil)
	
	// NewRecorder() returns an initialized ResponseRecorder
	// ResponseRecorder is an implementation of http.ResponseWriter that
    // records its mutations for later inspection in tests.
	rec := httptest.NewRecorder()
	
	// create request and response context
	ctx := e.NewContext(req, rec)
	c := NewContext(ctx)


    // with query string perpage=100, page=5, fields=name,email, orderby=id,-user.name, embeds=user,entities,usergroups and condition= name:alif, email:alifamri@qasico.com|AndNot.id:1
	// perpage=100 and page=5, that's mean there will be 100 data perpage and the total of data is 500 data
	// orderby is used to sort a data, expression - is used to sort by descending and default is ascending
	// embeds meaning table user ,entities and usergroup is related
	// condition where name is alif and email is ... and exclude id=1
	// expression "|" is used as alternative in condition, e.g: hex_digit     = "0" … "9" | "A" … "F" | "a" … "f" .
	// below is some notation expression in go:
	                 //  |   alternation
                     // ()  grouping
                     // []  option (0 or 1 times)
                     // {}  repetition (0 to n times)
                     
    // create a requestquery from context request                 
	qs := c.RequestQuery()
	assert.Equal(t, 100, qs.Limit) // expected qs.limit is 100 (will be limited by perpage)
	assert.Equal(t, 401, qs.Offset) // expected qs.offset is 401 (offset= limit*(page-1)+1)
	assert.Equal(t, 2, len(qs.Fields)) // expected qs.fields is 2 (name and email)
	assert.Equal(t, 2, len(qs.OrderBy)) // expected qs.orderby is 2 (id,-user.name)
	assert.Equal(t, 3, len(qs.Embeds)) // expected qs.embeds is 3 (user,entities,usergroups)
	assert.Equal(t, 2, len(qs.Conditions)) // expected qs.condition is 2 (name:alif, email:alifamri@qasico.com|AndNot.id:1)
```
**Here is an example to get all data using requestquery with method get**
```go
// this program usually located in handler API
    // set request query by the query string context
    ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()
	
	// call a func GetRepo in repository with param requestquery
	// will return all data, count data and err (nil if success)
	data, count, e := repository.GetRepo(rq)
```
```go
// this is a func GetRepo that located in repository
    // this func will return all data based on model whatever, data count and error
    func GetRepo(rq *orm.RequestQuery) (*[]model.Whatever, int64, error) {
        //declare a variable that will act as container for data
    	var m []model.Whatever
    	var count int64
    	var err error
        
        // create a new orm so that orm func can be use
    	o := orm.NewOrm()
    	// make a queryseter (can be string or struct) for table operation with new built in func of model(return *type)
    	// with condition of requestquery
    	qs := o.QueryTable(new(model.Whatever)).SetCond(rq.GetCondition())
    
        // if table is related, than join table
    	if len(rq.Embeds) > 0 {
    		qs = qs.RelatedSel(rq.GetJoin())
    	}
        // Count() will return QuerySeter execution result number
    	if count, err = qs.Count(); err != nil || count == 0 {
    		return nil, count, err
    	}
        // sorting the data with limit and offset from requestquery
    	qs = qs.OrderBy(rq.OrderBy...).Limit(rq.Limit, rq.Offset)
    	//query all data and map to containers m
    	if _, err = qs.All(&m, rq.Fields...); err == nil {
    		return &m, count, nil
    	}
    	return nil, count, err
    }
```

## Cuxs

### Cuxs func
- **func New()**<br />
  creates an instance of Echo.
- **func StartServer(e *echo.Echo)**<br />
  to starting echo servers with routes
- **func HTTPErrorHandler(err error, c echo.Context)**<br /> 
  to invokes the default HTTP error handler.


## Database
```go
// DbSetup registering database connection
// by reading database config variables.
func DbSetup() error {
orm.Debug = IsDebug()
orm.DebugLog = log.Log
orm.DefaultTimeLoc = time.Local
orm.DefaultRelsDepth = 3

ds := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", Config.DbUser, Config.DbPassword, Config.DbHost, Config.DbName, "charset=utf8&loc=Asia%2FJakarta")
return orm.RegisterDataBase("default", Config.DbEngine, ds)
```

## Errors

### Struct Errors
```go
// DataNotExistsError error for inexisting data
   type DataNotExistsError struct {
      Message string
      Errors  map[string]string
   }

// DataDuplicateError error that data is duplicate.
    type DataDuplicateError struct {
      Message string
      Errors  map[string]string
   }
```
### 
```go
// ErrDataNotExists error interface for inexists data
// this will cause error 422 with errors data fail should has 2 string,
// first is the key follow with values
// e := ErrDataNotExists("email", "This email is not exists.")
func ErrDataNotExists(fail ...string) error {
	e := &DataNotExistsError{}
	e.Errors = map[string]string{fail[0]: fail[1]}
	e.Message = fail[1]

	return e
}

// ErrDataExists error interface for already exists data
// this will cause error 422 with errors data fail should has 2 string,
// first is the key follow with values
// e := ErrDataNotExists("email", "This email is not exists.")
func ErrDataExists(fail ...string)error {
	e := &DataDuplicateError{}
	e.Errors = map[string]string{fail[0]: fail[1]}
	e.Message = fail[1]

	return e
}
```

## Response

### Struct Response
```go
const (
	// HTTPResponseSuccess default status for success responses
	HTTPResponseSuccess = "success"

	// HTTPResponseFail default status when responses has an errors.
	HTTPResponseFail = "fail"
)

// ResponseFormat is standart response formater of the applicatin.
type ResponseFormat struct {
	Code    int               `json:"-"`
	Status  string            `json:"status,omitempty"`
	Message interface{}       `json:"message,omitempty"`
	Data    interface{}       `json:"data,omitempty"`
	Total   int64             `json:"total,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}
```

### Func Response
- **func NewResponse()**<br />
  return new instances of response formater.
- **func (r *ResponseFormat) SetData(d interface{}, t ...int64)**<br />
  to fill data and total into response formater.
- **func (r *ResponseFormat) SetError(err error)**<br />
  to set an error into response formater.
- **func (r *ResponseFormat) reset()** <br />
  reset all data in response formater

### Example
```go
     // NewResponse return new instances of response formater.
       a := cuxs.NewResponse()

       // Standard response format
       a.Code = 12345
       a.Status = "succes"
       a.Message = "update succesfully"
       a.Data = "this data"
       a.Total = 10
       a.Errors = make(map[string]string)
       a.Errors["Message"] = "Error"

       fmt.Println( a.Code, "\n", a.Status, "\n", a.Message, "\n", a.Data, "\n", a.Total, "\n", a.Errors["Message"])
}
```

## Utils
for further instruction 
See: https://jwt.io/introduction<br />
See `JWTConfig.TokenLookup`

### Func in Utils
- **func Authorized() echo.MiddlewareFunc**<br />
  Authorized returns a JSON Web Token (JWT) auth middleware. For valid token, it sets the user in context and calls next handler.
  For invalid token, it returns "401 - Unauthorized" error. For empty token, it returns "400 - Bad Request" error.
- **func JwtKey()**<br />
  return byte of jwt secret keys.
- **func JwtToken(k string, v interface{})** <br />
  to generated a JWT token keys and values from the claims. the return will become a valid token with a life time 72 hours from the time generated.
- **func listRoutes(e *echo.Echo)**<br />
  ListRoutes print all route available, only show on debug mode.

### Example Authorized token
```go
type Token struct{}

// URLMapping declare endpoint with handler function.
func (h *Token) URLMapping(r *echo.Group) {
	//check the header with authorized token
	r.GET("", h.get,cuxs.Authorized())
}
```

### Example to generate a token
```go
func main() {
		token := cuxs.JwtToken("id", 1)
		validAuth := "Bearer " + token
// expected Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0ODY5NzQ3NzksImlkIjoxfQ.IUDeVJesTiYWtivBKjGPOKjK1WXUiGxpHNWnRqk1YfI
	fmt.Println(validAuth)
}
```
 :thumbsup:  :thumbsup:  :thumbsup:  :thumbsup:  :thumbsup:  :thumbsup:  :thumbsup:  :thumbsup:
