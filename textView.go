package rui

import (
	"fmt"
	"strings"
)

// TextView - text View
type TextView interface {
	View
}

type textViewData struct {
	viewData
	// TODO textShadow
}

// NewTextView create new TextView object and return it
func NewTextView(session Session, params Params) TextView {
	view := new(textViewData)
	view.Init(session)
	setInitParams(view, params)
	return view
}

func newTextView(session Session) View {
	return NewTextView(session, nil)
}

// Init initialize fields of TextView by default values
func (textView *textViewData) Init(session Session) {
	textView.viewData.Init(session)
	textView.tag = "TextView"
}

func (textView *textViewData) Get(tag string) interface{} {
	return textView.get(strings.ToLower(tag))
}

func (textView *textViewData) Remove(tag string) {
	textView.remove(strings.ToLower(tag))
}

func (textView *textViewData) remove(tag string) {
	textView.viewData.remove(tag)
	switch tag {
	case Text:
		updateInnerHTML(textView.htmlID(), textView.session)

	case TextOverflow:
		textView.textOverflowUpdated()
	}
}

func (textView *textViewData) Set(tag string, value interface{}) bool {
	return textView.set(strings.ToLower(tag), value)
}

func (textView *textViewData) set(tag string, value interface{}) bool {
	switch tag {
	case Text:
		switch value := value.(type) {
		case string:
			textView.properties[Text] = value

		case fmt.Stringer:
			textView.properties[Text] = value.String()

		case float32:
			textView.properties[Text] = fmt.Sprintf("%g", float64(value))

		case float64:
			textView.properties[Text] = fmt.Sprintf("%g", value)

		case []rune:
			textView.properties[Text] = string(value)

		case bool:
			if value {
				textView.properties[Text] = "true"
			} else {
				textView.properties[Text] = "false"
			}

		default:
			if n, ok := isInt(value); ok {
				textView.properties[Text] = fmt.Sprintf("%d", n)
			} else {
				notCompatibleType(tag, value)
				return false
			}
		}
		updateInnerHTML(textView.htmlID(), textView.session)
		return true

	case TextOverflow:
		if textView.viewData.set(tag, value) {
			textView.textOverflowUpdated()
		}
	}

	return textView.viewData.set(tag, value)
}

func (textView *textViewData) textOverflowUpdated() {
	session := textView.Session()
	if n, ok := enumProperty(textView, TextOverflow, textView.session, 0); ok {
		values := enumProperties[TextOverflow].cssValues
		if n >= 0 && n < len(values) {
			updateCSSProperty(textView.htmlID(), TextOverflow, values[n], session)
			return
		}
	}
	updateCSSProperty(textView.htmlID(), TextOverflow, "", session)
}

func (textView *textViewData) htmlSubviews(self View, buffer *strings.Builder) {
	if value, ok := stringProperty(textView, Text, textView.Session()); ok {
		if !GetNotTranslate(textView, "") {
			value, _ = textView.session.GetString(value)
		}

		text := strings.ReplaceAll(value, `"`, `\"`)
		text = strings.ReplaceAll(text, "\n", `\n`)
		text = strings.ReplaceAll(text, "\r", `\r`)
		buffer.WriteString(strings.ReplaceAll(text, `'`, `\'`))
	}
}

// GetTextOverflow returns a value of the "text-overflow" property:
// TextOverflowClip (0) or TextOverflowEllipsis (1).
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetTextOverflow(view View, subviewID string) int {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view == nil {
		return SingleLineText
	}
	t, _ := enumStyledProperty(view, TextOverflow, SingleLineText)
	return t
}