package rui

import "strings"

const (
	// KeyDown is the constant for "key-down-event" property tag.
	// The "key-down-event" event is fired when a key is pressed.
	// The main listener format:
	//   func(View, KeyEvent).
	// The additional listener formats:
	//   func(KeyEvent), func(View), and func().
	KeyDownEvent = "key-down-event"

	// KeyPp is the constant for "key-up-event" property tag.
	// The "key-up-event" event is fired when a key is released.
	// The main listener format:
	//   func(View, KeyEvent).
	// The additional listener formats:
	//   func(KeyEvent), func(View), and func().
	KeyUpEvent = "key-up-event"
)

type ControlKeyMask int
type KeyCode string

const (
	// AltKey is the mask of the "alt" key
	AltKey ControlKeyMask = 1
	// CtrlKey is the mask of the "ctrl" key
	CtrlKey ControlKeyMask = 2
	// ShiftKey is the mask of the "shift" key
	ShiftKey ControlKeyMask = 4
	// MetaKey is the mask of the "meta" key
	MetaKey ControlKeyMask = 8

	KeyA              KeyCode = "KeyA"
	KeyB              KeyCode = "KeyB"
	KeyC              KeyCode = "KeyC"
	KeyD              KeyCode = "KeyD"
	KeyE              KeyCode = "KeyE"
	KeyF              KeyCode = "KeyF"
	KeyG              KeyCode = "KeyG"
	KeyH              KeyCode = "KeyH"
	KeyI              KeyCode = "KeyI"
	KeyJ              KeyCode = "KeyJ"
	KeyK              KeyCode = "KeyK"
	KeyL              KeyCode = "KeyL"
	KeyM              KeyCode = "KeyM"
	KeyN              KeyCode = "KeyN"
	KeyO              KeyCode = "KeyO"
	KeyP              KeyCode = "KeyP"
	KeyQ              KeyCode = "KeyQ"
	KeyR              KeyCode = "KeyR"
	KeyS              KeyCode = "KeyS"
	KeyT              KeyCode = "KeyT"
	KeyU              KeyCode = "KeyU"
	KeyV              KeyCode = "KeyV"
	KeyW              KeyCode = "KeyW"
	KeyX              KeyCode = "KeyX"
	KeyY              KeyCode = "KeyY"
	KeyZ              KeyCode = "KeyZ"
	Digit0Key         KeyCode = "Digit0"
	Digit1Key         KeyCode = "Digit1"
	Digit2Key         KeyCode = "Digit2"
	Digit3Key         KeyCode = "Digit3"
	Digit4Key         KeyCode = "Digit4"
	Digit5Key         KeyCode = "Digit5"
	Digit6Key         KeyCode = "Digit6"
	Digit7Key         KeyCode = "Digit7"
	Digit8Key         KeyCode = "Digit8"
	Digit9Key         KeyCode = "Digit9"
	SpaceKey          KeyCode = "Space"
	MinusKey          KeyCode = "Minus"
	EqualKey          KeyCode = "Equal"
	IntlBackslashKey  KeyCode = "IntlBackslash"
	BracketLeftKey    KeyCode = "BracketLeft"
	BracketRightKey   KeyCode = "BracketRight"
	SemicolonKey      KeyCode = "Semicolon"
	CommaKey          KeyCode = "Comma"
	PeriodKey         KeyCode = "Period"
	QuoteKey          KeyCode = "Quote"
	BackquoteKey      KeyCode = "Backquote"
	SlashKey          KeyCode = "Slash"
	EscapeKey         KeyCode = "Escape"
	EnterKey          KeyCode = "Enter"
	TabKey            KeyCode = "Tab"
	CapsLockKey       KeyCode = "CapsLock"
	DeleteKey         KeyCode = "Delete"
	InsertKey         KeyCode = "Insert"
	HelpKey           KeyCode = "Help"
	BackspaceKey      KeyCode = "Backspace"
	PrintScreenKey    KeyCode = "PrintScreen"
	ScrollLockKey     KeyCode = "ScrollLock"
	PauseKey          KeyCode = "Pause"
	ContextMenuKey    KeyCode = "ContextMenu"
	ArrowLeftKey      KeyCode = "ArrowLeft"
	ArrowRightKey     KeyCode = "ArrowRight"
	ArrowUpKey        KeyCode = "ArrowUp"
	ArrowDownKey      KeyCode = "ArrowDown"
	HomeKey           KeyCode = "Home"
	EndKey            KeyCode = "End"
	PageUpKey         KeyCode = "PageUp"
	PageDownKey       KeyCode = "PageDown"
	F1Key             KeyCode = "F1"
	F2Key             KeyCode = "F2"
	F3Key             KeyCode = "F3"
	F4Key             KeyCode = "F4"
	F5Key             KeyCode = "F5"
	F6Key             KeyCode = "F6"
	F7Key             KeyCode = "F7"
	F8Key             KeyCode = "F8"
	F9Key             KeyCode = "F9"
	F10Key            KeyCode = "F10"
	F11Key            KeyCode = "F11"
	F12Key            KeyCode = "F12"
	F13Key            KeyCode = "F13"
	NumLockKey        KeyCode = "NumLock"
	NumpadKey0        KeyCode = "Numpad0"
	NumpadKey1        KeyCode = "Numpad1"
	NumpadKey2        KeyCode = "Numpad2"
	NumpadKey3        KeyCode = "Numpad3"
	NumpadKey4        KeyCode = "Numpad4"
	NumpadKey5        KeyCode = "Numpad5"
	NumpadKey6        KeyCode = "Numpad6"
	NumpadKey7        KeyCode = "Numpad7"
	NumpadKey8        KeyCode = "Numpad8"
	NumpadKey9        KeyCode = "Numpad9"
	NumpadDecimalKey  KeyCode = "NumpadDecimal"
	NumpadEnterKey    KeyCode = "NumpadEnter"
	NumpadAddKey      KeyCode = "NumpadAdd"
	NumpadSubtractKey KeyCode = "NumpadSubtract"
	NumpadMultiplyKey KeyCode = "NumpadMultiply"
	NumpadDivideKey   KeyCode = "NumpadDivide"
	ShiftLeftKey      KeyCode = "ShiftLeft"
	ShiftRightKey     KeyCode = "ShiftRight"
	ControlLeftKey    KeyCode = "ControlLeft"
	ControlRightKey   KeyCode = "ControlRight"
	AltLeftKey        KeyCode = "AltLeft"
	AltRightKey       KeyCode = "AltRight"
	MetaLeftKey       KeyCode = "MetaLeft"
	MetaRightKey      KeyCode = "MetaRight"
)

