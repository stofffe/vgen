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

...
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

Types and fields can be annotated with tags, the tags inform the ```vgen``` compiler what code to generate

#### Example
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
| Eq | Equality |
| Gt | Greater than |
| Gte | Greater than or equal to |
| Lt | Less than |
| Lte | Less than or equal to |
| OneOf | In list of elements |
| NotOneOf | Not in list of elements |
| MapHasKey | Map has key |
| Nested | Validate nested vgen type |
| List | Validate each element of list |
| MapValue | Validate each value element of map |
| Deref | Validate Inner value of pointer |
| PrefixMessage | Prefix error message |
| CustomMessage | Replace error with custom message |
| Regex | Regex match |
| Ascii | All characters are ascii |
| UUIDv1 | UUID version 1 |
| UUIDv3 | UUID version 3 |
| UUIDv4 | UUID version 4 |
| UUIDv5 | UUID version 5 |
| ListLenEq | List len equality |
| ListLenGt | List len greater than |
| ListLenGte | List len greater than or equal to |
| ListLenLt | List len less than |
| ListLenLte | List len less than or equal to |
| MapLenEq | Map len equality |
| MapLenGt | Map len greater than |
| MapLenGte | Map len greater than or equal to |
| MapLenLt | Map len less than |
| MapLenLte | Map len less than or equal to |
| StringLenEq | String len equality |
| StringLenGt | String len greater than |
| StringLenGte | String len greater than or equal to |
| StringLenLt | String len less than |
| StringLenLte | String len less than or equal to |
