package rui

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// ListItemClickedEvent is the constant for "list-item-clicked" property tag.
	// The "list-item-clicked" event occurs when the user clicks on an item in the list.
	// The main listener format: func(ListView, int), where the second argument is the item index.
	ListItemClickedEvent = "list-item-clicked"
	// ListItemSelectedEvent is the constant for "list-item-selected" property tag.
	// The "list-item-selected" event occurs when a list item becomes selected.
	// The main listener format: func(ListView, int), where the second argument is the item index.
	ListItemSelectedEvent = "list-item-selected"
	// ListItemCheckedEvent is the constant for "list-item-checked" property tag.
	// The "list-item-checked" event occurs when a list item checkbox becomes checked/unchecked.
	// The main listener format: func(ListView, []int), where the second argument is the array of checked item indexes.
	ListItemCheckedEvent = "list-item-checked"
	// ListItemStyle is the constant for "list-item-style" property tag.
	// The "list-item-style" string property defines the style of an unselected item
	ListItemStyle = "list-item-style"
	// CurrentStyle is the constant for "current-style" property tag.
	// The "current-style" string property defines the style of the selected item when the ListView is focused.
	CurrentStyle = "current-style"
	// CurrentInactiveStyle is the constant for "current-inactive-style" property tag.
	// The "current-inactive-style" string property defines the style of the selected item when the ListView is unfocused.
	CurrentInactiveStyle = "current-inactive-style"
)

const (
	// VerticalOrientation is the vertical ListView orientation
	VerticalOrientation = 0
	// HorizontalOrientation is the horizontal ListView orientation
	HorizontalOrientation = 1

	// NoneCheckbox is value of "checkbox" property: no checkbox
	NoneCheckbox = 0
	// SingleCheckbox is value of "checkbox" property: only one item can be checked
	SingleCheckbox = 1
	// MultipleCheckbox is value of "checkbox" property: several items can be checked
	MultipleCheckbox = 2
)

// ListView - the list view interface
type ListView interface {
	View
	ParanetView
	// ReloadListViewData updates ListView content
	ReloadListViewData()

	getCheckedItems() []int
	getItemFrames() []Frame
}

type listViewData struct {
	viewData
	adapter           ListAdapter
	clickedListeners  []func(ListView, int)
	selectedListeners []func(ListView, int)
	checkedListeners  []func(ListView, []int)
	items             []View
	itemFrame         []Frame
	checkedItem       []int
}

// NewListView creates the new list view
func NewListView(session Session, params Params) ListView {
	view := new(listViewData)
	view.Init(session)
	setInitParams(view, params)
	return view
}

func newListView(session Session) View {
	return NewListView(session, nil)
}

// Init initialize fields of ViewsContainer by default values
func (listView *listViewData) Init(session Session) {
	listView.viewData.Init(session)
	listView.tag = "ListView"
	listView.systemClass = "ruiListView"
	listView.items = []View{}
	listView.itemFrame = []Frame{}
	listView.checkedItem = []int{}
	listView.clickedListeners = []func(ListView, int){}
	listView.selectedListeners = []func(ListView, int){}
	listView.checkedListeners = []func(ListView, []int){}
}

func (listView *listViewData) Views() []View {
	return listView.items
}

func (listView *listViewData) normalizeTag(tag string) string {
	tag = strings.ToLower(tag)
	switch tag {
	case HorizontalAlign:
		tag = ItemHorizontalAlign

	case VerticalAlign:
		tag = ItemVerticalAlign
	}
	return tag
}

func (listView *listViewData) Remove(tag string) {
	listView.remove(listView.normalizeTag(tag))
}

func (listView *listViewData) remove(tag string) {
	switch tag {
	case Checked:
		listView.checkedItem = []int{}
		updateInnerHTML(listView.htmlID(), listView.session)

	case Items:
		listView.adapter = nil
		updateInnerHTML(listView.htmlID(), listView.session)

	case Orientation, Wrap:
		delete(listView.properties, tag)
		updateCSSStyle(listView.htmlID(), listView.session)

	case Current:
		current := GetListViewCurrent(listView, "")
		delete(listView.properties, tag)
		updateInnerHTML(listView.htmlID(), listView.session)
		if current != -1 {
			for _, listener := range listView.selectedListeners {
				listener(listView, -1)
			}
		}

	case ItemWidth, ItemHeight, ItemHorizontalAlign, ItemVerticalAlign, ItemCheckbox,
		CheckboxHorizontalAlign, CheckboxVerticalAlign, ListItemStyle, CurrentStyle, CurrentInactiveStyle:

		delete(listView.properties, tag)
		updateInnerHTML(listView.htmlID(), listView.session)

	case ListItemClickedEvent:
		if len(listView.clickedListeners) > 0 {
			listView.clickedListeners = []func(ListView, int){}
		}

	case ListItemSelectedEvent:
		if len(listView.selectedListeners) > 0 {
			listView.selectedListeners = []func(ListView, int){}
		}

	case ListItemCheckedEvent:
		if len(listView.checkedListeners) > 0 {
			listView.checkedListeners = []func(ListView, []int){}
		}

	default:
		listView.viewData.remove(tag)
	}
}

