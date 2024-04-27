---
title: Go language binding for Web IDL
---

# Go language binding for Web IDL

💡 Inspired by [Java language binding for Web IDL](https://www.w3.org/TR/WebIDL-Java/) \
🛑 **Very unofficial.** Not associated with W3C or WHATWG in any way.

Copyright 2024 Jacob Hummer \
Licensed under [the CC-BY-4.0 license](LICENSE).

## Status of this document

**🛑 This is not a standards document.** This is just one guy's ideas of how to map some types. It's sorta like a mini reference book at this point for future me.

## Introduction

This document tries to outline some good ideas for translating Web IDL types used in Web IDL specifications like [the Fetch API](https://fetch.spec.whatwg.org/) and [the DOM API](https://dom.spec.whatwg.org/) into native and (at least to some degree) ergonomic Go types. For example, how would you map a `Promise<ArrayBuffer>` for `Response#arrayBuffer()` into Go? What about `Node`?

<table><td>

```webidl
interface Response {
    // ...
}
Response includes Body

interface mixin Body {
    Promise<ArrayBuffer> arrayBuffer();
    // ...
}
```

<td>

```go
type Response struct {
    // ...
}

var _ Body = (*Response)(nil)

interface Body {
    ArrayBuffer() <-chan Unpack[[]byte, error]
}
```

<tr><td>

```webidl
interface Node : EventTarget {
  const unsigned short ELEMENT_NODE = 1;
  const unsigned short ATTRIBUTE_NODE = 2;
  readonly attribute unsigned short nodeType;
  attribute DOMString? nodeValue;
  Node cloneNode(optional boolean deep = false);
  // ...
}
```

<td>

```go
const (
    NodeElementNode = 1
    NodeAttributeNode = 2
    // ...
)

type Node struct {
    EventTarget
    NodeType() uint16
    NodeValue() *string
    SetNodeValue(v *string)
    CloneNode(deep *bool) (*Node, error)
}
```

</table>

## Go binding

### Names

| Web IDL name | Go name |
| --- | --- |
| `innerHTML` | `InnerHTML` |
| `getElementById` | `GetElementByID` |
| `htmlFor` | `HTMLFor` |
| `url` | `URL` |

### Go type mapping

| Web IDL type | Go type |
| --- | --- |
| `any` | `any` |
| `void` | Omit or `struct{}` |
| `boolean` | `bool` |
| `byte` | `int8` |
| `octet` | `byte` |
| `short` | `int16` |
| `unsigned short` | `uint16` |
| `long` | `int32` |
| `unsigned long` | `uint32` |
| `float` | `float32` with `panic()` |
| `unrestricted float` | `float32` |
| `double` | `float64` with documentation |
| `unrestricted double` | `float64` |
| `object` | `any` with `panic()` |

In places where `undefined` should be treated like `void` (like `Promise<undefined>`) use `struct{}`. Otherwise if it makes sense use `nil`; usually coupled with a `nil`-able type like in a `(DOMString or undefined)` either-or return type.

#### Interface types

<table><td>

```webidl
[Exposed=(Window,Worker)]
interface Headers {
  constructor(optional HeadersInit init);
  undefined append(ByteString name, ByteString value);
  undefined delete(ByteString name);
  // ...
};
```

<td>

```go
type Headers struct{}
func NewHeaders(...) *Headers
func (h *Headers) Append(name string, value string)
func (h *Headers) Delete(name string)
// ...
```

</table>

You could alternatively model these types as interfaces without concrete types like this:

```go
type Headers interface {
    Append(name string, value string)
    Delete(name string)
    // ...
}
func NewHeaders(...) Headers
```

But this comes with the minor downside that it just _feels_ a bit off to be using a non-concrete `interface` type for a concrete type that you construct and pass around.

#### `interface mixin`

<table><td>

```webidl
interface mixin Body {
  readonly attribute boolean bodyUsed;
  [NewObject] Promise<ArrayBuffer> arrayBuffer();
  // ...
};
```

<td>

```go
type Body interface {
  BodyUsed() bool
  ArrayBuffer() <-chan result[[]byte]
	// ...
}
```
		       
</table>

You can use the `var _ SomeInterface = (*MyStruct)(nil)` hack to do a compile-time assertion that `*MyStruct` is assignable to `SomeInterface` thus validating that `MyStruct` correctly implements `SomeInterface`.

#### Dictionary types
	
<table><td>

```webidl
dictionary ResponseInit {
  unsigned short status = 200;
	// ...
};
```

<td>

```go
type ResponseInit struct {
  Status *uint16
	// ...
}
```
	
</table>

All properties of the dictionary that are optional (i.e. not marked with `required`) should be `nil`-able. For primitive types and struct types this means using pointers.

Dictionary objects actually would make more sense if they could be expressed in terms of _any object that has `.Status` as a `*uint16` but Go doesn't have that. The next best thing along that route is `interface` stuff. But that comes with the problem of needing to define functions for these properties which quickly gets in the way of properties named the same thing in struct literals. The TLDR is that this doesn't compile:

```go
type TheOptionsI interface {
    Required() *bool
}

type TheOptions struct {
    Required *bool
}

func (t *TheOptions) Required() *bool {
    return t.Required
}

func DoThing(options TheOptionsI) {
    // ...
}

func main() {
    DoThing(&TheOptions{
        Required: lo.ToPtr(true)
    })
}
```

```
./prog.go:19:22: field and method with the same name Required
	./prog.go:16:2: other declaration of Required
./prog.go:20:9: cannot use t.Required (value of type func() *bool) as *bool value in return statement
```

So supporting arbitrary objects that fit the field name requirements of the interface (as functions) like `.Required()` doesn't work **because it interferes with providing those fields as struct literals**.

You also can't get around this with generics. Generic parameter interfaces must be `interface{T|U|V}` unions and normal method interfaces don't work.

```go
type TheOptionsI interface {
	Required() *bool
}

type TheOptions struct {
	Required *bool
}

func DoThing[T *TheOptions | TheOptionsI](options T) {
	// ...
}

func main() {
	fmt.Println("Hello, 世界")
	DoThing(&TheOptions{
		Required: lo.ToPtr(true),
	})
}
```

```
./prog.go:19:30: cannot use main.TheOptionsI in union (main.TheOptionsI contains methods)
```

#### Enumeration types

Similar JavaScript and Java, there is just a type. There's no accompanying list of valid enum values exposed as a programmatic object or static list. The valid values are provided in documentation and you should provide them as string literals. If the provided string isn't a valid enum value you should `panic()` instead of returning an `error`.

<table><td>

```webidl
enum RequestPriority { "high", "low", "auto" };
```

<td>

```go
// enum RequestPriority { "high", "low", "auto" };
type RequestPriority = string
```

</table>

```go
func DoThing(priority RequestPriority) {
    switch priority {
    case "high": log.Print("WOW this is important")
    case "low": log.Print("zzzzzz")
    case "auto": log.Print("when i get the chance")
    // make sure you always handle the invalid case!
    default: panic(fmt.Errorf(`%s is not in "high"|"low"|"auto"`, priority))
    }
}
```

#### `Promise<T>`

**`Promise<undefined>` with no `error`** can be a `<-chan struct{}` that emits a single `struct{}{}` value.

```go
func WaitABit() <-chan struct{} {
  c := make(chan struct{})
  go func(){
    defer close(c)
    time.Sleep(2000 * time.Millisecond)
    c <- struct{}{}
  }()
  return c
}
```

**`Promise<T>` with no `error`** can be modeled as a `<-chan T`. The `chan` should provide **a single value**.

```go
func Gimme() <-chan int {
  c := make(chan int)
  go func(){
    defer close(c)
    time.Sleep(200 * time.Millisecond)
    c <- 42
  }()
  return c
}
```

**`Promise<undefined>` with `error`** uses `<-chan error`; same as `Promise<T>` where `T` would be the Go `error` type. 

```go
func FailSoon() <-chan error {
  c := make(chan error)
  go func(){
    defer close(c)
    time.Sleep(200 * time.Millisecond)
    c <- errors.New("timed out")
  }()
  return c
}
```

**`Promise<T>` with `error`** should return a `<-chan R` where `R` is some kind of struct, interface, etc. that has a `.Unpack()` method which returns the needed `(T, error)` multivalue.

```go
// Don't actually need to export this.
type packed[A any, B any] struct {
	A A
    B B
}

func (r packed[A, B]) Unpack() (A, B) {
	return r.A, r.B
}

func AskUser() <-chan packed[string, error] {
	c := make(chan packed[string, error])
	go func() {
		defer close(c)
		time.Sleep(200 * time.Millisecond)
		if rand.Intn(2) == 1 {
			c <- packed[string, error]{"i like go", nil}
		} else {
			c <- packed[string, error]{"", errors.New("bad luck!")}
			return
		}
	}()
	return c
}

func main() {
	res, err := (<-AskUser()).Unpack()
	fmt.Println(res, err)
}
```

Note that we don't really care what the `R` result-like type is as long as it has a `.Unpack()` method on it which returns the multivalues.
