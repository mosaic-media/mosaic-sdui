// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 the Mosaic authors

// Package ui is a declarative authoring layer for the Mosaic SDUI — a proof of
// the "reads like a widget tree" ergonomics (Flutter/Compose/Vaadin, the
// gomponents pattern in Go). A component takes ...El, and children, props and
// slots are all Els that intermix, so a screen reads as a tree rather than a
// builder with option bags. The tree compiles to the same protobuf UINode the
// wire uses (ADR 0044) at Build(); only the authoring changes, not the payload.
package ui

import (
	"encoding/json"

	"google.golang.org/protobuf/types/known/structpb"

	sduiv1 "github.com/mosaic-media/sdui/gen/mosaic/sdui/v1"
	"github.com/mosaic-media/sdui/sdui"
)

// Node is the compiled, protobuf UI node the transport carries.
type Node = *sduiv1.UINode

// El is anything that composes into an element: a child, a prop, or a slot.
// A component accepts ...El and lets them intermix.
type El interface{ applyTo(*Element) }

// Element is a UI node under construction. It is itself an El — placing it in a
// parent adds it as a child — and Build() compiles it (and its subtree) to the
// protobuf Node.
type Element struct {
	typ      string
	id       string
	props    map[string]any
	children []*Element
	slots    map[string][]*Element
}

func (e *Element) applyTo(parent *Element) { parent.children = append(parent.children, e) }

// Build compiles the element tree into the protobuf UINode.
func (e *Element) Build() Node {
	if e == nil {
		return nil
	}
	n := &sduiv1.UINode{Type: e.typ, Id: e.id}
	if len(e.props) > 0 {
		n.Props = toStruct(e.props)
	}
	for _, c := range e.children {
		n.Children = append(n.Children, c.Build())
	}
	if len(e.slots) > 0 {
		n.Slots = make(map[string]*sduiv1.NodeList, len(e.slots))
		for name, kids := range e.slots {
			list := &sduiv1.NodeList{}
			for _, c := range kids {
				list.Nodes = append(list.Nodes, c.Build())
			}
			n.Slots[name] = list
		}
	}
	return n
}

// opt adapts a function into an El that modifies the element it lands in.
type opt func(*Element)

func (o opt) applyTo(e *Element) { o(e) }

// compose runs a component: it applies every El to a fresh element of typ.
func compose(typ string, base map[string]any, els []El) *Element {
	e := &Element{typ: typ, props: base}
	for _, el := range els {
		if el != nil {
			el.applyTo(e)
		}
	}
	return e
}

func setProp(e *Element, key string, val any) {
	if e.props == nil {
		e.props = map[string]any{}
	}
	e.props[key] = val
}

// ── containers ───────────────────────────────────────────────────────────────

// Screen is the root of a server-defined page.
func Screen(els ...El) *Element { return compose("Screen", nil, els) }

// Section is a titled band.
func Section(title string, els ...El) *Element {
	return compose("Section", map[string]any{"title": title}, els)
}

// Carousel is a horizontal snap-scrolling rail.
func Carousel(els ...El) *Element { return compose("Carousel", nil, els) }

// Grid is a responsive auto-fill grid.
func Grid(els ...El) *Element { return compose("Grid", nil, els) }

// Stack arranges children; direction is "horizontal" or "vertical".
func Stack(direction string, gap int, els ...El) *Element {
	return compose("Stack", map[string]any{"direction": direction, "gap": gap}, els)
}

// ── components ───────────────────────────────────────────────────────────────

// Hero is a featured banner. Fill its call-to-action row with Actions(…).
func Hero(title string, els ...El) *Element {
	return compose("HeroBanner", map[string]any{"title": title}, els)
}

// PosterCard renders a work/item card.
func PosterCard(title, mediaType string, els ...El) *Element {
	return compose("PosterCard", map[string]any{"title": title, "mediaType": mediaType}, els)
}

// Button carries an action; variant is primary/secondary/ghost/danger.
func Button(label, variant string, els ...El) *Element {
	return compose("Button", map[string]any{"label": label, "variant": variant}, els)
}

// Badge is a small pill; tone is one of the Tone constants.
func Badge(label, tone string, els ...El) *Element {
	return compose("Badge", map[string]any{"label": label, "tone": tone}, els)
}

// DetailHeader renders a node's metadata (title, meta, genres).
func DetailHeader(title string, els ...El) *Element {
	return compose("DetailHeader", map[string]any{"title": title}, els)
}

