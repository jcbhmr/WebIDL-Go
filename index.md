---
title: Go language binding for Web IDL
---

- **This version:** https://jcbhmr.me/WebIDL-Go/
- **Issue tracking:** [GitHub](https://github.com/jcbhmr/WebIDL-Go/issues)
- **Editors:** [Jacob Hummer](https://jcbhmr.me)

[CC0-1.0 License](https://github.com/jcbhmr/WebIDL-Go/blob/main/LICENSE)

<small>Not affiliated with Go, Google, W3C, or WHATWG.</small>

---

## Abstract

Outline conventions & ideas for mapping Web IDL concepts & types into native Go patterns, conventions, and types.

## Introduction

_This section is non-normative._

💡 This document was inspired by [Java language binding for Web IDL](https://www.w3.org/TR/WebIDL-Java/) and [Web IDL JavaScript Binding](https://webidl.spec.whatwg.org/#javascript-binding). This is just a rough collection of ideas on how to map types from Web IDL to Go. It's like a miniature reference book. If you have a better way to do things, chances are you're right! [Open an Issue!](https://github.com/jcbhmr/WebIDL-Go/issues/new) I'm open to more egonomic & better ways to translate Web IDL concepts to Go. ❤️

### Use case: ergonomic Go ↔ Web API bindings

There currently is a lack of high-quality Go bindings to Web APIs for use when writing Go code that targets `js/wasm`.

[dominikh/go-js-dom](https://github.com/dominikh/go-js-dom) provides a lot of bindings to access [DOM](https://dom.spec.whatwg.org/) APIs.

- It doesn't let you cast very easily. It's inconsistent with some types being `struct` and others being `interface`. The DX and ergonomics could be improved by using `interface` everywhere.
- It doesn't provide modern DOM features like `ParentNode#append()`. It's a bit outdated. This could be improved with a Web IDL-to-Go scaffolding generator and more formalized Go ↔ Web IDL rules.
- It doesn't handle optional primitives at all. `Node#nodeValue` is a `DOMString?` in Web IDL but is defined as `NodeValue() string` in go-js-dom. Using `nil`-able types would be more idiomatic Go.
- It doesn't work well with complex function signatures like `EventTarget#addEventListener()` which has a callback interface and a union type. A standardized well-known way to idiomatically convert Web IDL overloads, union types, and optional parameters to Go would be helpful.

[realPy/hogosuru](https://github.com/realPy/hogosuru) is a Go web framework that offers some Web API bindings too.

- It doesn't embrace union types. `EventTarget#addEventListener()` is represented only as `AddEventListener(string, func(Event))` with no options parameter. An ergonomic way to idiomatically represent overloads and optional parameters in Go would be nice.
- Packages are fragmented. There are so many packages! One for each DOM definition is too many. A more cohesive API might group things based on which specification ([DOM](https://dom.spec.whatwg.org/), [Fetch](https://fetch.spec.whatwg.org/), [Web Share API](https://w3c.github.io/web-share/), etc.) that they come from.
- Everything returns a `(T, error)`. This is just a gripe of how granular I think the error handling is.
- `Node#nodeValue` is represented as `NodeValue() string`. The `.nodeValue` of a node is actually `DOMString?` though. A standard way to represent nullable [Web IDL] values in Go (particularily primitives) would help.

### Use case: standard-ish way to port Web APIs to Go

[URL Pattern](https://urlpattern.spec.whatwg.org/), [DOM](https://dom.spec.whatwg.org/), [Fetch](https://fetch.spec.whatwg.org/), [Web Share API](https://w3c.github.io/web-share/), [Clipboard API and events](https://www.w3.org/TR/clipboard-apis/), and more are all cool Web APIs that could be ported to Go. Not because being Web IDL-based is better than a custom API, but because it gives developers a very concrete API surface to cover. This is particularily useful for web developers who may already be familiar with a particular API from the browser. As a side effect it makes it very easy to bind to the existing JavaScript bindings of the same Web API when compiling for `js/wasm`.

🌎 JavaScript APIs are also very familiar to most web developers and it's very cool 😎 when you can use the same concepts and libraries uniformly across domains and languages. It's just cool.

### The gist

- All attributes are getter & setter functions. Getters have the same name as the property, setters have a `Set` prefix.
- All names are converted to Go style.
- Namespaces and static methods are free functions.
- `Window` and other global scope interfaces stuff becomes top-level.
- Always use Go modules.
- Mixins use wrapper types like `type Mixin Actual`.
- Promises are `func (T) Wait() (U, error)`. Each package defines its own private promise type. It must have a `.Wait()` method that returns `(T, error)`.
- Extend through composition. Use a `this any` field to maintain a reference to narrower types.
- Use `nil`-able types for `T?` and other optionals when it's a single field. Use `(T, bool)` for return types.
- `undefined` is omitted where possible, `nil` where otherwise possible, and `struct{}{}` where necessary.
- `Uint8Array` and friends are all mapped to `[]T` of the appropriate type. `ArrayBuffer` is `[]byte`.
- Union types are `any` with body logic to handle each acceptable type, otherwise `panic()`.
- Additional overloads of a function are denoted with a number suffix.

## Go binding

1. You should use Go 1.11 modules when developing Go Web IDL bindings. That means `go mod init` and not `GOPATH`.
2. Try to group things based on which specification they are a part of. This helps segement the vast collection of browser APIs into develop-able chunks.
3. Put everything on one level. Don't nest things in sub-packages; put it all on the root. Some developers will want to `import . "github.com/octocat/go-fetch"` to make `Fetch()` and all associated things top-level.
4. The `Window`, `WorkerGlobalScope`, etc. interfaces are all ✨special since they are global interfaces. Anything that would be defined as a method on `Window` or other global scope interface should be defined **directly at the top level of the Go package** instead of as methods on a `Window` `struct`/`interface`. You should also **omit** the `Window` and other global interfaces from the types that you define. Think of `Window` as an invisible package-level interface instead of a concrete type.
6. Many Web IDL definitions (particularily in [HTML](https://html.spec.whatwg.org/multipage/)) will make reference to "the global document" or "the current window". There is no "current window" in a typical Go program. You should try to follow the spirit if not the letter of these definitions.

**Example:**

```webidl
partial interface Window {
  undefined alert(DOMString message);
};
```

```go
func Alert(message string) {
  // ...
}
```

### Names

[Remember: to be exported from a package a Go name must start with an uppercase letter.](https://go.dev/tour/basics/3)

Use your own best judgement to split the words and PascalCase-ify them to make them conform to the Go convention. Note that initialisms should be all-caps.

For const & static members of interfaces, use the interface name followed by the member name.

For Web IDL constructors: use `New` followed by the interface name.

Getters are to be converted into accessor functions with the same name. Setters should be converted to single-parameter setter functions with the prefix `Set` added.

**Examples:**

| Web IDL | Go |
| --- | --- |
| `Element#innerHTML` | `Element#InnerHTML()` & `Element#SetInnerHTML()` |
| `Document#getElementById()` | `Document#GetElementByID()` |
| `HTMLLabelElement#htmlFor` | `HTMLLabelElement#HTMLFor()` & `HTMLLabelElement#SetHTMLFor()` |
| `Response#url` | `Response#URL()` |
| `XMLHttpRequest` | `XMLHTTPRequest` |
| `HTMLHtmlElement` | `HTMLHTMLElement` |
| `HTMLPreElement` | `HTMLPreElement` |
| `Node.TEXT_NODE` | `NodeTextNode` |
| `Document#doctype` | `Document#Doctype()` |
| `Element#id` | `Element#ID()` & `Element#SetID()` |
| `new Document()` | `NewDocument()` |

### Go type mapping

💡 Can't find a special Web IDL type defined in this section? [Add it yourself](https://github.com/jcbhmr/WebIDL-Go) or [open an issue](https://github.com/jcbhmr/WebIDL-Go/issues/new). ❤️

The following list of types doesn't require much explaining.

| Web IDL | Go |
| --- | --- |
| `any` | `any` |
| `boolean` | `bool` |
| `byte` | `int8` |
| `octet` | `byte` |
| `short` | `int16` |
| `unsigned short` | `uint16` |
| `long` | `int32` |
| `unsigned long` | `uint32` |
| `long long` | `int64` |
| `unsigned long long` | `uint64` |
| `unrestricted float` | `float32` |
| `unrestricted double` | `float64` |
| `bigint` | [`*math/big.Int`](https://pkg.go.dev/math/big#Int) |
| `DOMString` | `string` |
| `ByteString` | `string` |

#### Restricted float and double

Reminder: the `unrestricted double` type means that the value *can* be `NaN`, `+Infinity`, or `-Infinity`. The counterpart is the plain old `double` type which *cannot* be `NaN`, `+Infinity`, or `-Infinity`.

Web IDL `float` values in Go are represented as Go `float32` values. Web IDL `double` values in Go are represented as Go `float64` values.

Implementers should `panic()` if their input Go `float32` or `float64` value is `NaN`, `+Infinity`, or `-Infinity` when the Web IDL signature uses a restricted `float` or `double` type.

**Example:**

```webidl
partial interface Window {
  double add(double a, double b);
  unrestricted double addUnrestricted(unrestricted double a, unrestricted double b);
};
```

```go
func Add(a, b float64) float64 {
  if math.IsNaN(a) || math.IsInf(a, 0) { panic("a is not a double") }
  if math.IsNaN(b) || math.IsInf(b, 0) { panic("b is not a double") }
  return a + b
}
func AddUnrestricted(a, b float64) float64 {
  return a + b
}
```

#### `undefined`

**In function return position** `undefined` should be omitted from the Go function signature.

**Example:**

```webidl
interface Storage {
  setter undefined setItem(DOMString name, DOMString value);
  // ...
};
```

```go
type Storage interface {
  SetItem(name string, value string)
}
```

**When used in `sequence<undefined>`** or other places where a type is required use `struct{}`.

```webidl
interface A {
  sequence<undefined> doThing();
};
```

```go
type A interface {
  DoThing() []struct{}
}
```

**When part of a type union in return position** treat it like returning `T?`. That means using the `(T value, bool ok)` return pattern.

```webidl
interface CustomElementRegistry {
  (CustomElementConstructor or undefined) get(DOMString name);
  // ...
};
callback CustomElementConstructor = HTMLElement ();
```

```go
type CustomElementRegistry interface {
  Get(name string) (CustomElementConstructor, bool)
}
type CustomElementConstructor = func() HTMLElement
```

**When part of a type union in parameter position** treat it like `T?`. That means accepting a `nil`-able type or making the type `nil`-able via a pointer.

```webidl
interface B {
  undefined doThing(a (long or undefined))
}
```

```go
type B interface {
  DoThing(a *int32)
}
```

**Otherwise** treat `undefined` as though it were a `struct{}`. This also applies in situations like `(long? or undefined)`.

```webidl
interface C {
  // accepts 42, null, or undefined.
  undefined doThing(c (long? or undefined))
}
```

```go
type C interface {
  // differentiates between 42 (int32), nil, and struct{}
  DoThing(c any)
}
```

#### `USVString`

> The USVString type corresponds to scalar value strings. Depending on the context, these can be treated as sequences of either 16-bit unsigned integer code units or scalar values.

Honestly I'm not entirely sure how to represent `USVString` in Go.

If you have a better explaination or some other insight into what `DOMString` vs `USVString` means in the context of a Go string I'm all ears. 😊

**Example:**

```go
func DoThing(v string) {
  if !utf8.Valid(v) { panic("v is not a well-formed UTF-8 string") }
}
```

#### `object`

`object` is similar to `any` but it **cannot be `null` or a primitive**. In Go-land that means it's like `any` just with a `panic()` check to make sure it's not `nil` or a `string|int|bool|etc...` type.

**Example:**

```go
func DoThing(v any) {
  if v == nil || reflect.ValueOf(v).Kind() < reflect.Array {
    panic("v is not object type")
  }
}
```

#### `symbol`

Not yet defined.

#### Interface types

To support downcasting (`Node` to `HTMLButtonElement`) easily Web IDL types are represented as Go `interface` instead of `struct`. Each Web IDL `interface` type should have a corresponding Go `interface` type. Unless it's `Window` or another global interface, in which case it should be omitted and all its members should go in the package scope.

**Example:**

```webidl
interface Node : EventTarget {
  readonly attribute DOMString nodeName;
};
interface Window {
  undefined alert(DOMString message);
};
```

```go
type Node interface {
  EventTarget
  NodeName() string
}
func Alert(message string) {}
```

You can't downcast a `struct` `*Node` into a `struct` `*HTMLAudioElement` in Go. You can only cast with `interface`-s. This means that if `.QuerySelector()` returns a `*Element` then there's no way to try to `instanceof HTMLAnchorElement` or `instanceof HTMLVideoElement` and then cast it to those types to use `.Click()` or `.Play()`.

```go
element, err := document.QuerySelector("audio")
// This won't work
// So how do you convert element (which is *Element) into an *HTMLAudioElement?
audio := element.(*HTMLAudioElement)
audio.Play()
```

```go
element, err := document.QuerySelector("audio")
// This is how you do it using interfaces
audio := element.(HTMLAudioElement)
audio.Play()
```

One trick to achieve self-introspection even though nested structs cannot see their parent container types is to [use the `This` trick](https://jcbhmr.me/2024/12/26/go-abstract-class/).

**Interface constructors** are pulled out into a `New___()` function where you fill in the blank with the interface name. For example the `Request` interface has a constructor `NewRequest()`. For interfaces without a constructor you should omit this `New___()` function. Treat interface constructors as though they were their own `Window` function and apply all the same union, overload, and optional rules.

**Example:**

```webidl
typedef (sequence<sequence<ByteString>> or record<ByteString, ByteString>) HeadersInit;
[Exposed=(Window,Worker)]
interface Headers {
  constructor(optional HeadersInit init);
  // ...
};
```

```go
type HeadersInit = any
func NewHeaders(init HeadersInit) Headers {}
```

##### `interface mixin`

Mixins are used by Web IDL to define a set of functions implemented by multiple interfaces. The best example of this is the `Body` mixin which defines `Body#arrayBuffer()`, `Body#text()`, and `Body#json()` for the `Request` and `Response` types.

In Go, that could corrolate to an extended/embedded interface. The one problem with this is that it requires each type to know which mixins it incorporates instead of leaving the door open to attach more mixins to a type at any time from any package. Instead, use a wrapper type that can be defined outside the package where the primary interface was defined.

```webidl
interface mixin Body {
  readonly attribute boolean bodyUsed;
  [NewObject] Promise<ArrayBuffer> arrayBuffer();
  // ...
};
Request includes Body;
Response includes Body;
```

```go
type Body interface {
  BodyUsed() bool
  ArrayBuffer() func() ([]byte, error)
  // ...
  isBody()
}
type BodyImpl struct {
  X__Inner any
  X__Body ConceptBody
}
// x: Request | Response from otherpackage
func WrapBody(x any) *BodyImpl {
  if v, ok := x.(*otherpackage.RequestImpl); ok {
    return &BodyImpl{v, v.X__Body}
  } else if v, ok := x.(*otherpackage.ResponseImpl); ok {
    return &BodyImpl{v, v.X__Body}
  } else {
    panic("x not Request | Response from otherpackage")
  }
}
// Implement .BodyUsed(), .ArrayBuffer(), etc.
```

🤔 I don't really know what the best conventions are here. The crux is the need to separate the `Navigator` object's core interface definition (defined in HTML) and `NavigatorStorage` which is mixed-in to `Navigator` (defined in [Storage](https://storage.spec.whatwg.org/)).

#### Callback interface types

Callback interfaces are a way to codify the below into Web IDL:

```js
// Both of these work!
globalThis.addEventListener("load", { handleEvent: () => console.log("handleEvent") })
globalThis.addEventListener("load", () => console.log("function"))
```

It's rarely used.

Treat callback interfaces as though they were just a union type of the function signature and an unnamed interface with that single method defined. Practically, this means that callback interfaces are represented as `any`.

**Example:**

```webidl
interface Thing {
  undefined addEventListener(listener EventListener);
};
callback interface EventListener {
  undefined handleEvent(Event event);
};
```

```go
type EventListener = any
func (e *thingImpl) AddEventListener(listener EventListener) {
  if b, ok := a.(interface { HandleEvent(event Event) }); ok {
    b.HandleEvent(event)
  } else if b, ok := a.(func(event Event)); ok {
    b(event)
  } else {
    panic("listener is not interface { HandleEvent(event Event) } or func(event Event)")
  }
}
```

**Why not use something like [`net/http.HandlerFunc`](https://pkg.go.dev/net/http#HandlerFunc)?** Because that would mean two types. You'd have to explicitly wrap your `func(event Event) {}` callback in an `EventListenerFunc()` conversion.

```go
type EventListener interface {
  HandleEvent(event Event)
}
type EventListenerFunc func(event Event)
func (e EventListenerFunc) HandleEvent(event Event) {
  e(event)
}
func AddEventListener(listener EventListener) {}
func main() {
  // Ugh.
  AddEventListener(EventListenerFunc(func(event Event) {
    fmt.Println(event)
  }))
}
```

This is better for type-safety. I could be convinced that this way is better. 🤷‍♀️

#### Dictionary types

Web IDL interfaces map quite well to Go `struct`. The biggest obstacle here is optional `struct` fields. Since Go doesn't have a Rust `Option<T>` it means we need to use something else. The most Go-like way of doing things is with `nil`-able types. For interfaces this is no sweat: interfaces are always `nil`-able. For primitives it means passing by `*string` instead of `string`.

Note: Go doesn't let you to take the `&` address of literal values. You can use a helper function like this for that:

```go
func ptr[T any](v T) *T {
  return &v
}
func DoThing(v *string) {}
func main() {
  // DoThing(&"no") // 😢
  DoThing(ptr("yes")) // 😄
}
```

**Example:**

```webidl
dictionary ResponseInit {
  unsigned short status = 200;
  ByteString statusText = "";
  HeadersInit headers;
};
```

```go
type ResponseInit struct {
  Status *uint16
  StatusText *string
  Headers any
}
```

Dictionaries that extend other dictionaries should use Go's struct embedding. This is a bit clunky from a consumer's perspective since some fields are nested behind long names but it's conceptually the most similar mapping of this Web IDL concept onto Go types.

**Example:**

```webidl
dictionary AddEventListenerOptions : EventListenerOptions {
  boolean passive = false;
  boolean once = false;
};
dictionary EventListenerOptions {
  boolean capture = false;
};
```

```go
type AddEventListenerOptions struct {
  EventListenerOptions // 👈 Embedded struct!
  Passive *bool
  Once *bool
}
type EventListenerOptions struct {
  Capture *bool
}
```

```go
AddEventListener("click", listener, AddEventListenerOptions{
  EventListenerOptions: EventListenerOptions{
    Capture: ptr(true),
  },
  Passive: ptr(true),
  Once: ptr(true),
})
```

Maybe flattened structs are better. Tell me your experience!

#### Enumeration types



### Callback function types ### {#go-callback-function}

Use a `type TheType func(...)...` newtype.

<div class=example>
```webidl
callback TheCallback = long (long a, long b);
```
```go
type TheCallback func(a int32, b int32) int32

func useTheCallback(cb TheCallback) {
  fmt.Printf("Result: %d", cb(1, 2))
}

func main() {
  useTheCallback(func(a int32, b int32) int32 {
    return a + b
  })
}
```
</div>

### Nullable types ### {#go-nullable-type}

Nullable types become `nil`-able types in Go. All Web IDL interface types should be `nil`-able by default since all Go `interface` types are `nil`-able. Go primitives like `string`, `int32`, etc. can all be made `nil`-able through a pointer. Go structs can be made `nil`-able through a pointer as well.

<div class=example>
```webidl
partial interface Window {
  undefined doThing(DOMString? message);
};
```
```go
func DoThing(message *string) {
  fmt.Println(message)
}
```
</div>

### Sequences ### {#go-sequence}

For now Web IDL sequences are analagous to `Array<T>` types. Each `sequence<T>` type can be exactly mapped into a Go `[]T` slice.

In the future iterables might be accepted.

<div class=example>
```webidl
partial interface Window {
  undefined doThing(sequence<DOMString> items);
};
```
```go
func DoThing(items []string) {
  for i, v := range items {
    fmt.Printf("item %d: %s", i, v)
  }
}
```
</div>

### Records ### {#go-record}

`record<K, V>` maps to `map[K]V` in Go.

<div class=example>
Usually records will use some kind of string key like {{DOMString}} or {{USVString}}.

```webidl
partial interface Window {
  undefined doThing(record<DOMString, long> things);
};
```
```go
func DoThing(things map[string]int32) {
  for k, v := range things {
    fmt.Println(k, v)
  }
}
```
</div>

So far all use-cases for `record<K, V>` have been **unordered** so to Go random `map[K]V` ordering is OK.

### Promise types ### {#go-promise}

Promises are a bit tricky in Go. Unlike `<-chan T`, a `Promise<T>` can be `await`-ed multiple times and will always spit out the same cached value. What this means is that just `<-chan T` isn't enough to represent a `Promise<T>` in Go. It falls apart specifically when you have a static `Promise<T>` that you want to `await` multiple times like {{ServiceWorkerContainer/ready}}.

Instead, the solution is to use a getter function. The current idiom that seems favored by Go is to *return something with a `.Wait()` method* and then call that to "await" the function's result.

<pre class=simpledef>
Promise&lt;T>: interface{ Wait() (T, error) }
Promise&lt;undefined>: interface{ Wait() error }
Infalliable Promise&lt;T>: interface{ Wait() T }
Infalliable Promise&lt;undefined>: interface{ Wait() }
</pre>

<div class=example>
```webidl
partial interface Window {
  Promise<Response> fetch();
};
```
```go
func Fetch() interface{Wait() (Response, error)} {
  // ...
}
```
</div>

An <dfn export>infalliable promise</dfn> is a promise that never errors. These are **extremely rare** in Web APIs since most things that are async call system-level operations that can often fail for arbitrary reasons. If you see a `Promise<T>` in Web IDL it's almost always a `Promise<T>` which should be a `interface { Wait() (T, error) }` in Go.

<div class=example>
{{ServiceWorkerContainer/ready}} is an example of the rare case where a `Promise<T>` will never reject.
```webidl
partial interface ServiceWorkerContainer {
  readonly attribute Promise<ServiceWorkerRegistration> ready;
};
```
```go
type ServiceWorkerContainer interface {
  Ready() interface{ Wait() ServiceWorkerRegistration }
}
```
</div>

In all cases the returned object should support being called concurrently. `sync.*` friends handle this for you.

The work that the operation is doing should be done in a `go func(){}()` Go-routine in an async fashion. The `Promise<T>` exposes the result to the caller.

### Union ### {#js-union}

Web IDL union types are represented as `any` in Go. Then at runtime they are `if b, ok := a.(T); ok {}` checked to see if they fit the union terminal types. There's no recursion for union types that contain union types; all of said types must be individually checked for at runtime using type conversion checks.

You should provide a stub `type UnionType = any` **type alias** (not a newtype) with a documentation comment outlining which values are accepted in the union.

<div class=note>

Why not `type UnionType any` as a newtype? Because then it becomes an opaque type that is impossible to downcast to any of the terminal types.

```webidl
typedef (DOMString or long) TheUnion;
partial interface Window {
  undefined doThing(TheUnion u);
};
```
```go
type TheUnion any
func DoThing(u TheUnion) {
  if uString, ok := u.(string); ok {
    fmt.Println("got string")
  } else if uInt32, ok := u.(int32); ok {
    fmt.Println("got int32")
  } else {
    panic("not string|int32")
  }
}

func main() {
  DoThing("wasd") // will panic!
  DoThing(1) // will panic!

  var myVariable TheUnion = "qwerty"
  // ☝ This has a type of TheUnion, not string.

  // ✅
  _ = myVariable.(TheUnion)

  // ❌
  _ = myVariable.(string)
}
```

There's not a good way to get the `string` or `int32` value out of `TheUnion` in this example when using `type U any`. If you can find one let me know! I'm sure there's a way.

On a semantics note: the union type should not behave as a distinct type; it should act like the _union_ of all the terminal types. `type U = any` is the best fit for that.

</div>

### Buffer source types ### {#go-buffer-source-types}

<pre class=simpledef>
{{Uint8Array}}: []uint8
{{Uint8ClampedArray}}: []uint8
{{Int8Array}}: []int8
{{Uint16Array}}: []uint16
{{Int16Array}}: []int16
{{Uint32Array}}: []uint32
{{Int32Array}}: []int32
{{BigUint64Array}}: []uint64
{{BigInt64Array}}: []int64
{{Float32Array}}: []float32
{{Float64Array}}: []float64
</pre>

An {{ArrayBuffer}} needs to support `[AllowResizable]` so it's represented in Go as a [bytes.Buffer](https://pkg.go.dev/bytes#Buffer). {{SharedArrayBuffer}} is a bit trickier. There's no concurrent `bytes.Buffer` in the Go standard library (that I know of) so for now it's just `bytes.Buffer` with the note that it should be used with caution.

<div class=example>
```webidl
partial interface Window {
  Float64Array doThing(ArrayBuffer b);
};
```
```go
func DoThing(b *bytes.Buffer) []float64 {
  return unsafe.Slice((*float64)(unsafe.Pointer(&b.Bytes()[0])), len(b.Bytes())/unsafe.Sizeof(float64(0)))
}
```
</div>

A {{DataView}} is roughly equivalent to a `[]byte` and can then be used with other encoding functions to use and mutate the slice's contents.

### Frozen arrays ### {#go-frozen-array}

Use a normal Go slice and just make sure you emphasize that it's intended to be *immutable*.

### Observable arrays ### {#go-observable-array}

This is fairly rare. {{DocumentOrShadowRoot/adoptedStyleSheets}} is one case where I know this is used. Don't know how to model this in Go.

## Extended attributes ## {#go-extended-attributes}

<pre class=simpledef>
[AllowResizable]: Add doc comment to note that this bytes.Buffer instance can be resized.
[AllowShared]: Make sure to emphasize that this bytes.Buffer should be behind a concurrency lock.
[Clamp]: Not relevant with exactly matching number types.
[CrossOriginIsolated]: Assume that we are always cross-origin isolated. 🤷‍♀️
[Default]: I don't know what this does.
[EnforceRange]: Not relevant. Go has exactly matching number types.
[Exposed]: Not relevant. every Web IDL definition is exposed to other Go code.
[Global]: Any interface that has this should expose all its members in the package scope and **not provide a concrete type itself**.
[NewObject]: Must always return a new instance such that `f() != f()`.
[Replaceable]: Not relevant. Go has a stricter type system.
[SameObject]: Must always return the same instance such that `f() == f()`.
[SecureContext]: Assume that we are always in a secure context. 🤷‍♀️
[Unscopable]: Not relevant. Go has no `with` block.
</pre>

<small>`-tags=coi` might be something that happens in the future to disable cross-origin isolated features by default.</small>

### [PutForwards] ### {#PutForwards}

<div class=note>

This attribute was designed so that this would work:

```webidl
[Exposed=Window]
interface Name {
  attribute DOMString full;
  attribute DOMString family;
  attribute DOMString given;
};

[Exposed=Window]
interface Person {
  [PutForwards=full] readonly attribute Name name;
  attribute unsigned short age;
};
```
```js
var p = getPerson();           // Obtain an instance of Person.

p.name = 'John Citizen';       // This statement...
p.name.full = 'John Citizen';  // ...has the same behavior as this one.
```
</div>

Since this attribute can only be applied on `readonly` attributes, the port to Go is simple: delegate the setter of the parent attribute to the specified child attribute. There's no union type or any type deduction.

<div class=example>
```webidl
[Exposed=Window]
interface Name {
  attribute DOMString full;
  attribute DOMString family;
  attribute DOMString given;
};

[Exposed=Window]
interface Person {
  [PutForwards=full] readonly attribute Name name;
  attribute unsigned short age;
};
```
```go
type Name interface {
  Full() string
  SetFull(v string)
  Family() string
  SetFamily(v string)
  Given() string
  SetGiven(v string)
  isName()
}

type Person interface {
  Name() Name
  SetName(v string)
  Age() uint16
  SetAge(v uint16)
  isPerson()
}

// ✅ The getter stays the same as always. This is a readonly property.
func (p *personImpl) Name() Name {
  return p.name
}

// 👇 We delegate to the p.name.SetFull() just like [PutForwards=full] said!
func (p *personImpl) SetName(v string) {
  p.name.SetFull(v)
}
```
</div>

<p class=note>
This attribute is used in {{ElementCSSInlineStyle}}/{{ElementCSSInlineStyle/style}}.
</p>

## Legacy extended attributes ## {#go-legacy-extended-attributes}

<pre class=simpledef>
[LegacyFactoryFunction]: Define this function in the same way JavaScript would.
[LegacyLenientSetter]: Not relevant.
[LegacyLenientThis]: Not relevant.
[LegacyNamespace]: Interfaces will have the namespace name prefixed to their name. Used extensively in the WebAssembly APIs.
[LegacyNoInterfaceObject]: Make the type private. Used in WebGL specifications.
[LegacyNullToEmptyString]: Not implemented.
[LegacyOverrideBuiltIns]: Not implemented.
[LegacyTreatNonObjectAsNull]: Not implemented.
[LegacyUnenumerableNamedProperties]: Not relevant. Go has no concept of iterable object properties.
[LegacyUnforgeable]: Not relevant. All Go properties are unforgeable (not overridable) by default.
[LegacyWindowAlias]: Do a `type T = U` type alias. Usually used for vendor-specific prefixes.
</pre>

## Overloads ## {#go-overloads}

Normally Web IDL overloads would be routed under a single JavaScript function. In Go the overload convention is different; you define overloads as separate functions with slightly different names. The naming convention for these extra functions usually focuses either on a type like `Match` vs `MatchString` or a named parameter like `Split` vs `SplitAfterN`.

The convention that this standard has adopted is to use TODO

## Interfaces ## {#go-interfaces}

TODO

### Constants ### {#go-constants}

TODO

### Attributes ### {#go-attributes}

### Stringifiers ### {#go-stringifier}

### Iterable declarations ### {#go-iterable}

### Asynchronous iterable declarations ### {#go-asynchronous-iterable}

### Maplike declarations ### {#go-maplike}

### Setlike declarations ### {#go-setlike}

## Namespaces ## {#go-namespaces}

## Exceptions ## {#go-exceptions}
