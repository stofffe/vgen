# vgen

A type safe validation library using code generation

## Motivation

When decoding json data it is impossible to differentiate between zero and nil values. 
One solution is to use pointers which become nil in the case of abscent fields. But using pointers 
can be tedious and creating new types leads to a lot of boilerplate.

vgen uses code generation to generate this boilerplate for you. Code generation also allows
for compile time checking of validation rules

## Installation

### Go

Run the following command to download the binary using ```go install```

```zsh
go install github.com/stofffe/vgen/cmd/vgen@latest
```

## Usage

Annotate a struct with the [tag](#tags) ```vgen(include)```

```go
// vgen(include)
type Person struct {
    name string
    age int
}
```

Run the ```vgen generate <filename>``` to generate the following types

```go
type PersonVgen struct { ... }
type PersonRules struct { ... }

func (v PersonVgen) Validate(rules PersonRules) ErrorMap { ... }
func (v PersonVgen) Convert() Person { ... }
func (v PersonVgen) ValidatedConvert() (Person, ErrorMap) { ... }
```

Use type safe validation

```go
func main() {
    // Create person manually or decode from json
    person := PersonVgen { ... }

	rules := PersonRules{
		name: vgen.RulesRequired(
			vgen.Eq("bob"),
		),
		age: vgen.RulesOptional(
			vgen.Gte(18),
		),
	}

	result, err := person.ValidatedConvert(rules)
	if err != nil {
		log.Fatal(err.Debug())
	}
	fmt.Println(result)
}
```

Output

```zsh
{
    name: [
        "required",
    ],
    age: [
        "must be greater than or equal to 18",
    ]
}
```

More examples can be found in [examples](examples/)

## Tags

Additional functionality can be added to structs and fields using tags 

```go
// vgen(i)
type Person struct {
    NAME string // vgen(alias=name)
}
```

#### Type Tags
| Tag | Short |Description  |
|-|-|-|
| include | i | Include type in code generation  |

#### Field Tags
| Tag | Short |Description  |
|-|-|-|
| nested | n | Allow for nested vgen types |
| alias | a | Give field different name in errors |

## Rules

#### Rule

vgen validation is based of the ```Rule``` interface

```go
type Rule[T any] interface {
	Validate(key string, input T) ErrorMap
}
```

#### RuleFunc

Custom rules can be created by implementing this ```Rule```, or by creating a ```RuleFunc```
which automatically implements the interface. See [custom example](examples/custom/custom.go) for more information

```go
type RuleFunc[T any] func(key string, input T) ErrorMap
```

#### Implemented rules

vgen includes a lot of default rules, here is the currently implemented rules

| Rule | Description |
|-|-|
| Nested | Validate nested vgen types |
| List | Validate each element of list |
| MapValue | Validate each value element of map |
| MapHasKey | Validate map has key |
| Eq | Validate equality |
| Gt | Validate greater than |
| Gte | Validate greater than or equal to |
| Lt | Validate less than |
| Lte | Validate less than or equal to |
| LenEq | Validate len equality |
| LenGt | Validate len greater than |
| LenGte | Validate len greater than or equal to |
| LenLt | Validate len less than |
| LenLte | Validate len less than or equal to |
| OneOf | Validate in list of elements |
| NotOneOf | Validate not in list of elements |
| Deref | Validate inner value of pointer |
| PrefixMessage | Prefix error message |
| CustomMessage | Replace error with custom message |
