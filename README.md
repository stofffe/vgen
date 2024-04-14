# vgen

vgen is a code generation tool that creates validation logic using a simple comments.

### Why

I was playing around with a web server and had to validate the incomming json requests. I tried different validation libraries but had one big problem with all of them, **required fields**. When decoding json, fields take the zero value for absent fields, so its impossible to determine between absent and zero valued fields. One solution is to decode to pointers which will become nil in the case of an absent field. Creating these new types with every field as a pointer can be tedious and is what this tool is aiming to solve.

### Getting started

```go
// vgen:[i]
type Email struct {
	Title  string `json:"title"` // vgen:[ req, not_empty, len_lt=10 ]
	Text   string `json:"text"` // vgen:[ req, not_empty, len_lt=200 ]
	Sender string `json:"sender"` // vgen:[ req, not_empty, len_lt=20 ]
}
```

Running the tool on the following type would generate a new type and method

```go
type EmailVgen struct {
    Title *string
    Text *string
    Sender *string
}
func (t EmailVgen) Validate() (Email, error) { ... }
func EmailFromJson(bytes []byte) (Email, error) { ... }
```

This would be the response for an example json request

```json
{
    "title": "this is a hello message that is too long",
    "text": "",
}
=>
{
  "title": [
    "len must be < 10"
  ],
  "text": [
    "can not be empty"
  ],
  "sender": [
    "required"
  ]
}
```

### Lists

Lists are also supported. Rules will be applied to all elements. 
```go
// vgen:[i]
type TicTacToe struct {
    Board [][]int `json:"board"` // vgen:[gte=0, lte=1 ]
}
```

### Nested types
Nesting custom types also works with vgen  
By adding the ```i``` rule the validation will dive into the nested type
```go
// vgen:[i]
type Person struct {
    Name string `json:"name"` // vgen:[ req, not_empty ]
}

// vgen:[i]
type House struct {
    Address string `json:"address"` // vgen:[ req, not_empty ]
    Owner Person `json:"owner"` // vgen:[ req, i ]
}
```

### Custom rules
If a rule is not built into vgen you can create custom rules

```go
// vgen:[i]
type LoginCredentials struct {
    Name string `json:"name"` // vgen:[ req, not_empty ]
    Password string `json:"password"` // vgen:[ req, custom=passwordCheck]
}

func passwordCheck(password string) error {
    if password != "123" {
        return fmt.Errorf("invalid password")
    }
    return nil
}
```

### More examples

For more advanced examples look in the [examples folder](examples)