func (listView *listViewData) Set(tag string, value interface{}) bool {
	return listView.set(listView.normalizeTag(tag), value)
}

func (listView *listViewData) set(tag string, value interface{}) bool {
	if value == nil {
		listView.remove(tag)
		return true
	}

	result := false

	switch tag {

	case ListItemClickedEvent:
		listeners := listView.valueToItemListeners(value)
		if listeners == nil {
			notCompatibleType(tag, value)
			return false
		}
		listView.clickedListeners = listeners
		return true

	case ListItemSelectedEvent:
		listeners := listView.valueToItemListeners(value)
		if listeners == nil {
			notCompatibleType(tag, value)
			return false
		}
		listView.selectedListeners = listeners
		return true

	case ListItemCheckedEvent:
		return listView.setItemCheckedEvent(value)

	case Checked:
		return listView.setChecked(value)

	case Items:
		result = listView.setItems(value)

	case Current:
		oldCurrent := GetListViewCurrent(listView, "")
		if listView.setIntProperty(Current, value) {
			current := GetListViewCurrent(listView, "")
			if oldCurrent != current {
				updateInnerHTML(listView.htmlID(), listView.session)
				for _, listener := range listView.selectedListeners {
					listener(listView, current)
				}
			}
			return true
		}

	case Orientation, Wrap:
		if listView.viewData.set(tag, value) {
			updateCSSStyle(listView.htmlID(), listView.session)
			return true
		}

	case ItemWidth, ItemHeight:
		result = listView.setSizeProperty(tag, value)

	case ItemHorizontalAlign, ItemVerticalAlign, ItemCheckbox, CheckboxHorizontalAlign, CheckboxVerticalAlign:
		result = listView.setEnumProperty(tag, value, enumProperties[tag].values)

	case ListItemStyle, CurrentStyle, CurrentInactiveStyle:
		switch value := value.(type) {
		case string:
			listView.properties[tag] = value
			result = true

		default:
			notCompatibleType(tag, value)
			return false
		}

	default:
		return listView.viewData.set(tag, value)
	}

	if result {
		updateInnerHTML(listView.htmlID(), listView.session)
	}

	return result
}

func (listView *listViewData) setItemCheckedEvent(value interface{}) bool {
	switch value := value.(type) {
	case func(ListView, []int):
		listView.checkedListeners = []func(ListView, []int){value}

	case func([]int):
		fn := func(view ListView, date []int) {
			value(date)
		}
		listView.checkedListeners = []func(ListView, []int){fn}

	case []func(ListView, []int):
		listView.checkedListeners = value

	case []func([]int):
		listeners := make([]func(ListView, []int), len(value))
		for i, val := range value {
			if val == nil {
				notCompatibleType(ListItemCheckedEvent, val)
				return false
			}

			listeners[i] = func(view ListView, date []int) {
				val(date)
			}
		}
		listView.checkedListeners = listeners

	case []interface{}:
		listeners := make([]func(ListView, []int), len(value))
		for i, val := range value {
			if val == nil {
				notCompatibleType(ListItemCheckedEvent, val)
				return false
			}

			switch val := val.(type) {
			case func(ListView, []int):
				listeners[i] = val

			case func([]int):
				listeners[i] = func(view ListView, checked []int) {
					val(checked)
				}

			default:
				notCompatibleType(ListItemCheckedEvent, val)
				return false
			}
		}
		listView.checkedListeners = listeners
	}
	return true
}

func (listView *listViewData) Get(tag string) interface{} {
	return listView.get(listView.normalizeTag(tag))
}

func (listView *listViewData) get(tag string) interface{} {
	switch tag {
	case ListItemClickedEvent:
		return listView.clickedListeners

	case ListItemSelectedEvent:
		return listView.selectedListeners

	case ListItemCheckedEvent:
		return listView.checkedListeners

	case Checked:
		return listView.checkedItem

	case Items:
		return listView.adapter

	case ListItemStyle:
		return listView.listItemStyle()

	case CurrentStyle:
		return listView.currentStyle()

	case CurrentInactiveStyle:
		return listView.currentInactiveStyle()
	}
	return listView.viewData.get(tag)
}

