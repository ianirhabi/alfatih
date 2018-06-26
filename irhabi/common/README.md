# cuxs/common

[![build status](https://git.qasico.com/cuxs/common/badges/master/build.svg)](https://git.qasico.com/cuxs/common/commits/master) [![coverage report](https://git.qasico.com/cuxs/common/badges/master/coverage.svg)](https://git.qasico.com/cuxs/common/commits/master)

## Feature Overview
- [X] Logger using logrus
- [X] Pretend, tester packages to pretend http requests
- [X] Converter functions
- [X] Go common functions
- [x] Now time package helper (ripped from jinzhu)

## Contribute

**Use issues for everything**

- Report issues
- Discuss on chat before sending a pull request
- Suggest new features or enhancements
- Improve/fix documentation

## Installation
```bash
- go get git.qasico.com/cuxs/common
```

## Func common
```bash
- func RandomStr(n int) // RandomStr return random string with defined length.
- func Encrypt(n interface{}) // Encrypt perform simple encryption and decription values.
- func Decrypt(v interface{}) // Decrypt return real values of encripted values.
- func PasswordHash(hashed string, plain string) // PasswordHash compares hashed password with its possible
- func PasswordHasher(p string) // PasswordHasher returns the bcrypt hash of the password
```

### Basic usage RandomStr
```bash
	var length int
	var strings string
	// define string length is 5
	length=5
	//RandomStr will generate a strin from letterbyte of "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	//generate a random string
	strings= common.RandomStr(length)
	fmt.Println(strings)
```

### Basic usage Encrypt
```bash
	var test_int int
	var test_string_num string
	var test_string string
	test_int=1
	test_string_num="1"
	test_string="test"

	encrypted_int:= common.Encrypt(test_int)
	encrypted_num:= common.Encrypt(test_string_num)
	encrypted_string:= common.Encrypt(test_string)
	//expected '65536' '65536' '0'
	fmt.Println(encrypted_int,encrypted_num,encrypted_string)
```

### Basic usage Decrypt
```bash
	test:=10
	encrypted:= common.Encrypt(test)
	//expected '655360'
	fmt.Println(encrypted)
	Decrypted,_:= common.Decrypt(encrypted)
	//expected '10'
	fmt.Println(Decrypted)
```

### Basic usage  PasswordHash and PasswordHasher
```bash
	test:="password"
	pw:= "password"
	pws:="password1"
	//PasswordHasher returns the bcrypt hash of the password
	pwd, err:= common.PasswordHasher(test)
	//expected $2a$10$Gm6bZaazPAPZS55XqEeaBOgxD1Tvuq84rBZp6uQbkPF8qexgewtEm <nil>
	fmt.Println(pwd,err)
	//PasswordHash compares hashed password with its possible
	e:=common.PasswordHash(pwd,pw)
	//expected <nil>
	fmt.Println(e)
	//PasswordHash compares hashed password with its possible
	er:=common.PasswordHash(pwd,pws)
	//expected crypto/bcrypt: hashedPassword is not the hash of the given password
	fmt.Println(er)	
```

## Func converter
```bash
- func ToString(value interface{}) //convert the input to a string.
- func ToJSON(value interface{}) // convert the input to a valid JSON string
- func ToFloat(value interface{}) //convert the input string to a float, or 0.0 if the input is not a float.
- func ToInt(value interface{}) //convert the input string to an integer, or 0 if the input is not an integer.
- func ToBoolean(value interface{}) //convert the input string to a boolean.
- func ToLower(value interface{}) //convert the value string into lowercase format.
- func ToUpper(value interface{}) //convert the value string into uppercase format.
- func ToCamelCase(s string) //converts from underscore separated form to camel case form.
- func ToUnderscore(str string) //converts from camel case form to underscore separated form.
- func LeftTrim(str, chars string) //trim characters from the left-side of the input.
- func RightTrim(str, chars string) //trim characters from the right-side of the input.
- func Trim(str, chars string) //trim characters from both sides of the input.
```

### Basic usage  ToString
```bash
	y:=common.ToString("str123")
	c:=common.ToString(123)
	d:=common.ToString(12.3)
	f:=common.ToString(1.5+10i)
	//expected "str123","123","12.3","1.5+10i"
	fmt.Println(y,c,d,f)
```

### Basic usage  Toint
```bash
	tests := []string{"1000", "-123", "abcdef", "100000000000000000000000000000000000000000000", " 1"}
	//expected {1000, -123, 0, 0, 1}
	for _,i:=range tests{
		// convert string to int
		c:=common.ToInt(i)
		fmt.Println(c)
	}
	//convert float to int
	d:=common.ToInt(0.1)
	//expected 0
	fmt.Println(d)
```

### Basic usage  ToJSON
```bash
	str:="test"
	st:=common.ToJSON(str)
```

### Basic usage  ToFloat
```bash
	test:=[]interface{} {"", "123", "-.01", "10.", "string", "1.23e3", ".23e10", []string{"asd"}, 0.1,76}

	for _,i:=range test {
		fl:=common.ToFloat(i)
		//expected '0' '123' '-0.01' '10' '0' '1230' '2.3e+09' '0' '0.1' '76'
		fmt.Println(fl)
	}
```

### Basic usage  ToBoolean
```bash
	test := []string{"true", "1", "True", "false", "0", "abcdef"}

	for _,i:=range test {
		bl:=common.ToBoolean(i)
		//expected true, true, true, false, false, false
		fmt.Println(bl)
	}
```

### Basic usage  ToLower
```bash
	test := []string{"A_B_C", "MY_FUNC", "1AB_CD", "0", "abcDef"}

	for _,i:=range test {
		lw:=common.ToLower(i)
		//expected "a_b_c", "my_func", "1ab_cd", "0","abcdf"
		fmt.Println(lw)
	}
```

### Basic usage  ToUpper
```bash
	test := []string{"a_b_c", "my_func", "1ab_cd", "0","abcdf"}

	for _,i:=range test {
		up:=common.ToUpper(i)
		//expected "A_B_C", "MY_FUNC", "1AB_CD", "0", "abcDef"
		fmt.Println(up)
	}
```

### Basic usage  ToCamelCase
```bash
	test := []string{"a_b_c", "my_func", "1ab_cd", "0","ab_cd f"}
	for _,i:=range test {
		cl:=common.ToCamelCase(i)
		//expected "ABC", "MyFunc", "1abCd", "0", "AbCdF"
		fmt.Println(cl)
	}
```

### Basic usage  ToUnderscore
```bash
	test := []string{"ABC", "MyFunc", "1abCd", "0", "AbCdF"}
	for _,i:=range test {
		us:=common.ToUnderscore(i)
		//expected "a_b_c", "my_func", "1ab_cd", "0","ab_cd_f"
		fmt.Println(us)
	}
```

### Basic usage  LeftTrim
```bash
	lft:="  rrrrfoossss   "
	str:=common.LeftTrim(lft,"  rrrr")
	//expected foossss
	fmt.Println(str)
```

### Basic usage  RightTrim
```bash
	rght:="  rrrrfoossss   "
	str:=common.RightTrim(rght,"  ssss")
	//expected rrrrfoo
	fmt.Println(str)
```

### Basic usage  Trim
```bash
	trm:="1234567890987654321"
	str:=common.Trim(trm,"1-8")
	//expected "909"
	fmt.Println(str)
```

## Package tester

### Struct tester
```bash
// H is HTTP Header Type
type H map[string]string

// D is HTTP Data Type
type D map[string]interface{}

// RequestConfig provide user input request structure
type RequestConfig struct {
	Method  string
	Path    string
	Body    string
	Headers H
	Cookies H
	Debug   bool
}
```

### Func in tester
```bash
-New() // New supply initial structure
-SetDebug(enable bool) // SetDebug supply enable debug mode.
- GET(path string) // GET is request method.
- POST(path string) // POST is request method.
- PUT(path string) // PUT is request method.
- DELETE(path string) // DELETE is request method.
- PATCH(path string) // PATCH is request method.
- HEAD(path string) // HEAD is request method.
- OPTIONS(path string) // OPTIONS is request method.
- SetHeader(headers H) // SetHeader supply http header what you defined.
- SetJSON(body D) // SetJSON supply JSON body.
- SetForm(body H) // SetForm supply form body.
- SetQuery(query H) // SetQuery supply query string.
- SetBody(body string) // SetBody supply raw body.
- SetCookie(cookies H) // SetCookie supply cookies what you defined.
- Run(r http.Handler, response ResponseFunc) // Run execute http request
```

### Example
```bash
type echoContent struct {
	Hello string `json:"hello"`
	Foo   string `json:"foo"`
	A     string `json:"a"`
	B     string `json:"b"`
	C     string `json:"c"`
	D     string `json:"d"`
}

// Binding from JSON
type echoJSONContent struct {
	A int `json:"a" binding:"required"`
	B int `json:"b" binding:"required"`
}


func echoQueryHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		text := c.QueryParam("text")
		foo := c.QueryParam("foo")

		return c.JSON(http.StatusOK, &echoContent{
			Hello: text,
			Foo:   foo,
		})
	}
}

func echoPostFormHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		a := c.FormValue("a")
		b := c.FormValue("b")

		return c.JSON(http.StatusOK, &echoContent{
			A: a,
			B: b,
		})
	}
}

func echoJSONHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		json := new(echoJSONContent)
		err := c.Bind(json)

		if err != nil {
			log.Error(err)
		}

		return c.JSON(http.StatusOK, json)
	}
}
// EchoEngine is echo router.
func engine() *echo.Echo {
	e := echo.New()
	
	e.GET("/query", echoQueryHandler())
	e.POST("/form", echoPostFormHandler())
	e.POST("/json", echoJSONHandler())
	return e
}
func EchoQuery(t *testing.T) {
	// New supply initial structure
	r := New()
	// GET is request method.
	r.GET("/query?text=world&foo=bar").
		// Run execute http request
		Run(engine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())

			hello, _ := jsonparser.GetString(data, "hello")
			foo, _ := jsonparser.GetString(data, "foo")
			//expected get result string "world" and "bar"
			assert.Equal(t, "world", hello)
			assert.Equal(t, "bar", foo)
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func EchoPostFormData(t *testing.T) {
	// New supply initial structure
	r := New()
	// POST is request method.
	r.POST("/form").
		// SetBody supply raw body.
		SetBody("a=1&b=2").
		// Run execute http request
		Run(engine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())

			a, _ := jsonparser.GetString(data, "a")
			b, _ := jsonparser.GetString(data, "b")
			//expected result "1" and "2" from post method
			assert.Equal(t, "1", a)
			assert.Equal(t, "2", b)
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func EchoPostJSONData(t *testing.T) {
	r := New()
	r.POST("/json").
		// SetJSON supply JSON body.
		SetJSON(D{
			"a": 1,
			"b": 2,
		}).
		Run(engine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())

			a, _ := jsonparser.GetInt(data, "a")
			b, _ := jsonparser.GetInt(data, "b")
			//expected result 1 and 2 from POST with JSON
			assert.Equal(t, 1, int(a))
			assert.Equal(t, 2, int(b))
			assert.Equal(t, http.StatusOK, r.Code)
		})
}
```

### Basic usage Now

```go
time.Now() // 2013-11-18 17:51:49.123456789 Mon

now.BeginningOfMinute()   // 2013-11-18 17:51:00 Mon
now.BeginningOfHour()     // 2013-11-18 17:00:00 Mon
now.BeginningOfDay()      // 2013-11-18 00:00:00 Mon
now.BeginningOfWeek()     // 2013-11-17 00:00:00 Sun
now.FirstDayMonday = true // Set Monday as first day, default is Sunday
now.BeginningOfWeek()     // 2013-11-18 00:00:00 Mon
now.BeginningOfMonth()    // 2013-11-01 00:00:00 Fri
now.BeginningOfQuarter()  // 2013-10-01 00:00:00 Tue
now.BeginningOfYear()     // 2013-01-01 00:00:00 Tue

now.EndOfMinute()         // 2013-11-18 17:51:59.999999999 Mon
now.EndOfHour()           // 2013-11-18 17:59:59.999999999 Mon
now.EndOfDay()            // 2013-11-18 23:59:59.999999999 Mon
now.EndOfWeek()           // 2013-11-23 23:59:59.999999999 Sat
now.FirstDayMonday = true // Set Monday as first day, default is Sunday
now.EndOfWeek()           // 2013-11-24 23:59:59.999999999 Sun
now.EndOfMonth()          // 2013-11-30 23:59:59.999999999 Sat
now.EndOfQuarter()        // 2013-12-31 23:59:59.999999999 Tue
now.EndOfYear()           // 2013-12-31 23:59:59.999999999 Tue


// Use another time
t := time.Date(2013, 02, 18, 17, 51, 49, 123456789, time.Now().Location())
now.New(t).EndOfMonth()   // 2013-02-28 23:59:59.999999999 Thu


// Don't want be bothered with the First Day setting, Use Monday, Sunday
now.Monday()              // 2013-11-18 00:00:00 Mon
now.Sunday()              // 2013-11-24 00:00:00 Sun (Next Sunday)
now.EndOfSunday()         // 2013-11-24 23:59:59.999999999 Sun (End of next Sunday)

t := time.Date(2013, 11, 24, 17, 51, 49, 123456789, time.Now().Location()) // 2013-11-24 17:51:49.123456789 Sun
now.New(t).Monday()       // 2013-11-18 00:00:00 Sun (Last Monday if today is Sunday)
now.New(t).Sunday()       // 2013-11-24 00:00:00 Sun (Beginning Of Today if today is Sunday)
now.New(t).EndOfSunday()  // 2013-11-24 23:59:59.999999999 Sun (End of Today if today is Sunday)
```

#### Parse String

```go
time.Now() // 2013-11-18 17:51:49.123456789 Mon

// Parse(string) (time.Time, error)
t, err := now.Parse("12:20")            // 2013-11-18 12:20:00, nil
t, err := now.Parse("1999-12-12 12:20") // 1999-12-12 12:20:00, nil
t, err := now.Parse("99:99")            // 2013-11-18 12:20:00, Can't parse string as time: 99:99

// MustParse(string) time.Time
now.MustParse("2013-01-13")             // 2013-01-13 00:00:00
now.MustParse("02-17")                  // 2013-02-17 00:00:00
now.MustParse("2-17")                   // 2013-02-17 00:00:00
now.MustParse("8")                      // 2013-11-18 08:00:00
now.MustParse("2002-10-12 22:14")       // 2002-10-12 22:14:00
now.MustParse("99:99")                  // panic: Can't parse string as time: 99:99
```

Extend `now` to support more formats is quite easy, just update `TimeFormats` variable with `time.Format` like time layout

```go
now.TimeFormats = append(now.TimeFormats, "02 Jan 2006 15:04")
```
