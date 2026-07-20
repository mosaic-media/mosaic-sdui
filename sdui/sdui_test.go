package sdui_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/mosaic-media/mosaic-sdui/sdui"
)

// homeScreen is the kind of tree the Platform's emit-side will build.
func homeScreen() sdui.Node {
	return sdui.Screen(
		sdui.Child(
			sdui.HeroBanner("Spirited Away",
				sdui.Meta("2001", "Anime Film", "PG"),
				sdui.Overview("A young girl wanders into a world of spirits."),
				sdui.Slot("actions",
					sdui.Button("Play", "primary", sdui.Play("demo-part")),
					sdui.Button("Details", "secondary", sdui.Navigate("detail", map[string]any{"title": "Spirited Away"})),
				),
			),
			sdui.Section("Continue watching",
				sdui.Child(sdui.Carousel(sdui.Child(
					sdui.PosterCard("Cowboy Bebop", "Anime Series",
						sdui.Subtitle("S1 · E12"), sdui.Progress(0.6), sdui.BadgeText("12 min left"),
						sdui.Act(sdui.Navigate("detail", map[string]any{"title": "Cowboy Bebop"}))),
					sdui.PosterCard("Dune", "Film", sdui.Progress(0.75)),
				))),
			),
		),
	)
}

func TestHomeScreenMarshals(t *testing.T) {
	b, err := json.Marshal(homeScreen())
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	s := string(b)
	for _, want := range []string{
		`"type":"Screen"`,
		`"type":"HeroBanner"`,
		`"type":"Carousel"`,
		`"type":"PosterCard"`,
		`"kind":"playPart"`,
		`"partId":"demo-part"`,
		`"progress":0.6`,
		`"slots":{"actions":`,
	} {
		if !strings.Contains(s, want) {
			t.Errorf("marshalled tree missing %q\n%s", want, s)
		}
	}
}

func TestActionsAreCleanPerKind(t *testing.T) {
	cases := map[string]sdui.Action{
		`{"kind":"navigate","screen":"home"}`:                      sdui.Navigate("home", nil),
		`{"kind":"playPart","partId":"p1"}`:                        sdui.Play("p1"),
		`{"kind":"toast","message":"hi","tone":"success"}`:         sdui.Toast("hi", sdui.ToneSuccess),
		`{"kind":"invoke","mutation":"importContent"}`:             sdui.Invoke("importContent", nil),
		`{"kind":"back"}`:                                          sdui.Back(),
	}
	for want, a := range cases {
		b, err := json.Marshal(a)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		if string(b) != want {
			t.Errorf("action marshalled to %s, want %s", b, want)
		}
	}
}

func TestRoundTrip(t *testing.T) {
	in := homeScreen()
	b, _ := json.Marshal(in)
	var out sdui.Node
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.Type != "Screen" || len(out.Children) != 2 {
		t.Fatalf("round-trip lost structure: type=%q children=%d", out.Type, len(out.Children))
	}
}