func (listView *listViewData) setItems(value interface{}) bool {
	switch value := value.(type) {
	case []string:
		listView.adapter = NewTextListAdapter(value, nil)

	case []DataValue:
		hasObject := false
		for _, val := range value {
			if val.IsObject() {
				hasObject = true
				break
			}
		}

		if hasObject {
			items := make([]View, len(value))
			for i, val := range value {
				if val.IsObject() {
					if view := CreateViewFromObject(listView.session, val.Object()); view != nil {
						items[i] = view
					} else {
						return false
					}
				} else {
					items[i] = NewTextView(listView.session, Params{Text: val.Value()})
				}
			}
			listView.adapter = NewViewListAdapter(items)
		} else {
			items := make([]string, len(value))
			for i, val := range value {
				items[i] = val.Value()
			}
			listView.adapter = NewTextListAdapter(items, nil)
		}

	case []interface{}:
		items := make([]View, len(value))
		for i, val := range value {
			switch value := val.(type) {
			case View:
				items[i] = value

			case string:
				items[i] = NewTextView(listView.session, Params{Text: value})

			case fmt.Stringer:
				items[i] = NewTextView(listView.session, Params{Text: value.String()})

			case float32:
				items[i] = NewTextView(listView.session, Params{Text: fmt.Sprintf("%g", float64(value))})

			case float64:
				items[i] = NewTextView(listView.session, Params{Text: fmt.Sprintf("%g", value)})

			default:
				if n, ok := isInt(val); ok {
					items[i] = NewTextView(listView.session, Params{Text: strconv.Itoa(n)})
				} else {
					notCompatibleType(Items, value)
					return false
				}
			}
		}
		listView.adapter = NewViewListAdapter(items)

	case []View:
		listView.adapter = NewViewListAdapter(value)

	case ListAdapter:
		listView.adapter = value

	default:
		notCompatibleType(Items, value)
		return false
	}

	size := listView.adapter.ListSize()
	listView.items = make([]View, size)
	listView.itemFrame = make([]Frame, size)

	return true
}

func (listView *listViewData) valueToItemListeners(value interface{}) []func(ListView, int) {
	if value == nil {
		return []func(ListView, int){}
	}

	switch value := value.(type) {
	case func(ListView, int):
		return []func(ListView, int){value}

	case func(int):
		fn := func(view ListView, index int) {
			value(index)
		}
		return []func(ListView, int){fn}

	case []func(ListView, int):
		return value

	case []func(int):
		listeners := make([]func(ListView, int), len(value))
		for i, val := range value {
			if val == nil {
				return nil
			}
			listeners[i] = func(view ListView, index int) {
				val(index)
			}
		}
		return listeners

	case []interface{}:
		listeners := make([]func(ListView, int), len(value))
		for i, val := range value {
			if val == nil {
				return nil
			}
			switch val := val.(type) {
			case func(ListView, int):
				listeners[i] = val

			case func(int):
				listeners[i] = func(view ListView, index int) {
					val(index)
				}

			default:
				return nil
			}
		}
		return listeners
	}

	return nil
}

func (listView *listViewData) setChecked(value interface{}) bool {
	var checked []int
	if value == nil {
		checked = []int{}
	} else {
		switch value := value.(type) {
		case int:
			checked = []int{value}

		case []int:
			checked = value

		default:
			return false
		}
	}

	switch GetListViewCheckbox(listView, "") {
	case SingleCheckbox:
		count := len(checked)
		if count > 1 {
			return false
		}

		if len(listView.checkedItem) > 0 &&
			(count == 0 || listView.checkedItem[0] != checked[0]) {
			listView.updateCheckboxItem(listView.checkedItem[0], false)
		}

		if count == 1 {
			listView.updateCheckboxItem(checked[0], true)
		}

	case MultipleCheckbox:
		inSlice := func(n int, slice []int) bool {
			for _, n2 := range slice {
				if n2 == n {
					return true
				}
			}
			return false
		}

		for _, n := range listView.checkedItem {
			if !inSlice(n, checked) {
				listView.updateCheckboxItem(n, false)
			}
		}

		for _, n := range checked {
			if !inSlice(n, listView.checkedItem) {
				listView.updateCheckboxItem(n, true)
			}
		}

	default:
		return false
	}

	listView.checkedItem = checked
	for _, listener := range listView.checkedListeners {
		listener(listView, listView.checkedItem)
	}
	return true
}

func (listView *listViewData) Focusable() bool {
	return true
}