type KeyEvent struct {
	// TimeStamp is the time at which the event was created (in milliseconds).
	// This value is time since epoch—but in reality, browsers' definitions vary.
	TimeStamp uint64

	// Key is the key value of the key represented by the event. If the value has a printed representation,
	// this attribute's value is the same as the char property. Otherwise, it's one of the key value strings
	// specified in Key values. If the key can't be identified, its value is the string "Unidentified".
	Key string

	// Code holds a string that identifies the physical key being pressed. The value is not affected
	// by the current keyboard layout or modifier state, so a particular key will always return the same value.
	Code KeyCode

	// Repeat == true if a key has been depressed long enough to trigger key repetition, otherwise false.
	Repeat bool

	// CtrlKey == true if the control key was down when the event was fired. false otherwise.
	CtrlKey bool

	// ShiftKey == true if the shift key was down when the event was fired. false otherwise.
	ShiftKey bool

	// AltKey == true if the alt key was down when the event was fired. false otherwise.
	AltKey bool

	// MetaKey == true if the meta key was down when the event was fired. false otherwise.
	MetaKey bool
}

func (event *KeyEvent) init(data DataObject) {
	getBool := func(tag string) bool {
		if value, ok := data.PropertyValue(tag); ok && value == "1" {
			return true
		}
		return false
	}

	event.Key, _ = data.PropertyValue("key")
	code, _ := data.PropertyValue("code")
	event.Code = KeyCode(code)
	event.TimeStamp = getTimeStamp(data)
	event.Repeat = getBool("repeat")
	event.CtrlKey = getBool("ctrlKey")
	event.ShiftKey = getBool("shiftKey")
	event.AltKey = getBool("altKey")
	event.MetaKey = getBool("metaKey")
}

