// Package sdui is the Go binding of the Mosaic Server-Driven-UI contract — the
// producer side. The Platform and Modules build a tree of Nodes carrying Action
// envelopes; a client renders it.
//
// The wire types are GENERATED from the single source of truth,
// schema/sdui.schema.json, into the contract subpackage — see contract.gen.go
// (do not hand-edit it; run scripts/generate.sh). This file provides the
// producer-facing ergonomics on top of the generated types: friendly aliases,
// constants, and the constructors in actions.go / components.go. A conformance
// test validates what the builders produce against the schema, so the
// hand-written layer cannot drift from the generated contract.
package sdui

//go:generate bash ../scripts/generate.sh

import "github.com/mosaic-media/sdui/sdui/contract"

// Contract types, re-exported from the generated package so producers import
// only "sdui".
type (
	Node                = contract.UINode
	Action              = contract.Action
	ActionKind          = contract.ActionKind
	Tone                = contract.Tone
	Surface             = contract.Surface
	ComponentDefinition = contract.ComponentDefinition
)

// Props is a component's open property bag.
type Props = map[string]any

// Action kinds (from the schema's ActionKind enum).
const (
	KindNavigate     = contract.Navigate
	KindBack         = contract.Back
	KindOpenURL      = contract.OpenURL
	KindInvoke       = contract.Invoke
	KindQuery        = contract.Query
	KindOpenOverlay  = contract.OpenOverlay
	KindCloseOverlay = contract.CloseOverlay
	KindPlayPart     = contract.PlayPart
	KindToast        = contract.Toast
	KindSequence     = contract.Sequence
)

// Tones.
const (
	ToneNeutral = contract.Neutral
	ToneAccent  = contract.Accent
	ToneSuccess = contract.Success
	ToneWarning = contract.Warning
	ToneDanger  = contract.Danger
	ToneInfo    = contract.Info
)

// Overlay surfaces.
const (
	SurfaceModal  = contract.Modal
	SurfaceSheet  = contract.Sheet
	SurfaceDrawer = contract.Drawer
)