func (listView *listViewData) ReloadListViewData() {
	itemCount := 0
	if listView.adapter != nil {
		itemCount = listView.adapter.ListSize()

		if itemCount != len(listView.items) {
			listView.items = make([]View, itemCount)
			listView.itemFrame = make([]Frame, itemCount)
		}

		for i := 0; i < itemCount; i++ {
			listView.items[i] = listView.adapter.ListItem(i, listView.Session())
		}
	} else if len(listView.items) > 0 {
		listView.items = []View{}
		listView.itemFrame = []Frame{}
	}

	updateInnerHTML(listView.htmlID(), listView.session)
}

func (listView *listViewData) getCheckedItems() []int {
	return listView.checkedItem
}

func (listView *listViewData) getItemFrames() []Frame {
	return listView.itemFrame
}

func (listView *listViewData) htmlProperties(self View, buffer *strings.Builder) {
	buffer.WriteString(`onfocus="listViewFocusEvent(this, event)" onblur="listViewBlurEvent(this, event)"`)
	buffer.WriteString(` onkeydown="listViewKeyDownEvent(this, event)" data-focusitemstyle="`)
	buffer.WriteString(listView.currentStyle())
	buffer.WriteString(`" data-bluritemstyle="`)
	buffer.WriteString(listView.currentInactiveStyle())
	buffer.WriteString(`"`)
	current := GetListViewCurrent(listView, "")
	if listView.adapter != nil && current >= 0 && current < listView.adapter.ListSize() {
		buffer.WriteString(` data-current="`)
		buffer.WriteString(listView.htmlID())
		buffer.WriteRune('-')
		buffer.WriteString(strconv.Itoa(current))
		buffer.WriteRune('"')
	}
}

func (listView *listViewData) cssStyle(self View, builder cssBuilder) {
	listView.viewData.cssStyle(self, builder)

	if GetListWrap(listView, "") != WrapOff {
		switch GetListOrientation(listView, "") {
		case TopDownOrientation, BottomUpOrientation:
			builder.add(`max-height`, `100%`)
		default:
			builder.add(`max-width`, `100%`)
		}
	}
}

func (listView *listViewData) itemAlign(self View, buffer *strings.Builder) {
	values := enumProperties[ItemHorizontalAlign].cssValues
	if hAlign := GetListItemHorizontalAlign(listView, ""); hAlign >= 0 && hAlign < len(values) {
		buffer.WriteString(" justify-items: ")
		buffer.WriteString(values[hAlign])
		buffer.WriteRune(';')
	}

	values = enumProperties[ItemVerticalAlign].cssValues
	if vAlign := GetListItemVerticalAlign(listView, ""); vAlign >= 0 && vAlign < len(values) {
		buffer.WriteString(" align-items: ")
		buffer.WriteString(values[vAlign])
		buffer.WriteRune(';')
	}
}

func (listView *listViewData) itemSize(self View, buffer *strings.Builder) {
	if itemWidth := GetListItemWidth(listView, ""); itemWidth.Type != Auto {
		buffer.WriteString(` min-width: `)
		buffer.WriteString(itemWidth.cssString(""))
		buffer.WriteRune(';')
	}

	if itemHeight := GetListItemHeight(listView, ""); itemHeight.Type != Auto {
		buffer.WriteString(` min-height: `)
		buffer.WriteString(itemHeight.cssString(""))
		buffer.WriteRune(';')
	}
}

