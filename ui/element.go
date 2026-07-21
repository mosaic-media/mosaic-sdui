// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 the Mosaic authors

// Package ui is a declarative authoring layer for the Mosaic SDUI — a "reads
// like a widget tree" ergonomics (Flutter/Compose/Vaadin, the gomponents
// pattern in Go). A component takes ...El, and children, props and slots are all
// Els that intermix, so a screen reads as a tree rather than a builder with
// option bags. The tree compiles to the same protobuf UINode the wire uses (ADR
// 0044) at Build(); only the authoring changes, not the payload.
//
// This file is the hand-written runtime — the Element machinery and the
// control-flow/escape-hatch options (Prop, ID, Slot, Group, When). The component
// constructors and the typed sugar (Hero, PosterCard, OnTap, Meta, …) are
// generated into components.gen.go from ui.spec.json by tools/genui; edit the
// spec, not the generated file.
package ui

import (
	"encoding/json"

	"google.golang.org/protobuf/types/known/structpb"

	sduiv1 "github.com/mosaic-media/sdui/gen/mosaic/sdui/v1"
	"github.com/mosaic-media/sdui/sdui"
)

// Node is the compiled, protobuf UI node the transport carries.
type Node = *sduiv1.UINode

// Action is a declarative behaviour envelope, re-exported from the producer
// binding; it rides the open props bag as JSON (ADR 0044). Author one with the
// generated constructors (Navigate, Invoke, Play).
type Action = sdui.Action

// El is anything that composes into an element: a child, a prop, or a slot. A
// component accepts ...El and lets them intermix.
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

// Prop sets an arbitrary prop — the escape hatch for anything without sugar.
func Prop(key string, val any) El { return opt(func(e *Element) { setProp(e, key, val) }) }

// ID sets a stable node id.
func ID(id string) El { return opt(func(e *Element) { e.id = id }) }

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

// Group bundles a slice of elements into one El, so a built slice composes
// inline alongside other elements — e.g. Screen(Title(t), Group(rows...)).
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
