// Package sdui is the Go binding of the Mosaic Server-Driven-UI contract — the
// producer side. The Platform and Modules build a tree of Nodes carrying Action
// envelopes; a client (the Shell, a native client) renders it. The wire form is
// JSON; these types marshal to exactly the shape a client expects.
//
// This package is the normative Go view of the schema in ../schema. It is
// deliberately thin: a Node's props are an open map (the vocabulary is open —
// a client that does not know a type renders a placeholder), and the standard
// components are constructor helpers (components.go), not closed message types.
package sdui

// Node is one element of a server-driven UI tree.
type Node struct {
	// Type is the component discriminator, e.g. "PosterCard". Open vocabulary.
	Type string `json:"type"`
	// ID is an optional stable identity (React key / targeted updates).
	ID string `json:"id,omitempty"`
	// Props is component-specific data. Open by design.
	Props Props `json:"props,omitempty"`
	// Children are ordered child nodes (the common container case).
	Children []Node `json:"children,omitempty"`
	// Slots are named regions for components that take structured areas
	// (e.g. a HeroBanner's "actions").
	Slots map[string][]Node `json:"slots,omitempty"`
}

// Props is a component's open property bag.
type Props map[string]any

// Action is a declarative behaviour envelope. It is data, never code: a client
// interprets the Kind. The flat shape marshals to exactly the per-kind JSON a
// client expects — each kind uses a subset of the fields.
type Action struct {
	Kind string `json:"kind"`

	// navigate
	Screen string         `json:"screen,omitempty"`
	Params map[string]any `json:"params,omitempty"`
	// openUrl
	URL string `json:"url,omitempty"`
	// invoke
	Mutation string         `json:"mutation,omitempty"`
	Input    map[string]any `json:"input,omitempty"`
	// query
	Query     string         `json:"query,omitempty"`
	Variables map[string]any `json:"variables,omitempty"`
	Into      string         `json:"into,omitempty"`
	// openOverlay
	Surface string `json:"surface,omitempty"`
	Node    *Node  `json:"node,omitempty"`
	// playPart
	PartID string `json:"partId,omitempty"`
	NodeID string `json:"nodeId,omitempty"`
	// toast
	Message string `json:"message,omitempty"`
	Tone    string `json:"tone,omitempty"`
	// sequence
	Actions []Action `json:"actions,omitempty"`
}

// Action kinds.
const (
	KindNavigate     = "navigate"
	KindBack         = "back"
	KindOpenURL      = "openUrl"
	KindInvoke       = "invoke"
	KindQuery        = "query"
	KindOpenOverlay  = "openOverlay"
	KindCloseOverlay = "closeOverlay"
	KindPlayPart     = "playPart"
	KindToast        = "toast"
	KindSequence     = "sequence"
)

// Tones, shared by feedback components and toasts.
const (
	ToneNeutral = "neutral"
	ToneAccent  = "accent"
	ToneSuccess = "success"
	ToneWarning = "warning"
	ToneDanger  = "danger"
	ToneInfo    = "info"
)

// Overlay surfaces.
const (
	SurfaceModal  = "modal"
	SurfaceSheet  = "sheet"
	SurfaceDrawer = "drawer"
)

// ── Action constructors ──────────────────────────────────────────────────────

// Navigate pushes another server-defined screen.
func Navigate(screen string, params map[string]any) Action {
	return Action{Kind: KindNavigate, Screen: screen, Params: params}
}

// Back pops the client's navigation stack.
func Back() Action { return Action{Kind: KindBack} }

// OpenURL opens an external URL (the client validates the scheme).
func OpenURL(url string) Action { return Action{Kind: KindOpenURL, URL: url} }

// Invoke runs a Platform mutation by name.
func Invoke(mutation string, input map[string]any) Action {
	return Action{Kind: KindInvoke, Mutation: mutation, Input: input}
}

// Query runs a Platform query, optionally refreshing a named region.
func Query(query string, variables map[string]any, into string) Action {
	return Action{Kind: KindQuery, Query: query, Variables: variables, Into: into}
}

// OpenOverlay presents a node as a modal/sheet/drawer.
func OpenOverlay(surface string, node Node) Action {
	return Action{Kind: KindOpenOverlay, Surface: surface, Node: &node}
}

// CloseOverlay dismisses the topmost overlay.
func CloseOverlay() Action { return Action{Kind: KindCloseOverlay} }

// Play asks the client to resolve and play a content Part.
func Play(partID string) Action { return Action{Kind: KindPlayPart, PartID: partID} }

// Toast shows a transient message.
func Toast(message, tone string) Action {
	return Action{Kind: KindToast, Message: message, Tone: tone}
}

// Sequence runs several actions in order.
func Sequence(actions ...Action) Action {
	return Action{Kind: KindSequence, Actions: actions}
}
