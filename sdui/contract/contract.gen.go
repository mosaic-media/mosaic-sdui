// Code generated from schema/sdui.schema.json by quicktype. DO NOT EDIT.
package contract

// The single source of truth for the Mosaic Server-Driven-UI contract. Language bindings
// (Go, TypeScript, Dart) are GENERATED from this file — do not hand-edit them. The root
// object exists only so a generator reaches every top-level type; the useful types are in
// $defs.

// A declarative behaviour envelope. Data, never code — the client interprets the kind. Each
// kind uses a subset of the fields.
type Action struct {
	Actions  []Action               `json:"actions,omitempty"`
	Input    map[string]interface{} `json:"input,omitempty"`
	Kind     ActionKind             `json:"kind"`
	Message  *string                `json:"message,omitempty"`
	Mutation *string                `json:"mutation,omitempty"`
	Node     *UINode                `json:"node,omitempty"`
	NodeID   *string                `json:"nodeId,omitempty"`
	Params   map[string]interface{} `json:"params,omitempty"`
	PartID   *string                `json:"partId,omitempty"`
	Screen   *string                `json:"screen,omitempty"`
	Surface  *Surface               `json:"surface,omitempty"`
	Tone     *Tone                  `json:"tone,omitempty"`
	URL      *string                `json:"url,omitempty"`
}

// One element of a server-driven UI tree. The `type` is an open vocabulary: a client that
// does not recognise a type renders a placeholder rather than failing.
type UINode struct {
	Children []UINode `json:"children,omitempty"`
	ID       *string  `json:"id,omitempty"`
	// Component-specific data. Open by design.
	Props map[string]interface{} `json:"props,omitempty"`
	Slots map[string][]UINode    `json:"slots,omitempty"`
	// Component discriminator, e.g. "PosterCard".
	Type string `json:"type"`
}

// A component expressed as data: a name, default params, and a template of primitives.
// Clients register definitions and expand them; this is how a module contributes a
// component without shipping client code. A template node's props may hold binding objects
// ({"$bind":"path"} / {"$match":{…}}) and control keys ($if / $ifNot / $each / $as); a node
// of type "Outlet" renders the caller's children or a named slot.
type ComponentDefinition struct {
	// The node type this definition provides.
	Name string `json:"name"`
	// Default param values, overridden by the caller's props.
	Params   map[string]interface{} `json:"params,omitempty"`
	Template UINode                 `json:"template"`
}

type ActionKind string

const (
	Back         ActionKind = "back"
	CloseOverlay ActionKind = "closeOverlay"
	Invoke       ActionKind = "invoke"
	Navigate     ActionKind = "navigate"
	OpenOverlay  ActionKind = "openOverlay"
	OpenURL      ActionKind = "openUrl"
	PlayPart     ActionKind = "playPart"
	Sequence     ActionKind = "sequence"
	Toast        ActionKind = "toast"
)

type Surface string

const (
	Drawer Surface = "drawer"
	Modal  Surface = "modal"
	Sheet  Surface = "sheet"
)

type Tone string

const (
	Accent  Tone = "accent"
	Danger  Tone = "danger"
	Info    Tone = "info"
	Neutral Tone = "neutral"
	Success Tone = "success"
	Warning Tone = "warning"
)
