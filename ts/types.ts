// The TypeScript binding of the Mosaic SDUI contract — the consumer side.
// The Shell (and any future TS client) imports these instead of re-declaring
// them. Kept in exact correspondence with ../schema and the Go binding in
// ../sdui. A generator from the JSON Schema will eventually emit this file; it
// is committed for now so the Shell can depend on it directly.

/** One element of a server-driven UI tree. `type` is an open vocabulary. */
export interface UINode {
  type: string;
  id?: string;
  props?: Record<string, unknown>;
  children?: UINode[];
  slots?: Record<string, UINode | UINode[]>;
}

/** A declarative behaviour envelope — data, never code. */
export type Action =
  | { kind: "navigate"; screen: string; params?: Record<string, unknown> }
  | { kind: "back" }
  | { kind: "openUrl"; url: string }
  | { kind: "invoke"; mutation: string; input?: Record<string, unknown> }
  | { kind: "query"; query: string; variables?: Record<string, unknown>; into?: string }
  | { kind: "openOverlay"; surface?: "modal" | "sheet" | "drawer"; node: UINode }
  | { kind: "closeOverlay" }
  | { kind: "playPart"; partId: string; nodeId?: string }
  | { kind: "toast"; message: string; tone?: Tone }
  | { kind: "sequence"; actions: Action[] };

export type Tone = "neutral" | "accent" | "success" | "warning" | "danger" | "info";

/** The Platform's fixed error categories; the client maps failures into these. */
export type PlatformErrorCategory =
  | "InvalidArgument"
  | "Unauthenticated"
  | "PermissionDenied"
  | "NotFound"
  | "Conflict"
  | "Unavailable"
  | "Internal";

/** A component expressed as data: registered and expanded by the client. */
export interface ComponentDefinition {
  name: string;
  params?: Record<string, unknown>;
  template: UINode;
}
