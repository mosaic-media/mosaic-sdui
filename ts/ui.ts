// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 the Mosaic authors

/**
 * A declarative authoring layer for the Mosaic SDUI — the TypeScript twin of the
 * Go `ui` package. It proves the same "reads like a widget tree" ergonomics
 * (Flutter/Compose/Vaadin) on the TS side: a component takes `...El`, and
 * children, props and slots are all `El`s that intermix, so a screen reads as a
 * tree rather than a hand-written bag of `UINode` JSON. `build()` compiles the
 * tree to the same `UINode` the wire carries (ADR 0044) — only the authoring
 * changes, not the payload. It is meant for authoring mock/fixture screens (the
 * Shell's mock payloads, storybook stories) with the same shape a producer emits.
 *
 * The API deliberately mirrors the Go `ui` package name-for-name, so a screen
 * written in one language transliterates to the other.
 */

import { ActionKind, type Action, type Tone, type UINode } from "./contract.gen.js";

export { ActionKind, Surface, Tone } from "./contract.gen.js";
export type { Action, UINode } from "./contract.gen.js";

/** A component's open property bag. */
export type Props = Record<string, unknown>;

/**
 * El is anything that composes into an element: a child, a prop, or a slot. A
 * component accepts `...El` and lets them intermix. `false`/`null`/`undefined`
 * are accepted and skipped, so `cond && Prop(...)` composes inline.
 */
export interface El {
  applyTo(parent: Element): void;
}

/** Elish is an El or a falsy placeholder that composes to nothing. */
export type Elish = El | false | null | undefined;

/**
 * Element is a UI node under construction. It is itself an El — placing it in a
 * parent adds it as a child — and `build()` compiles it (and its subtree) to the
 * `UINode` JSON.
 */
export class Element implements El {
  typ: string;
  id = "";
  props: Props = {};
  children: Element[] = [];
  slots: Record<string, Element[]> = {};

  constructor(typ: string, base?: Props) {
    this.typ = typ;
    if (base) this.props = { ...base };
  }

  applyTo(parent: Element): void {
    parent.children.push(this);
  }

  /** build compiles the element tree into the UINode JSON the wire carries. */
  build(): UINode {
    const node: UINode = { type: this.typ };
    if (this.id) node.id = this.id;
    if (Object.keys(this.props).length > 0) node.props = this.props;
    if (this.children.length > 0) node.children = this.children.map((c) => c.build());
    const slotNames = Object.keys(this.slots);
    if (slotNames.length > 0) {
      node.slots = {};
      for (const name of slotNames) {
        node.slots[name] = this.slots[name].map((c) => c.build());
      }
    }
    return node;
  }
}

/** opt adapts a function into an El that modifies the element it lands in. */
class Opt implements El {
  constructor(private readonly fn: (e: Element) => void) {}
  applyTo(e: Element): void {
    this.fn(e);
  }
}

function opt(fn: (e: Element) => void): El {
  return new Opt(fn);
}

/** compose runs a component: it applies every El to a fresh element of `typ`. */
function compose(typ: string, base: Props | undefined, els: Elish[]): Element {
  const e = new Element(typ, base);
  for (const el of els) {
    if (el) el.applyTo(e);
  }
  return e;
}

function setProp(e: Element, key: string, val: unknown): void {
  e.props[key] = val;
}

// ── containers ───────────────────────────────────────────────────────────────

/** Screen is the root of a server-defined page. */
export function Screen(...els: Elish[]): Element {
  return compose("Screen", undefined, els);
}

/** Section is a titled band. */
export function Section(title: string, ...els: Elish[]): Element {
  return compose("Section", { title }, els);
}

/** Carousel is a horizontal snap-scrolling rail. */
export function Carousel(...els: Elish[]): Element {
  return compose("Carousel", undefined, els);
}

/** Grid is a responsive auto-fill grid. */
export function Grid(...els: Elish[]): Element {
  return compose("Grid", undefined, els);
}

/** Stack arranges children; direction is "horizontal" or "vertical". */
export function Stack(direction: string, gap: number, ...els: Elish[]): Element {
  return compose("Stack", { direction, gap }, els);
}

// ── components ───────────────────────────────────────────────────────────────

/** Hero is a featured banner. Fill its call-to-action row with Actions(…). */
export function Hero(title: string, ...els: Elish[]): Element {
  return compose("HeroBanner", { title }, els);
}

/** PosterCard renders a work/item card. */
export function PosterCard(title: string, mediaType: string, ...els: Elish[]): Element {
  return compose("PosterCard", { title, mediaType }, els);
}

/** Button carries an action; variant is primary/secondary/ghost/danger. */
export function Button(label: string, variant: string, ...els: Elish[]): Element {
  return compose("Button", { label, variant }, els);
}