func valueToEventListeners[V View, E any](value any) ([]func(V, E), bool) {
	if value == nil {
		return nil, true
	}

	switch value := value.(type) {
	case func(V, E):
		return []func(V, E){value}, true

	case func(E):
		fn := func(_ V, event E) {
			value(event)
		}
		return []func(V, E){fn}, true

	case func(V):
		fn := func(view V, _ E) {
			value(view)
		}
		return []func(V, E){fn}, true

	case func():
		fn := func(V, E) {
			value()
		}
		return []func(V, E){fn}, true

	case []func(V, E):
		if len(value) == 0 {
			return nil, true
		}
		for _, fn := range value {
			if fn == nil {
				return nil, false
			}
		}
		return value, true

	case []func(E):
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(V, E), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			listeners[i] = func(_ V, event E) {
				v(event)
			}
		}
		return listeners, true

	case []func(V):
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(V, E), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			listeners[i] = func(view V, _ E) {
				v(view)
			}
		}
		return listeners, true

	case []func():
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(V, E), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			listeners[i] = func(V, E) {
				v()
			}
		}
		return listeners, true

	case []any:
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(V, E), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			switch v := v.(type) {
			case func(V, E):
				listeners[i] = v

			case func(E):
				listeners[i] = func(_ V, event E) {
					v(event)
				}

			case func(V):
				listeners[i] = func(view V, _ E) {
					v(view)
				}

			case func():
				listeners[i] = func(V, E) {
					v()
				}

			default:
				return nil, false
			}
		}
		return listeners, true
	}

	return nil, false
}

func valueToEventWithOldListeners[V View, E any](value any) ([]func(V, E, E), bool) {
	if value == nil {
		return nil, true
	}

	switch value := value.(type) {
	case func(V, E, E):
		return []func(V, E, E){value}, true

	case func(V, E):
		fn := func(v V, val, _ E) {
			value(v, val)
		}
		return []func(V, E, E){fn}, true

	case func(E, E):
		fn := func(_ V, val, old E) {
			value(val, old)
		}
		return []func(V, E, E){fn}, true

	case func(E):
		fn := func(_ V, val, _ E) {
			value(val)
		}
		return []func(V, E, E){fn}, true

	case func(V):
		fn := func(v V, _, _ E) {
			value(v)
		}
		return []func(V, E, E){fn}, true

	case func():
		fn := func(V, E, E) {
			value()
		}
		return []func(V, E, E){fn}, true

	case []func(V, E, E):
		if len(value) == 0 {
			return nil, true
		}
		for _, fn := range value {
			if fn == nil {
				return nil, false
			}
		}
		return value, true

	case []func(V, E):
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(V, E, E), count)
		for i, fn := range value {
			if fn == nil {
				return nil, false
			}
			listeners[i] = func(view V, val, _ E) {
				fn(view, val)
			}
		}
		return listeners, true

	case []func(E):
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(V, E, E), count)
		for i, fn := range value {
			if fn == nil {
				return nil, false
			}
			listeners[i] = func(_ V, val, _ E) {
				fn(val)
			}
		}
		return listeners, true

	case []func(E, E):
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(V, E, E), count)
		for i, fn := range value {
			if fn == nil {
				return nil, false
			}
			listeners[i] = func(_ V, val, old E) {
				fn(val, old)
			}
		}
		return listeners, true

	case []func(V):
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(V, E, E), count)
		for i, fn := range value {
			if fn == nil {
				return nil, false
			}
			listeners[i] = func(view V, _, _ E) {
				fn(view)
			}
		}
		return listeners, true

	case []func():
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(V, E, E), count)
		for i, fn := range value {
			if fn == nil {
				return nil, false
			}
			listeners[i] = func(V, E, E) {
				fn()
			}
		}
		return listeners, true

	case []any:
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(V, E, E), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			switch fn := v.(type) {
			case func(V, E, E):
				listeners[i] = fn

			case func(V, E):
				listeners[i] = func(view V, val, _ E) {
					fn(view, val)
				}

			case func(E, E):
				listeners[i] = func(_ V, val, old E) {
					fn(val, old)
				}

			case func(E):
				listeners[i] = func(_ V, val, _ E) {
					fn(val)
				}

			case func(V):
				listeners[i] = func(view V, _, _ E) {
					fn(view)
				}

			case func():
				listeners[i] = func(V, E, E) {
					fn()
				}

			default:
				return nil, false
			}
		}
		return listeners, true
	}

	return nil, false
}