func (listView *listViewData) getDivs(self View, checkbox, hCheckboxAlign, vCheckboxAlign int) (string, string, string) {
	session := listView.Session()

	contentBuilder := allocStringBuilder()
	defer freeStringBuilder(contentBuilder)

	contentBuilder.WriteString(`<div style="display: grid;`)
	listView.itemAlign(self, contentBuilder)

	onDivBuilder := allocStringBuilder()
	defer freeStringBuilder(onDivBuilder)

	if hCheckboxAlign == CenterAlign {
		if vCheckboxAlign == BottomAlign {
			onDivBuilder.WriteString(`<div style="grid-row: 2 / 3; grid-column: 1 / 2; display: grid; justify-items: center;`)
			contentBuilder.WriteString(` grid-row: 1 / 2; grid-column: 1 / 2;">`)
		} else {
			vCheckboxAlign = TopAlign
			onDivBuilder.WriteString(`<div style="grid-row: 1 / 2; grid-column: 1 / 2; display: grid; justify-items: center;`)
			contentBuilder.WriteString(` grid-row: 2 / 3; grid-column: 1 / 2;">`)
		}
	} else {
		if hCheckboxAlign == RightAlign {
			onDivBuilder.WriteString(`<div style="grid-row: 1 / 2; grid-column: 2 / 3; display: grid;`)
			contentBuilder.WriteString(` grid-row: 1 / 2; grid-column: 1 / 2;">`)
		} else {
			onDivBuilder.WriteString(`<div style="grid-row: 1 / 2; grid-column: 1 / 2; display: grid;`)
			contentBuilder.WriteString(` grid-row: 1 / 2; grid-column: 2 / 3;">`)
		}
		switch vCheckboxAlign {
		case BottomAlign:
			onDivBuilder.WriteString(` align-items: end;`)

		case CenterAlign:
			onDivBuilder.WriteString(` align-items: center;`)

		default:
			onDivBuilder.WriteString(` align-items: start;`)
		}
	}

	onDivBuilder.WriteString(`">`)

	offDivBuilder := allocStringBuilder()
	defer freeStringBuilder(offDivBuilder)

	offDivBuilder.WriteString(onDivBuilder.String())

	if checkbox == SingleCheckbox {
		offDivBuilder.WriteString(session.radiobuttonOffImage())
		onDivBuilder.WriteString(session.radiobuttonOnImage())
	} else {
		offDivBuilder.WriteString(session.checkboxOffImage())
		onDivBuilder.WriteString(session.checkboxOnImage())
	}

	onDivBuilder.WriteString("</div>")
	offDivBuilder.WriteString("</div>")

	return onDivBuilder.String(), offDivBuilder.String(), contentBuilder.String()
}

func (listView *listViewData) checkboxItemDiv(self View, checkbox, hCheckboxAlign, vCheckboxAlign int) string {
	itemStyleBuilder := allocStringBuilder()
	defer freeStringBuilder(itemStyleBuilder)

	itemStyleBuilder.WriteString(`<div style="display: grid; justify-items: stretch; align-items: stretch;`)

	if hCheckboxAlign == CenterAlign {
		if vCheckboxAlign == BottomAlign {
			itemStyleBuilder.WriteString(` grid-template-columns: 1fr; grid-template-rows: 1fr auto;`)
		} else {
			vCheckboxAlign = TopAlign
			itemStyleBuilder.WriteString(` grid-template-columns: 1fr; grid-template-rows: auto 1fr;`)
		}
	} else {
		if hCheckboxAlign == RightAlign {
			itemStyleBuilder.WriteString(` grid-template-columns: 1fr auto; grid-template-rows: 1fr;`)
		} else {
			itemStyleBuilder.WriteString(` grid-template-columns: auto 1fr; grid-template-rows: 1fr;`)
		}
	}

	if gap, ok := sizeConstant(listView.session, "ruiCheckboxGap"); ok && gap.Type != Auto {
		itemStyleBuilder.WriteString(` grid-gap: `)
		itemStyleBuilder.WriteString(gap.cssString("auto"))
		itemStyleBuilder.WriteRune(';')
	}

	itemStyleBuilder.WriteString(`">`)
	return itemStyleBuilder.String()

}

func (listView *listViewData) getItemView(index int) View {
	if listView.adapter == nil || index < 0 || index >= listView.adapter.ListSize() {
		return nil
	}

	size := listView.adapter.ListSize()
	if size != len(listView.items) {
		listView.items = make([]View, size)
	}

	if listView.items[index] == nil {
		listView.items[index] = listView.adapter.ListItem(index, listView.Session())
	}

	return listView.items[index]
}

func (listView *listViewData) listItemStyle() string {
	if value := listView.getRaw(ListItemStyle); value != nil {
		if style, ok := value.(string); ok {
			if style, ok = listView.session.resolveConstants(style); ok {
				return style
			}
		}
	}
	return "ruiListItem"
}

func (listView *listViewData) currentStyle() string {
	if value := listView.getRaw(CurrentStyle); value != nil {
		if style, ok := value.(string); ok {
			if style, ok = listView.session.resolveConstants(style); ok {
				return style
			}
		}
	}
	return "ruiListItemFocused"
}

func (listView *listViewData) currentInactiveStyle() string {
	if value := listView.getRaw(CurrentInactiveStyle); value != nil {
		if style, ok := value.(string); ok {
			if style, ok = listView.session.resolveConstants(style); ok {
				return style
			}
		}
	}
	return "ruiListItemSelected"
}

