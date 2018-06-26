# cuxs/mailer

[![build status](https://git.qasico.com/cuxs/mailer/badges/master/build.svg)](https://git.qasico.com/cuxs/mailer/commits/master) [![coverage report](https://git.qasico.com/cuxs/mailer/badges/master/coverage.svg)](https://git.qasico.com/cuxs/mailer/commits/master)

cuxs/mailer is a simple and efficient package to send emails.
it can only send emails using an SMTP server. But the API is flexible and it
is easy to implement other methods for sending emails using a local Postfix, an
API, etc.

## Features
- Attachments
- Embedded images
- HTML and text templates
- Automatic encoding of special characters
- SSL and TLS
- Sending multiple emails with the same SMTP connection

## Installation
```
    go get git.qasico.com/cuxs/mailer
```

and set environment variable for
```
// host of smtp server
SMTP_HOST=
// port number smtp
SMTP_PORT=
// authentication credential
SMTP_USERNAME=
SMTP_PASSWORD=
// from where email will appear in receipient
SMTP_SENDER=
SMTP_SENDER_NAME=
```

## Example Usage
```go
    package main
    
    import (
        "git.qasico.com/cuxs/mailer"
    )
    
    func main() {
        m := mailer.NewMessage()
        m.SetRecipient("bob@example.com", "cora@example.com")
        m.SetAddressHeader("Cc", "dan@example.com", "Dan")
        m.SetSubject("Hello!")
        m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")
        m.Attach("/home/Alex/lolcat.jpg")
        
        // initialing smtp dialer
        d := mailer.NewDialer()
        // Send the email to Bob, Cora and Dan.
        if err := d.DialAndSend(m); err != nil {
            panic(err)
        }
    }
```

### Sending multiple receipient
To set multiple receipient is easy, just type many receipient email
```go
    package main
        
        import (
            "git.qasico.com/cuxs/mailer"
        )
        
        func main() {
            m := mailer.NewMessage()
            m.SetRecipient("one@example.com", "two@example.com", "three@example.com")
        }
```

### Using email and name receipient
When you want to set name of receipient, you need to formating it first using
`FormatAddress()` method. it will append email and address on 1 receipient.
 ```go
     package main
         
         import (
             "git.qasico.com/cuxs/mailer"
         )
         
         func main() {
             m := mailer.NewMessage()
             m.SetRecipient(m.FormatAddress("one@example.com", "Mr One"), m.FormatAddress("two@example.com", "Mr Two"), m.FormatAddress("three@example.com", "Mr Three"))
         }
 ```
 
### Using HTML template and variable
cuxs/mailer juga mempunyai fungsi bawaan untuk mengcompile file HTML template dengan variable untuk 
dijadikan sebagai email body pada messages.

- Template file stored in /template/hi.html
```html
<html>
    <body>
        Hi, {{.Name}}
    </body>
</html>
```

- main file
```go
    package main
    
    import (
        "git.qasico.com/cuxs/mailer"
    )
    
    func main() {
        // struct as variable that will be used as data source on template file
        // you are free to use any struct to give a template an data.
        data := struct{ 
            Name string 
        }{
            Name: "Testing",
        }
        
        tpl := mailer.ParseTemplate("/template/hi.html")
    
        // initialing email messages
        m := mailer.NewMessage()
        m.SetRecipient(m.FormatAddress("one@example.com", "Mr One"))
        m.SetSubject("Hello!")
        
        // set body from template and compile with data source
        m.SetBody("text/html", m.FormatHTML(tpl, data))
        
        // initialing smtp dialer
        d := mailer.NewDialer()
        // Send the email to Bob, Cora and Dan.
        if err := d.DialAndSend(m); err != nil {
            panic(err)
        }
    }
```