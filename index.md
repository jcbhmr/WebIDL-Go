---
title: Go language binding for Web IDL
---

# Go language binding for Web IDL

- **This version:** https://jcbhmr.me/WebIDL-Go/
- **Issue tracking:** [GitHub](https://github.com/jcbhmr/WebIDL-Rust/issues)
- **Repository:** https://github.com/jcbhmr/WebIDL-Go
- **Editors:** [Jacob Hummer](https://jcbhmr.me)
- **License:** [CC-BY-4.0 License](https://github.com/jcbhmr/WebIDL-Rust/blob/main/LICENSE)

## Status of this document

💡 Inspired by [Java language binding for Web IDL](https://www.w3.org/TR/WebIDL-Java/). 🛑 **Very unofficial.** Not associated with W3C or WHATWG in any way. **This is not a standards document.** This is just a rough collection of ideas on how to map some types from Web IDL to Go. It's sorta like a mini reference book. If you have a better way to do things, chances are you're right! [Open an Issue!](https://github.com/jcbhmr/WebIDL-Rust/issues/new) I'm always looking for more egonomic & better ways to translate Web IDL concepts to Go. ❤️

## Table of contents

1. [Introduction](#introduction)
2. [Go binding](#go-binding)
    1. [Names](#names)
    2. [Go type mapping](#go-type-mapping)
        1. [Restricted `float` and `double`](#restricted-float-and-double)

## Introduction

**Use case 1:** Ergonomically exposing JavaScript APIs to Go. There are currently a few Go wrappers that attempt to provide a wrapper around JavaScript APIs like [the DOM API](https://dom.spec.whatwg.org/) or [the Fetch API](https://fetch.spec.whatwg.org/) so that you can write Go code that calls out to these JavaScript platform features without using `syscall/js` to `.Get("document")` `.Get("querySelector")` over and over. These interfaces all differ slightly in how they map certain Web IDL (and JavaScript-land) concepts into Go-land. This document aims to provide some conventions for how bindings between Web IDL (which is usually through a JavaScript & WASM layer) and Go-land should be represented on the Go side of things.

**Use case 2:** Outlining patterns and conventions to write Go-native implementations of web and web-adjacent Web IDL-based specifications. For example [the Web Share API](https://w3c.github.io/web-share/) could be brought to Go desktop & CLI apps with a `share()` function or similar. This is particularily poignant for writing implementations like [the Fetch API](https://fetch.spec.whatwg.org/) that use the `syscall/js` functions from the browser runtime when `GOOS=js` but also offer a native Go implementation for use outside JS+WASM environments.

> The original use case for this was to play around and implement the JavaScript `fetch()` function in Go. Problems arose when there were diverging ways to implement `Promise<Response>`. Should it be a `Promise[T]`? A channel? How would it integrate with other concurrency stuff? What if this theoretical `github.com/octocat/fetch` package could use the browser's native `fetch()` API if it were running in `GOOS=js`? Then it spiraled from there into questions like "well how do you represent DOM node inheritance?" and "what if I want to downcast a `Node` from `QuerySelector()` into an `HTMLAnchorElement`?". All of that meant I had to think and pick an arbitrary way to represent Web IDL types in Go-land.

&mdash; [@jcbhmr](https://jcbhmr.me/)

This document tries to outline some good ideas for translating Web IDL types used in Web IDL specifications like [the Fetch API](https://fetch.spec.whatwg.org/) and [the DOM API](https://dom.spec.whatwg.org/) into native and (at least to some degree) ergonomic Go types. For example, how would you map a `Promise<ArrayBuffer>` for `Response#arrayBuffer()` into Go? What about `Node` or `HTMLElement`?

**ℹ Example:** Here's some examples of some popular Web IDL interfaces and how they would map to Go-land types using the conventions outlined in this document.

<table><td>

```ts
interface Response {
    // ...
};
Response includes Body;

interface mixin Body {
    Promise<ArrayBuffer> arrayBuffer();
    Promise<any> json();
    Promise<USVString> text();
    // ...
};
```

<td>

```go
type Response interface {
  Body
  // ...
}

type Body interface {
  ArrayBuffer() promise[[]byte]
  JSON() promise[any]
  Text() promise[string]
  // ...
}
```

<tr><td>

```ts
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
  NodeElementNode uint16 = 1
  NodeAttributeNode uint16 = 2
  // ...
)

type Node interface {
  EventTarget
  NodeType() uint16
  NodeValue() *string
  SetNodeValue(v *string)
  CloneNode() (Node, error)
  CloneNodeDeep(deep bool) (Node, error)
}
```

<tr><td>

```ts
interface EventTarget {
  undefined addEventListener(
    DOMString type,
    EventListener? callback,
    optional (AddEventListenerOptions or boolean) options = {}
  );
  // ...
};

callback interface EventListener {
  undefined handleEvent(Event event);
};
```

<td>

```go
type EventTarget interface {
  AddEventListener(
    type_ string,
    callback EventListenerFunc
  )
  AddEventListenerEventListener(
    type_ string,
    callback EventListener
  )
  AddEventListenerOptions(
    type_ string,
    callback EventListener,
    options bool
  )
  AddEventListenerOptionsAddEventListenerOptions(
    type_ string,
    callback EventListener,
    options AddEventListenerOptions
  )
  AddEventListenerEventListenerOptions(
    type_ string,
    callbackObject EventListenerObject,
    options bool
  )
  AddEventListenerEventListenerOptionsAddEventListenerOptions(
    type_ string,
    callbackObject EventListenerObject,
    options AddEventListenerOptions
  )
  // ...
}

type EventListenerFunc = func(event Event)
type EventListener interface {
  HandleEvent(event Event)
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

Here's some of the quick ones that are pretty intuitive and don't require much explaination.

| Web IDL type | Go type |
| --- | --- |
| `any` | `any` |
| `boolean` | `bool` |
| `byte` | `int8` |
| `octet` | `byte` |
| `short` | `int16` |
| `unsigned short` | `uint16` |
| `long` | `int32` |
| `unsigned long` | `uint32` |
| `unrestricted float` | `float32` |
| `unrestricted double` | `float64` |
| `bigint` | `math/big.Int` |
| `DOMString` | `string` |
| `ByteString` | `string` |
| `symbol` | _Not specified_ |

#### Restricted `float` and `double`

Remember that `float64` includes the `NaN` and `Inf` values in its type domain. The Web IDL `unrestricted double` is directly equivalent to Go's `float64`. For `double` (the restricted variant) implementations **should `panic()`** if a non-number value (like `NaN` or `-Inf`) was passed. This applies to `float` and `double` Web IDL types.

**ℹ Example:**

```go
func Add(a float64, b float64) float64 {
  if math.IsNaN(a) || math.IsInf(a, 0) {
    panic(errors.New("a is not in restricted double range"))
  }
  if math.IsNaN(b) || math.IsInf(b, 0) {
    panic(errors.New("b is not in restricted double range"))
  }
  return a + b
}
```

#### `USVString`

> The USVString type corresponds to scalar value strings. Depending on the context, these can be treated as sequences of either 16-bit unsigned integer code units or scalar values.

Honestly I'm not entirely sure that this is the right way to do things but this is how I've been doing the Go `string` to Web IDL `USVString` conversion:

```go
func DoThing(v string) {
  if !utf8.Valid(v) {
    panic(errors.New("v is not a well-formed UTF-8 string"))
  }
}
```

If you have a better explaination or some other insight into what `DOMString` vs `USVString` means in the context of a Go string I'm all ears. 😊

#### `undefined`

**In function return position** `undefined` should be omitted from the Go function signature.

**ℹ Example:**

<table><td>

```ts
interface Storage {
  setter undefined setItem(DOMString name, DOMString value);
  // ...
};
```

<td>

```go
interface Storage {
  SetItem(name string, value string)
}
```

</table>

**When used in `Promise<undefined>`** or other places where a type is required use `struct{}`.

**When part of a type union** like `(DOMString or undefined)` use a `nil`-able variant of the other union members.

**ℹ Example:**

<table><td>

```ts
interface CustomElementRegistry {
  (CustomElementConstructor or undefined) get(DOMString name);
  // ...
};

callback CustomElementConstructor = HTMLElement ();
```

<td>

```go
interface CustomElementRegistry {
  // Since CustomElementConstructor is nil-able
  Get(name string) CustomElementConstructor
}

type CustomElementConstructor = func() HTMLElement
```

</table>

#### `object` type

`object` is similar to `any` but it **cannot be `null` or a primitive**. In Go-land that means it's like `any` just with a `panic()` check to make sure it's not `nil` or a `string|int|bool|etc...` type.

**ℹ Example:**

```go
func DoThing(v any) {
  if v == nil || reflect.ValueOf(v).Kind() < reflect.Array {
    panic(errors.New("v is not object type"))
  }
}
```

#### Interface types

Interfaces in Web IDL are the basic building blocks of the APIs that they define. They are equivalent to `class` types in JavaScript and Java. They act like `*T` class pointers in C++. In Go you might assume that we could represent these `interface` declarations in a similar way with `struct`. That has some a critical lacking feature though: downcasting.

You can't downcast a `struct` `*Node` into a `struct` `*HTMLAudioElement` in Go. You can only cast with `interface`-es. This means that if `.QuerySelector()` returns a `*Element` then there's no way to try to `instanceof HTMLAnchorElement` or `instanceof HTMLVideoElement` and then cast it to those types to use `.click()` or `.play()`.

**ℹ Example:** How would you do this in Go when using `*Node` as the return type?

```java
// This is very ergonomic in Java.
var aOrVideo = document.querySelector("#video")
if (aOrVideo instanceof HTMLVideoElement video) {
  video.play()
} else if (aOrVideo instanceof HTMLAnchorElement a) {
  a.click()
}
```

The solution is to use `interface` types instead. Go interfaces let us freely cast between things at runtime. A helpful trick is to **use _sealed_ interfaces** to brand your own types similar to how JavaScript Web IDL works.

**ℹ Example:**

```go
type Node interface {
  NodeValue() *string
  SetNodeValue(v *string)
  // ...
  sealedNode()
}
type nodeData struct {}
func (n *nodeData) sealedNode() {}
var _ Node = (*nodeData)(nil)

type Element interface {
  Node
  AppendChild(...) (Node, error)
  Remove() error
  InnerHTML() string
  // ...
  sealedElement()
}
type elementData struct {}
func (e *elementData) sealedNode() {}
var _ Element = (*elementData)(nil)
```

**Interface constructors** are pulled out into a `New___()` function where you fill in the blank with the interface name. For example the `Request` interface has a constructor `NewRequest()`. For interfaces without a constructor you should omit this `New___()` function. Use the same rules as described in function overloading to create multiple variants of the constructor if there are multiple signatures.

**ℹ Example:**

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
type Headers interface {
  Append(name string, value string)
  Delete(name string)
  // ...
}
func NewHeaders() Headers {}
func NewHeadersInit(init map[string]string) Headers {}
func NewHeadersInitStringsList(init [][]string) Headers {}
```

</table>

##### `interface mixin`

Mixins are used by Web IDL to define a set of functions implemented by multiple interfaces. The best example of this is the `Body` mixin which defines `.arrayBuffer()`, `.json()`, and `.text()` on `Request` and `Response` objects.

The corollary to this in Go is just another interface!

**ℹ Example:**

<table><td>

```webidl
interface mixin Body {
  readonly attribute boolean bodyUsed;
  [NewObject] Promise<ArrayBuffer> arrayBuffer();
  // ...
};

Request includes Body;
Response includes Body;
```

<td>

```go
type Body interface {
  BodyUsed() bool
  ArrayBuffer() promise[[]byte]
	// ...
}

type Request interface {
  Body
}
type Response interface {
  Body
}
```
		       
</table>

#### Callback interface types

This concept is a way to codify the below into Web IDL:

```js
// Both of these work!
globalThis.addEventListener("load", { handleEvent: () => console.log("handleEvent") })
globalThis.addEventListener("load", () => console.log("function"))
```

It's rarely used. Trait it as though it were a union of `type ___Func = func(...)...` and `type ___ interface {...}` types.

<table><td>

```ts
interface EventTarget {
  undefined addEventListener(
    DOMString type,
    EventListener? callback,
    optional (AddEventListenerOptions or boolean) options = {}
  );
  // ...
};

callback interface EventListener {
  undefined handleEvent(Event event);
};
```

<td>

```go
type EventTarget interface {
  AddEventListener(
    type_ string,
    callback EventListenerFunc
  )
  AddEventListenerEventListener(
    type_ string,
    callback EventListener
  )
  AddEventListenerOptions(
    type_ string,
    callback EventListener,
    options bool
  )
  AddEventListenerOptionsAddEventListenerOptions(
    type_ string,
    callback EventListener,
    options AddEventListenerOptions
  )
  AddEventListenerEventListenerOptions(
    type_ string,
    callbackObject EventListenerObject,
    options bool
  )
  AddEventListenerEventListenerOptionsAddEventListenerOptions(
    type_ string,
    callbackObject EventListenerObject,
    options AddEventListenerOptions
  )
  // ...
}

type EventListenerFunc = func(event Event)
type EventListener interface {
  HandleEvent(event Event)
}
```

</table>

#### Dictionary types

Web IDL interfaces map quite well to Go `struct`. The biggest obstacle here is optional `struct` fields. Since Go doesn't have a Rust `Option<T>` or a Java `null`-able primitive (`Integer`, `Double`, `String`) objects it means we need to use something else. The most Go-like way of doing things is with `nil`-able types. For interfaces this is no sweat! For structs it just means passing by `&theStruct`. For primitives it means passing by `*string`.

**💡 Tip:** Go doesn't allow you to take the `&` address of literal values. You can use a helper function like this for that:

```go
func ptr[T any](v T) *T {
  return &v
}

func DoThing(v *string) {}

func main() {
  DoThing(ptr("hello ptr magic!"))
}
```

**ℹ Example:**
	
<table><td>

```webidl
dictionary ResponseInit {
  unsigned short status = 200;
  ByteString statusText = "";
  HeadersInit headers;
};
```

<td>

```go
type ResponseInit struct {
  Status *uint16
  StatusText *string
  Headers map[string]string
  HeadersStringsList [][]string
}
```
	
</table>

#### Enumeration types

Similar JavaScript and Java there is just a type. There's no accompanying list of valid enum values exposed as a programmatic object or static list. The valid values are provided in documentation and you should provide them as string literals. If the provided string isn't a valid enum value you should `panic()` instead of returning an `error`.

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

```go
type future[T] interface {
  Wait() (T, error)
}
```