func (listView *listViewData) checkboxSubviews(self View, buffer *strings.Builder, checkbox int) {
	count := listView.adapter.ListSize()
	listViewID := listView.htmlID()

	hCheckboxAlign := GetListViewCheckboxHorizontalAlign(listView, "")
	vCheckboxAlign := GetListViewCheckboxVerticalAlign(listView, "")

	itemDiv := listView.checkboxItemDiv(self, checkbox, hCheckboxAlign, vCheckboxAlign)
	onDiv, offDiv, contentDiv := listView.getDivs(self, checkbox, hCheckboxAlign, vCheckboxAlign)

	current := GetListViewCurrent(listView, "")
	checkedItems := GetListViewCheckedItems(listView, "")
	for i := 0; i < count; i++ {
		buffer.WriteString(`<div id="`)
		buffer.WriteString(listViewID)
		buffer.WriteRune('-')
		buffer.WriteString(strconv.Itoa(i))
		buffer.WriteString(`" class="ruiView `)
		buffer.WriteString(listView.listItemStyle())
		if i == current {
			buffer.WriteRune(' ')
			buffer.WriteString(listView.currentInactiveStyle())
		}
		buffer.WriteString(`" onclick="listItemClickEvent(this, event)" data-left="0" data-top="0" data-width="0" data-height="0" style="display: grid; justify-items: stretch; align-items: stretch;`)
		listView.itemSize(self, buffer)
		buffer.WriteString(`">`)
		buffer.WriteString(itemDiv)

		checked := false
		for _, index := range checkedItems {
			if index == i {
				buffer.WriteString(onDiv)
				checked = true
				break
			}
		}
		if !checked {
			buffer.WriteString(offDiv)
		}
		buffer.WriteString(contentDiv)

		if view := listView.getItemView(i); view != nil {
			//view.setNoResizeEvent()
			viewHTML(view, buffer)
		} else {
			buffer.WriteString("ERROR: invalid item view")
		}

		buffer.WriteString(`</div></div></div>`)
	}
}

func (listView *listViewData) noneCheckboxSubviews(self View, buffer *strings.Builder) {
	count := listView.adapter.ListSize()
	listViewID := listView.htmlID()

	itemStyleBuilder := allocStringBuilder()
	defer freeStringBuilder(itemStyleBuilder)

	itemStyleBuilder.WriteString(`data-left="0" data-top="0" data-width="0" data-height="0" style="max-width: 100%; max-height: 100%; display: grid;`)

	listView.itemAlign(self, itemStyleBuilder)
	listView.itemSize(self, itemStyleBuilder)

	itemStyleBuilder.WriteString(`" onclick="listItemClickEvent(this, event)"`)
	itemStyle := itemStyleBuilder.String()

	current := GetListViewCurrent(listView, "")
	for i := 0; i < count; i++ {
		buffer.WriteString(`<div id="`)
		buffer.WriteString(listViewID)
		buffer.WriteRune('-')
		buffer.WriteString(strconv.Itoa(i))
		buffer.WriteString(`" class="ruiView `)
		buffer.WriteString(listView.listItemStyle())
		if i == current {
			buffer.WriteRune(' ')
			buffer.WriteString(listView.currentInactiveStyle())
		}
		buffer.WriteString(`" `)
		buffer.WriteString(itemStyle)
		buffer.WriteString(`>`)

		if view := listView.getItemView(i); view != nil {
			//view.setNoResizeEvent()
			viewHTML(view, buffer)
		} else {
			buffer.WriteString("ERROR: invalid item view")
		}

		buffer.WriteString(`</div>`)
	}
}

func (listView *listViewData) updateCheckboxItem(index int, checked bool) {

	checkbox := GetListViewCheckbox(listView, "")
	hCheckboxAlign := GetListViewCheckboxHorizontalAlign(listView, "")
	vCheckboxAlign := GetListViewCheckboxVerticalAlign(listView, "")
	onDiv, offDiv, contentDiv := listView.getDivs(listView, checkbox, hCheckboxAlign, vCheckboxAlign)

	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	buffer.WriteString(`updateInnerHTML('`)
	buffer.WriteString(listView.htmlID())
	buffer.WriteRune('-')
	buffer.WriteString(strconv.Itoa(index))
	buffer.WriteString(`', '`)

	buffer.WriteString(listView.checkboxItemDiv(listView, checkbox, hCheckboxAlign, vCheckboxAlign))
	if checked {
		buffer.WriteString(onDiv)
	} else {
		buffer.WriteString(offDiv)
	}
	buffer.WriteString(contentDiv)

	session := listView.Session()
	if listView.adapter != nil {
		if view := listView.getItemView(index); view != nil {
			view.setNoResizeEvent()
			viewHTML(view, buffer)
		} else {
			buffer.WriteString("ERROR: invalid item view")
		}
	}
	buffer.WriteString(`</div></div>');`)

	session.runScript(buffer.String())
}