func (view *viewData) setKeyListener(tag string, value any) bool {
	listeners, ok := valueToEventListeners[View, KeyEvent](value)
	if !ok {
		notCompatibleType(tag, value)
		return false
	}

	if listeners == nil {
		view.removeKeyListener(tag)
	} else {
		switch tag {
		case KeyDownEvent:
			view.properties.Store(tag, listeners)
			if view.created {
				view.session.updateProperty(view.htmlID(), "onkeydown", "keyDownEvent(this, event)")
			}

		case KeyUpEvent:
			view.properties.Store(tag, listeners)
			if view.created {
				view.session.updateProperty(view.htmlID(), "onkeyup", "keyUpEvent(this, event)")
			}

		default:
			return false
		}
	}

	return true
}

func (view *viewData) removeKeyListener(tag string) {
	view.properties.Delete(tag)
	if view.created {
		switch tag {
		case KeyDownEvent:
			if !view.Focusable() {
				view.session.removeProperty(view.htmlID(), "onkeydown")
			}

		case KeyUpEvent:
			view.session.removeProperty(view.htmlID(), "onkeyup")
		}
	}
}

func getEventWithOldListeners[V View, E any](view View, subviewID []string, tag string) []func(V, E, E) {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view != nil {
		if value := view.Get(tag); value != nil {
			if result, ok := value.([]func(V, E, E)); ok {
				return result
			}
		}
	}
	return []func(V, E, E){}
}

func getEventListeners[V View, E any](view View, subviewID []string, tag string) []func(V, E) {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view != nil {
		if value := view.Get(tag); value != nil {
			if result, ok := value.([]func(V, E)); ok {
				return result
			}
		}
	}
	return []func(V, E){}
}

func keyEventsHtml(view View, buffer *strings.Builder) {
	if len(getEventListeners[View, KeyEvent](view, nil, KeyDownEvent)) > 0 {
		buffer.WriteString(`onkeydown="keyDownEvent(this, event)" `)
	} else if view.Focusable() {
		if len(getEventListeners[View, MouseEvent](view, nil, ClickEvent)) > 0 {
			buffer.WriteString(`onkeydown="keyDownEvent(this, event)" `)
		}
	}

	if listeners := getEventListeners[View, KeyEvent](view, nil, KeyUpEvent); len(listeners) > 0 {
		buffer.WriteString(`onkeyup="keyUpEvent(this, event)" `)
	}
}

func handleKeyEvents(view View, tag string, data DataObject) {
	var event KeyEvent
	event.init(data)
	listeners := getEventListeners[View, KeyEvent](view, nil, tag)

	if len(listeners) > 0 {
		for _, listener := range listeners {
			listener(view, event)
		}
		return
	}

	if tag == KeyDownEvent && view.Focusable() && (event.Key == " " || event.Key == "Enter") && !IsDisabled(view) {
		if listeners := getEventListeners[View, MouseEvent](view, nil, ClickEvent); len(listeners) > 0 {
			clickEvent := MouseEvent{
				TimeStamp: event.TimeStamp,
				Button:    PrimaryMouseButton,
				Buttons:   PrimaryMouseMask,
				CtrlKey:   event.CtrlKey,
				AltKey:    event.AltKey,
				ShiftKey:  event.ShiftKey,
				MetaKey:   event.MetaKey,
				ClientX:   view.Frame().Width / 2,
				ClientY:   view.Frame().Height / 2,
				X:         view.Frame().Width / 2,
				Y:         view.Frame().Height / 2,
				ScreenX:   view.Frame().Left + view.Frame().Width/2,
				ScreenY:   view.Frame().Top + view.Frame().Height/2,
			}
			for _, listener := range listeners {
				listener(view, clickEvent)
			}
		}
	}
}

// GetKeyDownListeners returns the "key-down-event" listener list. If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetKeyDownListeners(view View, subviewID ...string) []func(View, KeyEvent) {
	return getEventListeners[View, KeyEvent](view, subviewID, KeyDownEvent)
}

// GetKeyUpListeners returns the "key-up-event" listener list. If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetKeyUpListeners(view View, subviewID ...string) []func(View, KeyEvent) {
	return getEventListeners[View, KeyEvent](view, subviewID, KeyUpEvent)
}
