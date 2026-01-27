# Hyper Control

Hyper Control's primary purpose is to provide human-readable XML wrappers to a REST interface, to guide the user in how to interact with the interface. Similar to HTML, but stripped down for readability; generally on the commandline using a tool like [HTTPie](https://httpie.io/docs/cli).

An XML element under the `c:` namespace is a hypermedia control provided by `hyprctl` that tells the user what interactions on the given resource are possible.

There are five main elements that hyprctl provides. They generally try to mimic the existing standard of HTML:

 - `<c:Form>`: analogous to HTML's `<form>`. It encloses a group of inputs, and generally describes which HTTP verb to use under the `method` attribute, e.g. `POST` or by default `GET`.
 - `<c:Input>`: analogous to HTML's `<input>` type. It represents a single name-value pair. It may have validation attributes, similar to HTML: e.g. `type`, `required`, `minlength`. An `<c:Input>` (or a `<c:Select>`, or a `<c:Map>`) outside of a `<c:Form>` is not a valid input.
 - `<c:Select>`: analogous to HTML's `<select>`. It represents an input with fixed options. The option list may be non-exaustive. It may take multiple options if the `multiple` attribute is set.
 - `<c:Link>`: analogous to HTML's `<a>` hyperlink. It provides directions to other relevant resources.
 - `<c:Map>`: Is the only element without an HTML analogue. A `<c:Map>` with `name="foo"` means that arbitrary name-value pairs may be provided under the namespace "foo" with bracket notation, e.g. `foo[bar]=baz`.

These elements can also be serialised to JSON for ease of querying, especially using [`jq`](https://jqlang.org/).
