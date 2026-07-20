package sdui

// Action constructors. The generated Action uses pointer fields for optionals
// (so the wire omits absent ones); these constructors hide that.

func strp(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// Navigate pushes another server-defined screen.
func Navigate(screen string, params map[string]any) Action {
	return Action{Kind: KindNavigate, Screen: strp(screen), Params: params}
}

// Back pops the client's navigation stack.
func Back() Action { return Action{Kind: KindBack} }

// OpenURL opens an external URL (the client validates the scheme).
func OpenURL(url string) Action { return Action{Kind: KindOpenURL, URL: strp(url)} }

// Invoke runs a Platform mutation by name.
func Invoke(mutation string, input map[string]any) Action {
	return Action{Kind: KindInvoke, Mutation: strp(mutation), Input: input}
}

// Query runs a Platform query, optionally refreshing a named region.
func Query(query string, variables map[string]any, into string) Action {
	return Action{Kind: KindQuery, Query: strp(query), Variables: variables, Into: strp(into)}
}

// OpenOverlay presents a node as a modal/sheet/drawer.
func OpenOverlay(surface Surface, node Node) Action {
	return Action{Kind: KindOpenOverlay, Surface: &surface, Node: &node}
}

// CloseOverlay dismisses the topmost overlay.
func CloseOverlay() Action { return Action{Kind: KindCloseOverlay} }

// Play asks the client to resolve and play a content Part.
func Play(partID string) Action { return Action{Kind: KindPlayPart, PartID: strp(partID)} }

// Toast shows a transient message.
func Toast(message string, tone Tone) Action {
	a := Action{Kind: KindToast, Message: strp(message)}
	if tone != "" {
		a.Tone = &tone
	}
	return a
}

// Sequence runs several actions in order.
func Sequence(actions ...Action) Action {
	return Action{Kind: KindSequence, Actions: actions}
}
