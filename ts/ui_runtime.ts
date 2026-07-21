// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 the Mosaic authors

/**
 * The hand-written runtime for the TypeScript ui authoring layer — the Element
 * machinery and the control-flow/escape-hatch options (Prop, ID, Slot, Group,
 * When). The component constructors and the typed sugar (Hero, PosterCard,
 * OnTap, Meta, …) are generated into ui.ts from ui.spec.json by tools/genui;
 * edit the spec, not the generated file. Consumers import everything from the
 * package's `./ui` entry, which re-exports this runtime.
 */

import type { UINode } from "./contract.gen.js";

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

/** opt wraps a mutator as an El. Exported for the generated layer; not public. */
export function opt(fn: (e: Element) => void): El {
  return new Opt(fn);
}

/** compose runs a component: it applies every El to a fresh element of `typ`.
 * Exported for the generated layer; not part of the public API. */
export function compose(typ: string, base: Props | undefined, els: Elish[]): Element {
  const e = new Element(typ, base);
  for (const el of els) {
    if (el) el.applyTo(e);
  }
  return e;
}

function setProp(e: Element, key: string, val: unknown): void {
  e.props[key] = val;
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

/**
 * Group bundles several elements into one El, so a built array composes inline
 * alongside other elements — e.g. `Screen(Title(t), Group(...rows))`.
 */
export function Group(...els: Elish[]): El {
  return opt((e) => {
    for (const el of els) {
      if (el) el.applyTo(e);
    }
  });
}

/**
 * When includes `el` only if `cond` holds; otherwise it composes to nothing. (In
 * TS `cond && el` works too, since components skip falsy — When reads clearer.)
 */
export function When(cond: boolean, el: Elish): El {
  return cond && el ? el : opt(() => {});
}
