// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 the Mosaic authors

package ui_test

import (
	"testing"

	"github.com/mosaic-media/contracts/ui"
)

// TestDeclarativeTreeBuilds proves the widget-tree authoring compiles to the
// expected protobuf structure — children as varargs, a slot, and an action in
// the open props bag.
func TestDeclarativeTreeBuilds(t *testing.T) {
	root := ui.Screen(
		ui.Hero("Spirited Away",
			ui.Overview("A young girl wanders into a world of spirits."),
			ui.Meta("2001", "Anime Film"),
			ui.Actions(
				ui.Button("Play", "primary", ui.OnTap(ui.Play("demo-part"))),
			),
		),
		ui.Section("Continue watching",
			ui.Carousel(
				ui.PosterCard("Cowboy Bebop", "Anime Series",
					ui.Subtitle("S1 · E12"), ui.Progress(0.6),
					ui.OnTap(ui.Navigate("detail", map[string]any{"title": "Cowboy Bebop"})),
				),
			),
		),
	).Build()

	if root.GetType() != "Screen" || len(root.GetChildren()) != 2 {
		t.Fatalf("root = %q with %d children, want Screen/2", root.GetType(), len(root.GetChildren()))
	}

	hero := root.GetChildren()[0]
	if hero.GetType() != "HeroBanner" || hero.GetProps().AsMap()["title"] != "Spirited Away" {
		t.Fatalf("hero = %+v", hero.GetProps().AsMap())
	}
	// The Play button lives in the actions slot, its action in the open props bag.
	acts := hero.GetSlots()["actions"].GetNodes()
	if len(acts) != 1 {
		t.Fatalf("hero actions = %d, want 1", len(acts))
	}
	action, _ := acts[0].GetProps().AsMap()["action"].(map[string]any)
	if action["kind"] != "playPart" || action["partId"] != "demo-part" {
		t.Fatalf("action = %v, want playPart/demo-part", action)
	}

	// The card carries a 0.6 progress and a navigate action.
	card := root.GetChildren()[1].GetChildren()[0].GetChildren()[0]
	if card.GetType() != "PosterCard" || card.GetProps().AsMap()["progress"] != 0.6 {
		t.Fatalf("card = %+v", card.GetProps().AsMap())
	}
}
