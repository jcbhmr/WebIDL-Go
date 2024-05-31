![🚧 Under construction 👷‍♂️](https://i.imgur.com/LEP2R3N.png)

- [x] Names
- [x] any
- [x] undefined
- [x] boolean
- [x] Integer types
- [x] Unrestricted float & double
- [x] Restricted float & double
- [x] bigint
- [x] DOMString
- [x] ByteString
- [x] USVString
- [x] object
- [ ] symbol
- [x] Interface types
- [x] Callback interface types
- [x] Dictionary types
- [ ] Enumeration types
- [x] Callback function types
- [x] Nullable types
- [x] Sequences
- [x] Records
- [x] Promise types
- [x] Union types
- [x] Buffer source types
- [x] Frozen arrays
- [ ] Observable arrays
- [ ] Extended attributes
- [ ] Legacy extended attributes
- [ ] Overloads
- [ ] Interface object
- [ ] Legacy factory functions
- [ ] Named properties object
- [ ] Constants
- [ ] Attributes
- [ ] Operations
- [ ] Stringifiers
- [ ] Iterable declarations
- [ ] Asynchronous iterable declarations
- [ ] Maplike declarations
- [ ] Setlike declarations
- [ ] Namespaces
- [ ] Exceptions

<p align=center>
  <b>You're probably looking for <a href="https://jcbhmr.me/WebIDL-Go/">jcbhmr.me/WebIDL-Go</a></b>
</p>

## Development

You'll need to have [Bikeshed](https://speced.github.io/bikeshed/) installed to build the spec. You can quickly install Bikeshed using `pipx`:

```sh
pipx install bikeshed
```

Once you have Bikeshed installed you can use `bikeshed serve` to start a local dev server.

```sh
bikeshed serve
```

Currently the dev server does not live reload on changes. [speced/bikeshed#1674](https://github.com/speced/bikeshed/issues/1674)

Here's a quick rundown of Bikeshed markup:

- You can use _basic_ Markdown syntax in most places. If you can't then just use HTML `<a>`, `<i>`, etc. 🤷‍♀️
- Use `<div class=example>` or `<p class=note>` to create automagic boxes with the fancy backgrounds.
- **`{{theThing}}`:** Autolink to the Web IDL definition of that symbol. Examples: `{{getSetCookie}}`, `{{Response/url}}`
- **`[[THE-SPEC inline]]`:** Automagically expands to the title of THE-SPEC. Examples: `[[FETCH inline]]`, `[[DOM inline]]`
