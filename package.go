// Package hyprctl provides types for building HATEOAS-driven REST APIs with
// progressive enhancement.
//
// Each type represents a hypermedia control—a semantic element that guides
// client state transitions. These controls can be serialized as JSON, XML,
// or wrapped in HTML templates depending on the client's capabilities.
//
// # Philosophy
//
// Traditional REST APIs serve data (JSON) and leave state transitions to
// out-of-band documentation. HATEOAS APIs serve hypermedia—data enriched
// with the controls needed to interact with it. This library provides the
// building blocks for those controls.
//
// The types in hyprctl are deliberately minimal, focusing on the semantic
// layer rather than presentation. They describe WHAT actions are available
// and HOW to invoke them, not how they should be styled or rendered.
//
// # Progressive Enhancement
//
// The same endpoint can serve multiple representations:
//
//   - JSON: Machine-readable, suitable for scripts and traditional API clients
//   - XML: Human-readable structure for CLI tools (curl, HTTPie + xmllint)
//   - HTML: Browser-ready interfaces when wrapped with your templates
//
// This enables a "write once, serve many" approach where a single handler
// can satisfy browsers, CLIs, and programmatic clients through content
// negotiation.
//
// # Usage
//
// Types are designed to be embedded in your domain structs:
//
//	type LoginPage struct {
//	    LoginForm hyprctl.Form[struct {
//	        Username hyprctl.Input
//	        Password hyprctl.Input
//	        Submit   hyprctl.Submit
//	    }]
//	    RegisterLink hyprctl.Link
//	}
//
// Marshal to JSON for API clients, XML for CLI tools, or wrap in HTML
// templates for browsers. Examples at ./examples/templates.
//
// Input validation is minimal and extensible—Validate() checks Required
// and MinLength, matching basic browser behavior. Extend by inspecting
// Input.Value and setting Input.Error for domain-specific rules.
//
// # Pairs Well With
//
// This library pairs naturally with github.com/Teajey/rsvp, which provides
// HTTP handlers with ergonomic content negotiation, making it easy to serve
// the same semantic data in multiple formats based on Accept headers.
//
// # What This Is Not
//
// This is not a complete HTML generation framework. You bring your own
// templates for presentation (you should use the examples as a starting point).
// This is not a client-side form library—it's server-side semantics.
// This is not a validation framework—it provides minimal checks and extension
// points for your domain rules.
package hyprctl
