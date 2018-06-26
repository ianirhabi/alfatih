# cuxs/validation

[![build status](https://git.qasico.com/cuxs/validation/badges/master/build.svg)](https://git.qasico.com/cuxs/validation/commits/master) [![coverage report](https://git.qasico.com/cuxs/validation/badges/master/coverage.svg)](https://git.qasico.com/cuxs/validation/commits/master)


## Installation
```
    go get git.qasico.com/cuxs/validation
```

## Field Validation
```go
    package main
    
    import (
        "fmt"
        "git.qasico.com/cuxs/validation"
    )
    
    func main() {
        e := "notvalid.com"
        
        v := validation.New()
        if r := v.Field(e, "required|email"); !r.Valid {
            fmt.Printf(r.Message(), "email")
        }
    }
```

## Struct Validation
```go
    package main
    
    import (
        "fmt"
        "git.qasico.com/cuxs/validation"
    )
    
    // User contains user information
    type User struct {
        FirstName string     `json:"fname" valid:"required"`
        LastName  string     `json:"lname"`
        Age       uint8      `valid:"range:17,45"`
        Email     string     `valid:"required|email"`
        Addresses []Address  `valid:"required"`
    }
    
    // Address houses a users address information
    type Address struct {
        Street string `valid:"required"`
        City   string `valid:"required|alpha"`
        Planet string `valid:"required"`
        Phone  string `valid:"required"`
    }
    
    func main() {
        v := validation.New()
        
        // build 'User' info, normally posted data etc...
        address := Address{
            Street: "Eavesdown Docks",
            Planet: "Persphone",
            Phone:  "none",
        }
        
        user := User{
            FirstName: "",
            LastName:  "",
            Age:       45,
            Email:     "Badger.Smith@gmail.com",
            Addresses: []Address{address},
        }
        
        if res := v.Struct(user); !res.Valid {
            // can return all failure messages as flatten map
            fmt.Println(res.Messages())
            
            // and can get individual failure messages
            // based on field of struct
            fmt.Println(res.Message("firstname"))
        }
        
        // save to database.
    }

```
## Request Validation
```go
    package main
    
    import (
        "fmt"
        "git.qasico.com/cuxs/validation"
    )
    
    // User contains user information
    type User struct {
        FirstName string    `json:"fname" valid:"required"`
        LastName  string    `json:"lname"`
        Age       uint8     `valid:"required|range:17,45"`
        Email     string    `valid:"required|email"`
        Addresses []Address `valid:"required"`
    }
    
    // Address houses a users address information
    type Address struct {
        Street string `valid:"required"`
        City   string `valid:"required|alpha"`
        Planet string `valid:"required"`
        Phone  string `valid:"required"`
    }
    
    func (u User) Validate() *validation.Output {
        o := &validation.Output{Valid: true}
        
        // you can put custom validation here.
        // checking databases or anything else
        if u.Email == "Badger.Smith@gmail.com" {
            o.Valid = false
            o.Failure("email", "Email is already registered.")
        }
        
        return o
    }
    
    func (u User) Messages() map[string]string {
        // set custom messages here
        return map[string]string{
            "firstname.required":        "Please type your first name.",
            "age.required":              "Please type your age.",
            "age.range":                 "Sorry your age is out of range.",
            "addresses.*.city.required": "You should have a city.",
        }
    }
    
    func main() {
        v := validation.New()
        
        // build 'User' info, normally posted data etc...
        address := Address{
            Street: "Eavesdown Docks",
            Planet: "Persphone",
            Phone:  "none",
        }
        
        user := User{
            FirstName: "",
            LastName:  "",
            Age:       45,
            Email:     "Badger.Smith@gmail.com",
            Addresses: []Address{address},
        }
        
        if res := v.Request(user); !res.Valid {
            // can return all failure messages as flatten map
            fmt.Println(res.Messages())
            
            // and can get individual failure messages
            // based on field of struct
            fmt.Println(res.Message("firstname"))
        }
        
        // save to database.
    }

```

## Value Checker
```go
    package main
    
    import (
        "fmt"
        "git.qasico.com/cuxs/validation"
    )
    
    func main() {
        e := "notvalidemail.com"
        if !validation.IsEmail(e) {
            fmt.Println("its not email")
        }
    }
```

## Available Tag (Validator)
Below is a list of all available validation rules and their function:

#### required
The field under validation must be present in the input data and not empty. A field is considered "empty" if one of the following conditions are true:
```
- The value is `null`.
- The value is an empty string.
- The value is an empty map or empty `Countable` slice.
```
#### numeric
The field under validation must be numeric.

#### alpha
The field under validation must be entirely alphabetic characters.

#### alpha_num
The field under validation must be entirely alpha-numeric characters.

#### alpha_num_space
The field under validation must be entirely alpha-numeric characters with spaces.

#### alpha_space
The field under validation must be entirely alphabetic characters with spaces.

#### email
The field under validation must be formatted as an e-mail address.

#### url
The field under validation must be a valid URL.

#### json
The field under validation must be a valid JSON string.

#### cc
The field under validation must be a valid credit card number format.

#### lte:max
The field under validation must have a maximum value. For string data, value corresponds to the number of characters. For numeric data, value corresponds to a given integer value. For an slice/map, size corresponds to the count of the slice.

#### gte:min
The field under validation must have a minimum value. For string data, value corresponds to the number of characters. For numeric data, value corresponds to a given integer value. For an slice/map, size corresponds to the count of the slice.
 
#### lt:max
The field under validation must lower than maximum value. For string data, value corresponds to the number of characters. For numeric data, value corresponds to a given integer value. For an slice/map, size corresponds to the count of the slice.

#### gt:min
The field under validation must lower than minimum value. For string data, value corresponds to the number of characters. For numeric data, value corresponds to a given integer value. For an slice/map, size corresponds to the count of the slice.

#### range:min,max
The field under validation must have a length between the given min and max. it will evaluate same as lte and gte rules.

#### contains:param
The field under validation must contains the param as substring.

#### match:pattern
The field under validation must match the given regular expression.

#### same:param
The given field must match the field under validation.

#### in:param1,param2...
The field under validation must be included in the given list of values.

#### not_in:param1,param2...
The field under validation must not be included in the given list of values.