func (listView *listViewData) htmlSubviews(self View, buffer *strings.Builder) {
	if listView.adapter == nil {
		return
	}
	if listView.adapter.ListSize() == 0 {
		return
	}

	if !listView.session.ignoreViewUpdates() {
		listView.session.setIgnoreViewUpdates(true)
		defer listView.session.setIgnoreViewUpdates(false)
	}

	checkbox := GetListViewCheckbox(listView, "")
	if checkbox == NoneCheckbox {
		listView.noneCheckboxSubviews(self, buffer)
	} else {
		listView.checkboxSubviews(self, buffer, checkbox)
	}
}

func (listView *listViewData) handleCommand(self View, command string, data DataObject) bool {
	switch command {
	case "itemSelected":
		if text, ok := data.PropertyValue(`number`); ok {
			if number, err := strconv.Atoi(text); err == nil {
				listView.properties[Current] = number
				for _, listener := range listView.selectedListeners {
					listener(listView, number)
				}
			}
		}

	case "itemUnselected":
		delete(listView.properties, Current)
		for _, listener := range listView.selectedListeners {
			listener(listView, -1)
		}

	case "itemClick":
		listView.onItemClick()

	default:
		return listView.viewData.handleCommand(self, command, data)
	}

	return true
}

func (listView *listViewData) onItemClick() {
	current := GetListViewCurrent(listView, "")
	if current >= 0 && !IsDisabled(listView) {
		checkbox := GetListViewCheckbox(listView, "")
	m:
		switch checkbox {
		case SingleCheckbox:
			if len(listView.checkedItem) == 0 {
				listView.checkedItem = []int{current}
				listView.updateCheckboxItem(current, true)
			} else if listView.checkedItem[0] != current {
				listView.updateCheckboxItem(listView.checkedItem[0], false)
				listView.checkedItem[0] = current
				listView.updateCheckboxItem(current, true)
			}

		case MultipleCheckbox:
			for i, index := range listView.checkedItem {
				if index == current {
					listView.updateCheckboxItem(index, false)
					count := len(listView.checkedItem)
					if count == 1 {
						listView.checkedItem = []int{}
					} else if i == 0 {
						listView.checkedItem = listView.checkedItem[1:]
					} else if i == count-1 {
						listView.checkedItem = listView.checkedItem[:i]
					} else {
						listView.checkedItem = append(listView.checkedItem[:i], listView.checkedItem[i+1:]...)
					}
					break m
				}
			}

			listView.updateCheckboxItem(current, true)
			listView.checkedItem = append(listView.checkedItem, current)
		}

		if checkbox != NoneCheckbox {
			for _, listener := range listView.checkedListeners {
				listener(listView, listView.checkedItem)
			}
		}
		for _, listener := range listView.clickedListeners {
			listener(listView, current)
		}
	}
}

func (listView *listViewData) onItemResize(self View, index int, x, y, width, height float64) {
	if index >= 0 && index < len(listView.itemFrame) {
		listView.itemFrame[index] = Frame{Left: x, Top: y, Width: width, Height: height}
	}
}

// GetVerticalAlign return the vertical align of a list: TopAlign (0), BottomAlign (1), CenterAlign (2), StretchAlign (3)
func GetVerticalAlign(view View) int {
	if align, ok := enumProperty(view, VerticalAlign, view.Session(), TopAlign); ok {
		return align
	}
	return TopAlign
}

// GetHorizontalAlign return the vertical align of a list: LeftAlign (0), RightAlign (1), CenterAlign (2), StretchAlign (3)
func GetHorizontalAlign(view View) int {
	if align, ok := enumProperty(view, HorizontalAlign, view.Session(), LeftAlign); ok {
		return align
	}
	return LeftAlign
}

// GetListItemClickedListeners returns a ListItemClickedListener of the ListView.
// If there are no listeners then the empty list is returned
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetListItemClickedListeners(view View, subviewID string) []func(ListView, int) {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if value := view.Get(ListItemClickedEvent); value != nil {
			if result, ok := value.([]func(ListView, int)); ok {
				return result
			}
		}
	}
	return []func(ListView, int){}
}

// GetListItemSelectedListeners returns a ListItemSelectedListener of the ListView.
// If there are no listeners then the empty list is returned
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetListItemSelectedListeners(view View, subviewID string) []func(ListView, int) {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if value := view.Get(ListItemSelectedEvent); value != nil {
			if result, ok := value.([]func(ListView, int)); ok {
				return result
			}
		}
	}
	return []func(ListView, int){}
}