/** Badge is a small pill; tone is one of the Tone values. */
export function Badge(label: string, tone: Tone | string, ...els: Elish[]): Element {
  return compose("Badge", { label, tone }, els);
}

/** DetailHeader renders a node's metadata (title, meta, genres). */
export function DetailHeader(title: string, ...els: Elish[]): Element {
  return compose("DetailHeader", { title }, els);
}

/** EpisodeRow renders one episode under a season. */
export function EpisodeRow(title: string, ...els: Elish[]): Element {
  return compose("EpisodeRow", { title }, els);
}

/** PersonChip is a cast/crew chip. */
export function PersonChip(name: string, ...els: Elish[]): Element {
  return compose("PersonChip", { name }, els);
}

/** GenreTag is a genre chip. */
export function GenreTag(label: string, ...els: Elish[]): Element {
  return compose("GenreTag", { label }, els);
}

/** EmptyState is a titled empty placeholder. */
export function EmptyState(icon: string, title: string): Element {
  return compose("EmptyState", { icon, title }, []);
}

/**
 * Component is the generic constructor for a type without a helper (a standard
 * component like SeasonSelector, or a module's own).
 */
export function Component(typ: string, ...els: Elish[]): Element {
  return compose(typ, undefined, els);
}

// ── slots ────────────────────────────────────────────────────────────────────

/** Slot fills a named slot with the given elements' nodes. */
export function Slot(name: string, ...els: Elish[]): El {
  return opt((parent) => {
    const scratch = new Element("");
    for (const el of els) {
      if (el) el.applyTo(scratch);
    }
    (parent.slots[name] ??= []).push(...scratch.children);
  });
}

/** Actions fills the "actions" slot (a hero's CTA row). */
export function Actions(...els: Elish[]): El {
  return Slot("actions", ...els);
}

/** Aside fills the "aside" slot (a hero's docked poster). */
export function Aside(...els: Elish[]): El {
  return Slot("aside", ...els);
}

// ── prop options ─────────────────────────────────────────────────────────────

/**
 * Group bundles several elements into one El, so a built array composes inline
 * alongside other elements — e.g. `Screen(Title(t), Group(rows))`.
 */
export function Group(...els: Elish[]): El {
  return opt((e) => {
    for (const el of els) {
      if (el) el.applyTo(e);
    }
  });
}

/**
 * When includes `el` only if `cond` holds; otherwise it composes to nothing. It
 * carries an element conditionally without breaking the declarative flow. (In TS
 * `cond && el` works too, since components skip falsy — When reads clearer.)
 */
export function When(cond: boolean, el: Elish): El {
  return cond && el ? el : opt(() => {});
}

/** Prop sets an arbitrary prop — the escape hatch for anything without sugar. */
export function Prop(key: string, val: unknown): El {
  return opt((e) => setProp(e, key, val));
}

/** ID sets a stable node id. */
export function ID(id: string): El {
  return opt((e) => {
    e.id = id;
  });
}

/** OnTap sets the node's primary action. */
export function OnTap(a: Action): El {
  return opt((e) => setProp(e, "action", a));
}

/** Subtitle, Poster, Backdrop, Logo, Overview, Progress, BadgeText — card/hero sugar. */
export function Subtitle(s: string): El {
  return Prop("subtitle", s);
}
export function Poster(url: string): El {
  return Prop("poster", url);
}
export function Backdrop(url: string): El {
  return Prop("backdrop", url);
}
export function Logo(url: string): El {
  return Prop("logo", url);
}
export function Overview(s: string): El {
  return Prop("overview", s);
}
export function Progress(f: number): El {
  return Prop("progress", f);
}
export function BadgeText(s: string): El {
  return Prop("badge", s);
}

/** Title sets a screen or component title. */
export function Title(s: string): El {
  return Prop("title", s);
}

/** Meta sets a hero's meta line (year · type · rating). */
export function Meta(...items: string[]): El {
  return Prop("meta", items);
}

/** Genres sets a DetailHeader's genre list. */
export function Genres(...items: string[]): El {
  return Prop("genres", items);
}

// ── actions (author with one import) ─────────────────────────────────────────
// Actions ride inside the open props bag as JSON (the faithful encoding, ADR
// 0044); each constructor hides the discriminator and the optional fields.

/** Navigate pushes a screen, optionally with params. */
export function Navigate(screen: string, params?: Props): Action {
  return { kind: ActionKind.Navigate, screen, ...(params ? { params } : {}) };
}

/** Invoke runs a server mutation, optionally with input. */
export function Invoke(mutation: string, input?: Props): Action {
  return { kind: ActionKind.Invoke, mutation, ...(input ? { input } : {}) };
}

/** Play starts playback of a part. */
export function Play(partId: string): Action {
  return { kind: ActionKind.PlayPart, partId };
}
