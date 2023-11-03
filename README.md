# vgen

vgen is a code generation tool that creates validation logic using a simple comments.

### Why
I was playing around with a web server and had to validate the incomming json requests. I tried different validation libraries but had one big problem with all of them, **required fields**. When decoding json, fields take the zero value for absent fields, so its impossible to determine between absent and zero valued fields. One solution is to decode to pointers which will become nil in the case of an absent field. Creating these new types with every field as a pointer can be tedious and is what this tool is aiming to solve. 

### Example
```go
// vgen:[i]
type Email struct {
	Title  string // vgen:[ req, not_empty, len_lt=10 ]
	Text   string // vgen:[ req, not_empty, len_lt=200 ]
	Sender string // vgen:[ req, not_empty, len_lt=20 ]
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
```
This would be the response for an example json request
```
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

### More examples
For more advanced examples look in the [examples folder](examples)