// GetListItemCheckedListeners returns a ListItemCheckedListener of the ListView.
// If there are no listeners then the empty list is returned
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetListItemCheckedListeners(view View, subviewID string) []func(ListView, []int) {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if value := view.Get(ListItemCheckedEvent); value != nil {
			if result, ok := value.([]func(ListView, []int)); ok {
				return result
			}
		}
	}
	return []func(ListView, []int){}
}

// GetListViewCurrent returns the index of the ListView selected item or <0 if there is no a selected item.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetListViewCurrent(view View, subviewID string) int {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if result, ok := intProperty(view, Current, view.Session(), -1); ok {
			return result
		}
	}
	return -1
}

// GetListItemWidth returns the width of a ListView item.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetListItemWidth(view View, subviewID string) SizeUnit {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		result, _ := sizeProperty(view, ItemWidth, view.Session())
		return result
	}
	return AutoSize()
}

// GetListItemHeight returns the height of a ListView item.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetListItemHeight(view View, subviewID string) SizeUnit {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		result, _ := sizeProperty(view, ItemHeight, view.Session())
		return result
	}
	return AutoSize()
}

// GetListViewCheckbox returns the ListView checkbox type: NoneCheckbox (0), SingleCheckbox (1), or MultipleCheckbox (2).
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetListViewCheckbox(view View, subviewID string) int {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		result, _ := enumProperty(view, ItemCheckbox, view.Session(), 0)
		return result
	}
	return 0
}

// GetListViewCheckedItems returns the array of ListView checked items.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetListViewCheckedItems(view View, subviewID string) []int {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if listView, ok := view.(ListView); ok {
			checkedItems := listView.getCheckedItems()
			switch GetListViewCheckbox(view, "") {
			case NoneCheckbox:
				return []int{}

			case SingleCheckbox:
				if len(checkedItems) > 1 {
					return []int{checkedItems[0]}
				}
			}

			return checkedItems
		}
	}
	return []int{}
}

// IsListViewCheckedItem returns true if the ListView item with index is checked, false otherwise.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func IsListViewCheckedItem(view View, subviewID string, index int) bool {
	for _, n := range GetListViewCheckedItems(view, subviewID) {
		if n == index {
			return true
		}
	}
	return false
}

// GetListViewCheckboxVerticalAlign returns the vertical align of the ListView checkbox:
// TopAlign (0), BottomAlign (1), CenterAlign (2)
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetListViewCheckboxVerticalAlign(view View, subviewID string) int {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if align, ok := enumProperty(view, CheckboxVerticalAlign, view.Session(), TopAlign); ok {
			return align
		}
	}
	return TopAlign
}

// GetListViewCheckboxHorizontalAlign returns the horizontal align of the ListView checkbox:
// LeftAlign (0), RightAlign (1), CenterAlign (2)
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetListViewCheckboxHorizontalAlign(view View, subviewID string) int {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if align, ok := enumProperty(view, CheckboxHorizontalAlign, view.Session(), LeftAlign); ok {
			return align
		}
	}
	return LeftAlign
}

// GetListItemVerticalAlign returns the vertical align of the ListView item content:
// TopAlign (0), BottomAlign (1), CenterAlign (2)
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetListItemVerticalAlign(view View, subviewID string) int {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if align, ok := enumProperty(view, ItemVerticalAlign, view.Session(), TopAlign); ok {
			return align
		}
	}
	return TopAlign
}

// ItemHorizontalAlign returns the horizontal align of the ListView item content:
// LeftAlign (0), RightAlign (1), CenterAlign (2), StretchAlign (3)
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetListItemHorizontalAlign(view View, subviewID string) int {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if align, ok := enumProperty(view, ItemHorizontalAlign, view.Session(), LeftAlign); ok {
			return align
		}
	}
	return LeftAlign
}

// GetListItemFrame - returns the location and size of the ListView item in pixels.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetListItemFrame(view View, subviewID string, index int) Frame {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if listView, ok := view.(ListView); ok {
			itemFrames := listView.getItemFrames()
			if index >= 0 && index < len(itemFrames) {
				return itemFrames[index]
			}
		}
	}
	return Frame{Left: 0, Top: 0, Width: 0, Height: 0}
}

// GetListViewAdapter - returns the ListView adapter.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetListViewAdapter(view View, subviewID string) ListAdapter {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if value := view.Get(Items); value != nil {
			if adapter, ok := value.(ListAdapter); ok {
				return adapter
			}
		}
	}
	return nil
}

// ReloadListViewData updates ListView content
// If the second argument (subviewID) is "" then content the first argument (view) is updated.
func ReloadListViewData(view View, subviewID string) {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if listView, ok := view.(ListView); ok {
			listView.ReloadListViewData()
		}
	}
}