# Mosaic SDUI

The published **Server-Driven-UI contract** for [Mosaic](https://github.com/mosaic-media). It is to the interface what [`mosaic-sdk`](https://github.com/mosaic-media/mosaic-sdk) is to content: the single, language-neutral surface a **producer** (the Platform, or a Module) emits and a **client** (the Shell, a future Flutter client) renders — so nobody re-writes it per language.

> Status: seed. The schema and the Go producer binding are usable now; the standard-definition library and tokens are seeded and grow as the Shell is wired to consume them.

## Why this exists

Two Go producers already need it — the Platform must emit its screens, and a Module (e.g. `mosaic-module-stremio`) emits its own — while the Shell consumes it in TypeScript, and a native client will consume it in Dart. One contract, generated/mirrored per language, is the only way that stays consistent. This is [ADR 0023](https://github.com/mosaic-media/mosaic-architecture/blob/main/docs/adr/0023-server-driven-ui-and-the-shell.md)'s "extract on the second consumer," and the second consumer is here.

## What's in it

```
schema/         the normative contract, as JSON Schema (2020-12)
  uinode.schema.json        the UI node tree (open vocabulary)
  action.schema.json        the declarative Action envelope
  definition.schema.json    a component expressed as data (template + markers)
sdui/           Go binding — the PRODUCER side (Platform + Modules build trees)
definitions/    the standard component library, as data (a client registers these)
tokens/         design tokens (W3C DTCG) — compiled to CSS vars / Flutter theme
ts/             TypeScript binding — the CONSUMER side (the Shell)
```

**JSON, not protobuf** — deliberately. The node tree is open (props are an untyped bag by design), it rides GraphQL as JSON, and the definitions and tokens are JSON data. A JSON-native contract keeps all of that cohesive; protobuf's typed-message strength is undercut where the data actually lives. See [ADR 0025](https://github.com/mosaic-media/mosaic-architecture/blob/main/docs/adr/0025-sdui-contract-repository.md).

## Using it — a Go producer

```go
import "github.com/mosaic-media/mosaic-sdui/sdui"

home := sdui.Screen(sdui.Child(
    sdui.HeroBanner("Spirited Away",
        sdui.Meta("2001", "Anime Film", "PG"),
        sdui.Slot("actions",
            sdui.Button("Play", "primary", sdui.Play("part-1")),
        ),
    ),
    sdui.Section("Continue watching", sdui.Child(
        sdui.Carousel(sdui.Child(
            sdui.PosterCard("Cowboy Bebop", "Anime Series",
                sdui.Progress(0.6),
                sdui.Act(sdui.Navigate("detail", map[string]any{"title": "Cowboy Bebop"})),
            ),
        )),
    )),
))
// json.Marshal(home) → exactly the payload the Shell renders.
```

Until this module is tagged and published, a producer wires it locally with a
`replace github.com/mosaic-media/mosaic-sdui => ../mosaic-sdui` in its `go.mod`
(the same pattern the SDK uses for local work).

## The standard definitions

The reusable components — `PosterCard`, `HeroBanner`, `Section`, `Badge`, … — live here as `ComponentDefinition` data, not per-client code. A client registers them; a producer emits `{ "type": "HeroBanner", … }` and it renders identically on every client, with the Module shipping **zero** UI code. A Module can also ship its own definitions the same way.

Only the irreducible **primitives** (`Box`/`Text`/`Image`/`Pressable`/inputs/…) are per-client native code — the vocabulary each client implements. Definitions compose only those ([ADR 0024](https://github.com/mosaic-media/mosaic-architecture/blob/main/docs/adr/0024-primitives-and-definitions.md)).

## Next

- Wire the Shell to import `ts/types.ts`, load `definitions/*.json`, and generate its CSS variables from `tokens/tokens.json` (retiring its local copies).
- Migrate the remaining standard definitions from the Shell into `definitions/`.
- A tokens generator (DTCG → CSS + Dart) and the light theme.
- Tag `v0.1.0` and switch producers from `replace` to a versioned require.

## Licence

**Apache-2.0** (see [`LICENSE`](LICENSE) and [`NOTICE`](NOTICE)). A contract surface must be permissive so a Module may build its UI against it under any licence, exactly as the SDK is ([ADR 0022](https://github.com/mosaic-media/mosaic-architecture/blob/main/docs/adr/0022-licensing.md)).