// EpisodeRow renders one episode under a season.
func EpisodeRow(title string, els ...El) *Element {
	return compose("EpisodeRow", map[string]any{"title": title}, els)
}

// PersonChip is a cast/crew chip.
func PersonChip(name string, els ...El) *Element {
	return compose("PersonChip", map[string]any{"name": name}, els)
}

// GenreTag is a genre chip.
func GenreTag(label string, els ...El) *Element {
	return compose("GenreTag", map[string]any{"label": label}, els)
}

// EmptyState is a titled empty placeholder.
func EmptyState(icon, title string) *Element {
	return compose("EmptyState", map[string]any{"icon": icon, "title": title}, nil)
}

// Component is the generic constructor for a type without a helper (a standard
// component like SeasonSelector, or a module's own).
func Component(typ string, els ...El) *Element { return compose(typ, nil, els) }

// ── slots ────────────────────────────────────────────────────────────────────

// Slot fills a named slot with the given elements' nodes.
func Slot(name string, els ...El) El {
	return opt(func(parent *Element) {
		scratch := &Element{}
		for _, el := range els {
			if el != nil {
				el.applyTo(scratch)
			}
		}
		if parent.slots == nil {
			parent.slots = map[string][]*Element{}
		}
		parent.slots[name] = append(parent.slots[name], scratch.children...)
	})
}

// Actions fills the "actions" slot (a hero's CTA row).
func Actions(els ...El) El { return Slot("actions", els...) }

// Aside fills the "aside" slot (a hero's docked poster).
func Aside(els ...El) El { return Slot("aside", els...) }

// ── prop options ─────────────────────────────────────────────────────────────

// Group bundles a slice of elements into one El, so a built slice composes
// inline alongside other elements — e.g. Screen(Title(t), Group(rows)).
func Group(els ...El) El {
	return opt(func(e *Element) {
		for _, el := range els {
			if el != nil {
				el.applyTo(e)
			}
		}
	})
}

// When includes el only if cond holds; otherwise it is a no-op. It lets a tree
// carry an element conditionally without breaking the declarative flow.
func When(cond bool, el El) El {
	if cond {
		return el
	}
	return opt(func(*Element) {})
}

// Prop sets an arbitrary prop — the escape hatch for anything without sugar.
func Prop(key string, val any) El { return opt(func(e *Element) { setProp(e, key, val) }) }

// ID sets a stable node id.
func ID(id string) El { return opt(func(e *Element) { e.id = id }) }

// OnTap sets the node's primary action.
func OnTap(a sdui.Action) El { return opt(func(e *Element) { setProp(e, "action", a) }) }

// Subtitle, Poster, Backdrop, Logo, Overview, Progress, BadgeText — card/hero sugar.
func Subtitle(s string) El   { return Prop("subtitle", s) }
func Poster(url string) El   { return Prop("poster", url) }
func Backdrop(url string) El { return Prop("backdrop", url) }
func Logo(url string) El     { return Prop("logo", url) }
func Overview(s string) El   { return Prop("overview", s) }
func Progress(f float64) El  { return Prop("progress", f) }
func BadgeText(s string) El  { return Prop("badge", s) }

// Title sets a screen or component title.
func Title(s string) El { return Prop("title", s) }

// Meta sets a hero's meta line (year · type · rating).
func Meta(items ...string) El { return Prop("meta", items) }

// Genres sets a DetailHeader's genre list.
func Genres(items ...string) El { return Prop("genres", items) }

// Tones (re-exported as strings, the open-bag encoding).
const (
	ToneNeutral = sdui.ToneNeutral
	ToneSuccess = sdui.ToneSuccess
	ToneWarning = sdui.ToneWarning
	ToneDanger  = sdui.ToneDanger
	ToneInfo    = sdui.ToneInfo
)

// ── re-exported actions (author with one import) ─────────────────────────────

// Action and its constructors come from the sdui producer package; they ride the
// open props bag as JSON (the faithful encoding, ADR 0044).
type Action = sdui.Action

var (
	Navigate = sdui.Navigate
	Invoke   = sdui.Invoke
	Play     = sdui.Play
)

// toStruct JSON-encodes the open props bag into a protobuf Struct.
func toStruct(props map[string]any) *structpb.Struct {
	b, err := json.Marshal(props)
	if err != nil {
		return nil
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil
	}
	s, err := structpb.NewStruct(m)
	if err != nil {
		return nil
	}
	return s
}
