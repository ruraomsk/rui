var sessionID = "0"
var socket
var socketUrl

var images = new Map();
var windowFocus = true

function sendMessage(message) {
	if (socket) {
		socket.send(message)
	}
}

window.onload = function() {
	socketUrl = document.location.protocol == "https:" ? "wss://" : "ws://" 
	socketUrl += document.location.hostname
	var port = document.location.port
	if (port) {
		socketUrl += ":" + port
	}
	socketUrl += window.location.pathname + "ws"

	socket = new WebSocket(socketUrl);
	socket.onopen = socketOpen;
	socket.onclose = socketClose;
	socket.onerror = socketError;
	socket.onmessage = function(event) {
		window.execScript ? window.execScript(event.data) : window.eval(event.data);
	};
};

function socketOpen() {

	const touch_screen = (('ontouchstart' in document.documentElement) || (navigator.maxTouchPoints > 0) || (navigator.msMaxTouchPoints > 0)) ? "1" : "0";
	var message = "startSession{touch=" + touch_screen 
	
	const style = window.getComputedStyle(document.body);
	if (style) {
		var direction = style.getPropertyValue('direction');
		if (direction) {
			message += ",direction=" + direction
		}
	}

	const lang = window.navigator.languages;
	if (lang) {
		message += ",languages=\"" + lang + "\"";
	}

	const darkThemeMq = window.matchMedia("(prefers-color-scheme: dark)");
	if (darkThemeMq.matches) {
		message += ",dark=1";
	} 

	const pixelRatio = window.devicePixelRatio;
	if (pixelRatio) {
		message += ",pixel-ratio=" + pixelRatio;
	}

	sendMessage( message + "}" );
}

function socketReopen() {
	sendMessage( "reconnect{session=" + sessionID + "}" );
}

function socketReconnect() {
	if (!socket) {
		socket = new WebSocket(socketUrl);
		socket.onopen = socketReopen;
		socket.onclose = socketClose;
		socket.onerror = socketError;
		socket.onmessage = function(event) {
			window.execScript ? window.execScript(event.data) : window.eval(event.data);
		};
	}
}

function socketClose(event) {
	console.log("socket closed")
	socket = null;
	if (!event.wasClean && windowFocus) {
		window.setTimeout(socketReconnect, 10000);
	}
	/*
	if (event.wasClean) {
		alert('Connection was clean closed');
	} else {
		alert('Connection was lost');
	}
	alert('Code: ' + event.code + ' reason: ' + event.reason);
	*/
}

function socketError(error) {
	console.log(error);
}

window.onresize = function() {
	scanElementsSize();
}

window.onbeforeunload = function(event) {
	sendMessage( "session-close{session=" + sessionID +"}" );
}

window.onblur = function(event) {
	windowFocus = false
	sendMessage( "session-pause{session=" + sessionID +"}" );
}

window.onfocus = function(event) {
	windowFocus = true
	if (!socket) {
		socketReconnect()
	} else {
		sendMessage( "session-resume{session=" + sessionID +"}" );
	}
}

function getIntAttribute(element, tag) {
	let value = element.getAttribute(tag);
	if (value) {
		return value;
	}
	return 0;
}

function scanElementsSize() {
	var views = document.getElementsByClassName("ruiView");
	if (views) {
		var message = "resize{session=" + sessionID + ",views=["
		var count = 0
		for (var i = 0; i < views.length; i++) {
			let element = views[i];
			let noresize = element.getAttribute("data-noresize");
			if (!noresize) {
				let rect = element.getBoundingClientRect();
				let top = getIntAttribute(element, "data-top");
				let left = getIntAttribute(element, "data-left");
				let width = getIntAttribute(element, "data-width");
				let height = getIntAttribute(element, "data-height");
				if (rect.width > 0 && rect.height > 0 && 
					(width != rect.width || height != rect.height || left != rect.left || top != rect.top)) {
					element.setAttribute("data-top", rect.top);
					element.setAttribute("data-left", rect.left);
					element.setAttribute("data-width", rect.width);
					element.setAttribute("data-height", rect.height);
					if (count > 0) {
						message += ",";
					}
					message += "view{id=" + element.id + ",x=" + rect.left + ",y=" + rect.top + ",width=" + rect.width + ",height=" + rect.height + 
						",scroll-x=" + element.scrollLeft + ",scroll-y=" + element.scrollTop + ",scroll-width=" + element.scrollWidth + ",scroll-height=" + element.scrollHeight + "}";
					count += 1;
				}
			}
		}

		if (count > 0) {
			sendMessage(message + "]}");
		}
	}
}

function scrollEvent(element, event) {
	sendMessage("scroll{session=" + sessionID + ",id=" + element.id + ",x=" + element.scrollLeft + 
		",y=" + element.scrollTop + ",width=" + element.scrollWidth + ",height=" + element.scrollHeight + "}");
}

function updateCSSRule(selector, ruleText) {
	var styleSheet = document.styleSheets[0];
	var rules = styleSheet.cssRules ? styleSheet.cssRules : styleSheet.rules
	selector = "." + selector
	for (var i = 0; i < rules.length; i++) {
		var rule = rules[i]
		if (!rule.selectorText) {
			continue;
		}
		if (rule.selectorText == selector) {
			if (styleSheet.deleteRule) {
				styleSheet.deleteRule(i)
			} else if (styleSheet.removeRule) {
				styleSheet.removeRule(i)
			}
			break;
		}
	}
	if (styleSheet.insertRule) {
		styleSheet.insertRule(selector + " { " + ruleText + "}")
	} else if (styleSheet.addRule) {
		styleSheet.addRule(selector, ruleText, rules.length)
	}
	scanElementsSize();
}

function updateCSSStyle(elementId, style) {
	var element = document.getElementById(elementId);
	if (element) {
		element.style = style;
		scanElementsSize();
	}
}

function updateCSSProperty(elementId, property, value) {
	var element = document.getElementById(elementId);
	if (element) {
		element.style[property] = value;
		scanElementsSize();
	}
}

function updateProperty(elementId, property, value) {
	var element = document.getElementById(elementId);
	if (element) {
		element.setAttribute(property, value);
		scanElementsSize();
	}
}

function removeProperty(elementId, property, value) {
	var element = document.getElementById(elementId);
	if (element && element.hasAttribute(property)) {
		element.removeAttribute(property);
		scanElementsSize();
	}
}

function updateInnerHTML(elementId, content) {
	var element = document.getElementById(elementId);
	if (element) {
		element.innerHTML = content;
		scanElementsSize();
	}
}

function appendToInnerHTML(elementId, content) {
	var element = document.getElementById(elementId);
	if (element) {
		element.innerHTML += content;
		scanElementsSize();
	}
}

function setDisabled(elementId, disabled) {
	var element = document.getElementById(elementId);
	if (element) {
		if ('disabled' in element) {
			element.disabled = disabled
		} else {
			element.setAttribute("data-disabled", disabled ? "1" : "0");
		}
		scanElementsSize();
	}
}

function focusEvent(element, event) {
	event.stopPropagation();
	sendMessage("focus-event{session=" + sessionID + ",id=" + element.id + "}");
}

function blurEvent(element, event) {
	event.stopPropagation();
	sendMessage("lost-focus-event{session=" + sessionID + ",id=" + element.id + "}");
}

function enterOrSpaceKeyClickEvent(event) {
	if (event.key) {
		return (event.key == " " || event.key == "Enter");
	} else if (event.keyCode) {
		return (event.keyCode == 32 || event.keyCode == 13);
	}
	return false;
}

function activateTab(layoutId, tabNumber) {
	var element = document.getElementById(layoutId);
	if (element) {
		var currentTabId = element.getAttribute("data-current");
		var newTabId = layoutId + '-' + tabNumber;
		if (currentTabId != newTabId) {
			function setTab(tabId, styleProperty, display) {
				var tab = document.getElementById(tabId);
				if (tab) {	
					tab.className = element.getAttribute(styleProperty);
					var page = document.getElementById(tab.getAttribute("data-view"));
					if (page) {
						page.style.display = display;
					}
				}
			}
			setTab(currentTabId, "data-inactiveTabStyle", "none")
			setTab(newTabId, "data-activeTabStyle", "");
			element.setAttribute("data-current", newTabId);
			scanElementsSize()
		}
	}
}

function tabClickEvent(layoutId, tabNumber, event) {
	event.stopPropagation();
	event.preventDefault();
	activateTab(layoutId, tabNumber)
	sendMessage("tabClick{session=" + sessionID + ",id=" + layoutId + ",number=" + tabNumber + "}");
}

function tabKeyClickEvent(layoutId, tabNumber, event) {
	if (enterOrSpaceKeyClickEvent(event)) {
		tabClickEvent(layoutId, tabNumber, event)
	}
}

function keyEvent(element, event, tag) {
	event.stopPropagation();

	var message = tag + "{session=" + sessionID + ",id=" + element.id;
	if (event.timeStamp) {
		message += ",timeStamp=" + event.timeStamp;
	}
	if (event.key) {
		message += ",key=\"" + event.key + "\"";
	}
	if (event.code) {
		message += ",code=\"" + event.code + "\"";
	}
	if (event.repeat) {
		message += ",repeat=1";
	}
	if (event.ctrlKey) {
		message += ",ctrlKey=1";
	}
	if (event.shiftKey) {
		message += ",shiftKey=1";
	}
	if (event.altKey) {
		message += ",altKey=1";
	}
	if (event.metaKey) {
		message += ",metaKey=1";
	}

	message += "}"
	sendMessage(message);
}

function keyDownEvent(element, event) {
	keyEvent(element, event, "key-down-event")
}

function keyUpEvent(element, event) {
	keyEvent(element, event, "key-up-event")
}

function mouseEventData(element, event) {
	var message = ""

	if (event.timeStamp) {
		message += ",timeStamp=" + event.timeStamp;
	}
	if (event.button) {
		message += ",button=" + event.button;
	}
	if (event.buttons) {
		message += ",buttons=" + event.buttons;
	}
	if (event.clientX) {
		var x = event.clientX;
		var el = element;
		if (el.parentElement) {
			x += el.parentElement.scrollLeft;
		}
		while (el) {
			x -= el.offsetLeft
			el = el.parentElement
		}

		message += ",x=" + x + ",clientX=" + event.clientX;
	}
	if (event.clientY) {
		var y = event.clientY;
		var el = element;
		if (el.parentElement) {
			y += el.parentElement.scrollTop;
		}
		while (el) {
			y -= el.offsetTop
			el = el.parentElement
		}

		message += ",y=" + y + ",clientY=" + event.clientY;
	}
	if (event.screenX) {
		message += ",screenX=" + event.screenX;
	}
	if (event.screenY) {
		message += ",screenY=" + event.screenY;
	}
	if (event.ctrlKey) {
		message += ",ctrlKey=1";
	}
	if (event.shiftKey) {
		message += ",shiftKey=1";
	}
	if (event.altKey) {
		message += ",altKey=1";
	}
	if (event.metaKey) {
		message += ",metaKey=1";
	}
	return message
}

function mouseEvent(element, event, tag) {
	event.stopPropagation();
	//event.preventDefault()

	var message = tag + "{session=" + sessionID + ",id=" + element.id + mouseEventData(element, event) + "}";
	sendMessage(message);
}

function mouseDownEvent(element, event) {
	mouseEvent(element, event, "mouse-down")
}

function mouseUpEvent(element, event) {
	mouseEvent(element, event, "mouse-up")
}

function mouseMoveEvent(element, event) {
	mouseEvent(element, event, "mouse-move")
}

function mouseOverEvent(element, event) {
	mouseEvent(element, event, "mouse-over")
}

function mouseOutEvent(element, event) {
	mouseEvent(element, event, "mouse-out")
}

function clickEvent(element, event) {
	mouseEvent(element, event, "click-event")
	event.preventDefault();
}

function doubleClickEvent(element, event) {
	mouseEvent(element, event, "double-click-event")
	event.preventDefault();
}

function contextMenuEvent(element, event) {
	mouseEvent(element, event, "context-menu-event")
	event.preventDefault();
}

function pointerEvent(element, event, tag) {
	event.stopPropagation();

	var message = tag + "{session=" + sessionID + ",id=" + element.id + mouseEventData(element, event);

	if (event.pointerId) {
		message += ",pointerId=" + event.pointerId;
	}
	if (event.width) {
		message += ",width=" + event.width;
	}
	if (event.height) {
		message += ",height=" + event.height;
	}
	if (event.pressure) {
		message += ",pressure=" + event.pressure;
	}
	if (event.tangentialPressure) {
		message += ",tangentialPressure=" + event.tangentialPressure;
	}
	if (event.tiltX) {
		message += ",tiltX=" + event.tiltX;
	}
	if (event.tiltY) {
		message += ",tiltY=" + event.tiltY;
	}
	if (event.twist) {
		message += ",twist=" + event.twist;
	}
	if (event.pointerType) {
		message += ",pointerType=" + event.pointerType;
	}
	if (event.isPrimary) {
		message += ",isPrimary=1";
	}

	message += "}";
	sendMessage(message);
}

function pointerDownEvent(element, event) {
	pointerEvent(element, event, "pointer-down")
}

function pointerUpEvent(element, event) {
	pointerEvent(element, event, "pointer-up")
}

function pointerMoveEvent(element, event) {
	pointerEvent(element, event, "pointer-move")
}

function pointerCancelEvent(element, event) {
	pointerEvent(element, event, "pointer-cancel")
}

function pointerOverEvent(element, event) {
	pointerEvent(element, event, "pointer-over")
}

function pointerOutEvent(element, event) {
	pointerEvent(element, event, "pointer-out")
}

function touchEvent(element, event, tag) {
	event.stopPropagation();

	var message = tag + "{session=" + sessionID + ",id=" + element.id;
	if (event.timeStamp) {
		message += ",timeStamp=" + event.timeStamp;
	}
	if (event.touches && event.touches.length > 0) {
		message += ",touches=["
		for (var i = 0; i < event.touches.length; i++) {
			var touch = event.touches.item(i)
			if (touch) {
				if (i > 0) {
					message += ","	
				}
				message += "touch{identifier=" + touch.identifier;

				var x = touch.clientX;
				var y = touch.clientY;
				var el = element;
				if (el.parentElement) {
					x += el.parentElement.scrollLeft;
					y += el.parentElement.scrollTop;
				}
				while (el) {
					x -= el.offsetLeft
					y -= el.offsetTop
					el = el.parentElement
				}
		
				message += ",x=" + x + ",y=" + y + ",clientX=" + touch.clientX + ",clientY=" + touch.clientY +
					",screenX=" + touch.screenX + ",screenY=" + touch.screenY + ",radiusX=" + touch.radiusX +
					",radiusY=" + touch.radiusY + ",rotationAngle=" + touch.rotationAngle + ",force=" + touch.force + "}"
			}
		}
		message += "]"
	}
	if (event.ctrlKey) {
		message += ",ctrlKey=1";
	}
	if (event.shiftKey) {
		message += ",shiftKey=1";
	}
	if (event.altKey) {
		message += ",altKey=1";
	}
	if (event.metaKey) {
		message += ",metaKey=1";
	}

	message += "}";
	sendMessage(message);
}

function touchStartEvent(element, event) {
	touchEvent(element, event, "touch-start")
}

function touchEndEvent(element, event) {
	touchEvent(element, event, "touch-end")
}

function touchMoveEvent(element, event) {
	touchEvent(element, event, "touch-move")
}

function touchCancelEvent(element, event) {
	touchEvent(element, event, "touch-cancel")
}

function dropDownListEvent(element, event) {
	event.stopPropagation();
	var message = "itemSelected{session=" + sessionID + ",id=" + element.id + ",number=" + element.selectedIndex.toString() + "}"
	sendMessage(message);
}

function selectDropDownListItem(elementId, number) {
	var element = document.getElementById(elementId);
	if (element) {
		element.selectedIndex = number;
		scanElementsSize();
	}
}

function listItemClickEvent(element, event) {
	event.stopPropagation();
	var selected = false;
	if (element.classList) {
		selected = (element.classList.contains("ruiListItemFocused") || element.classList.contains("ruiListItemSelected"));
	} else {
		selected = element.className.indexOf("ruiListItemFocused") >= 0 || element.className.indexOf("ruiListItemSelected") >= 0;
	}

	var list = element.parentNode
	if (list) {
		if (!selected) {
			selectListItem(list, element, true)
		}

		var message = "itemClick{session=" + sessionID + ",id=" + list.id + "}"
		sendMessage(message);
	}
}

function getListItemNumber(itemId) {
	var pos = itemId.indexOf("-")
	if (pos >= 0) {
		return parseInt(itemId.substring(pos+1))
	}
}

function selectListItem(element, item, needSendMessage) {
	var currentId = element.getAttribute("data-current");
	var message;
	var focusStyle = element.getAttribute("data-focusitemstyle");
	var blurStyle = element.getAttribute("data-bluritemstyle");

	if (!focusStyle) {
		focusStyle = "ruiListItemFocused"
	}
	if (!blurStyle) {
		blurStyle = "ruiListItemSelected"
	}

	if (currentId) {
		var current = document.getElementById(currentId);
		if (current) {
			if (current.classList) {
				current.classList.remove(focusStyle, blurStyle);
			} else { // IE < 10
				current.className = "ruiListItem";
			}
			if (sendMessage) {
				message = "itemUnselected{session=" + sessionID + ",id=" + element.id + "}";
			}
		}
	}

	if (item) {
		if (element === document.activeElement) {
			if (item.classList) {
				item.classList.add(focusStyle);
			} else { // IE < 10
				item.className = "ruiListItem " + focusStyle
			}
		} else {
			if (item.classList) {
				item.classList.add(blurStyle);
			} else { // IE < 10
				item.className = "ruiListItem " + blurStyle
			}
		}

		element.setAttribute("data-current", item.id);
		if (sendMessage) {
			var number = getListItemNumber(item.id)
			if (number != undefined) {
				message = "itemSelected{session=" + sessionID + ",id=" + element.id + ",number=" + number + "}";
			}
		}

		var left = item.offsetLeft - element.offsetLeft;
		if (left < element.scrollLeft) {
			element.scrollLeft = left;
		}

		var top = item.offsetTop - element.offsetTop;
		if (top < element.scrollTop) {
			element.scrollTop = top;
		}
		
		var right = left + item.offsetWidth;
		if (right > element.scrollLeft + element.clientWidth) {
			element.scrollLeft = right - element.clientWidth;
		}

		var bottom = top + item.offsetHeight
		if (bottom > element.scrollTop + element.clientHeight) {
			element.scrollTop = bottom - element.clientHeight;
		}
	}

	if (needSendMessage && message != undefined) {
		sendMessage(message);
	}
	scanElementsSize();
}

function findRightListItem(list, x, y) {
	var result;
	var count = list.childNodes.length;
	for (var i = 0; i < count; i++) {
		var item = list.childNodes[i];
		if (item.offsetLeft >= x) {
			if (result) {
				var result_dy = Math.abs(result.offsetTop - y);
				var item_dy = Math.abs(item.offsetTop - y);
				if (item_dy < result_dy || (item_dy == result_dy && (item.offsetLeft - x) < (result.offsetLeft - x))) {
					result = item;	
				}
			} else {
				result = item;
			}
		}
	}
	return result
}

function findLeftListItem(list, x, y) {
	var result;
	var count = list.childNodes.length;
	for (var i = 0; i < count; i++) {
		var item = list.childNodes[i];
		if (item.offsetLeft < x) {
			if (result) {
				var result_dy = Math.abs(result.offsetTop - y);
				var item_dy = Math.abs(item.offsetTop - y);
				if (item_dy < result_dy || (item_dy == result_dy && (x - item.offsetLeft) < (x - result.offsetLeft))) {
					result = item;	
				}
			} else {
				result = item;
			}
		}
	}
	return result
}

function findTopListItem(list, x, y) {
	var result;
	var count = list.childNodes.length;
	for (var i = 0; i < count; i++) {
		var item = list.childNodes[i];
		if (item.offsetTop < y) {
			if (result) {
				var result_dx = Math.abs(result.offsetLeft - x);
				var item_dx = Math.abs(item.offsetLeft - x);
				if (item_dx < result_dx || (item_dx == result_dx && (y - item.offsetTop) < (y - result.offsetTop))) {
					result = item;	
				}
			} else {
				result = item;
			}
		}
	}
	return result
}

function findBottomListItem(list, x, y) {
	var result;
	var count = list.childNodes.length;
	for (var i = 0; i < count; i++) {
		var item = list.childNodes[i];
		if (item.offsetTop >= y) {
			if (result) {
				var result_dx = Math.abs(result.offsetLeft - x);
				var item_dx = Math.abs(item.offsetLeft - x);
				if (item_dx < result_dx || (item_dx == result_dx && (item.offsetTop - y) < (result.offsetTop - y))) {
					result = item;	
				}
			} else {
				result = item;
			}
		}
	}
	return result
}

function listViewKeyDownEvent(element, event) {
	var key;
	if (event.key) {
		key = event.key;
	} else if (event.keyCode) {
		switch (event.keyCode) {
			case 13: key = "Enter"; break;
			case 32: key = " "; break;
			case 33: key = "PageUp"; break;
			case 34: key = "PageDown"; break;
			case 35: key = "End"; break;
			case 36: key = "Home"; break;
			case 37: key = "ArrowLeft"; break;
			case 38: key = "ArrowUp"; break;
			case 39: key = "ArrowRight"; break;
			case 40: key = "ArrowDown"; break;
		}
	}
	if (key) {
		var currentId = element.getAttribute("data-current");
		var current
		if (currentId) {
			current = document.getElementById(currentId);
			//number = getListItemNumber(currentId);
		}
		if (current) {
			var item
			switch (key) {
			case " ": 
			case "Enter":
				var message = "itemClick{session=" + sessionID + ",id=" + element.id + "}";
				sendMessage(message);
				break;

			case "ArrowLeft":
				item = findLeftListItem(element, current.offsetLeft, current.offsetTop);
				break;
			
			case "ArrowRight":
				item = findRightListItem(element, current.offsetLeft + current.offsetWidth, current.offsetTop);
				break;
	
			case "ArrowDown":
				item = findBottomListItem(element, current.offsetLeft, current.offsetTop + current.offsetHeight);
				break;

			case "ArrowUp":
				item = findTopListItem(element, current.offsetLeft, current.offsetTop);
				break;

			case "Home":
				item = element.childNodes[0];
				break;

			case "End":
				item = element.childNodes[element.childNodes.length - 1];
				break;

			case "PageUp":
				// TODO
				break;

			case "PageDown":
				// TODO
				break;

			default:
				return;
			}
			if (item && item !== current) {
				selectListItem(element, item, true);
			}
		}
	}

	event.stopPropagation();
	event.preventDefault();
}

function listViewFocusEvent(element, event) {
	var currentId = element.getAttribute("data-current");
	if (currentId) {
		var current = document.getElementById(currentId);
		if (current) {
			var focusStyle = element.getAttribute("data-focusitemstyle");
			var blurStyle = element.getAttribute("data-bluritemstyle");
			if (!focusStyle) {
				focusStyle = "ruiListItemFocused"
			}
			if (!blurStyle) {
				blurStyle = "ruiListItemSelected"
			}
			
			if (current.classList) {
				current.classList.remove(blurStyle);
				current.classList.add(focusStyle);
			} else { // IE < 10
				current.className = "ruiListItem " + focusStyle;
			}
		}
	}
}

function listViewBlurEvent(element, event) {
	var currentId = element.getAttribute("data-current");
	if (currentId) {
		var current = document.getElementById(currentId);
		if (current) {
			var focusStyle = element.getAttribute("data-focusitemstyle");
			var blurStyle = element.getAttribute("data-bluritemstyle");
			if (!focusStyle) {
				focusStyle = "ruiListItemFocused"
			}
			if (!blurStyle) {
				blurStyle = "ruiListItemSelected"
			}

			if (current.classList) {
				current.classList.remove(focusStyle);
				current.classList.add(blurStyle);
			} else { // IE < 10
				current.className = "ruiListItem " + blurStyle;
			}
		}
	}
}

function selectRadioButton(radioButtonId) {
	var element = document.getElementById(radioButtonId);
	if (element) {
		var list = element.parentNode
		if (list) {
			var current = list.getAttribute("data-current");
			if (current) {
				if (current === radioButtonId) {
					return
				}

				var mark = document.getElementById(current + "mark");
				if (mark) {
					//mark.hidden = true
					mark.style.visibility = "hidden"
				}
			}

			var mark = document.getElementById(radioButtonId + "mark");
			if (mark) {
				//mark.hidden = false
				mark.style.visibility = "visible"
			}
			list.setAttribute("data-current", radioButtonId);
			var message = "radioButtonSelected{session=" + sessionID + ",id=" + list.id + ",radioButton=" + radioButtonId + "}"
			sendMessage(message);
			scanElementsSize();
		}
	}
}

function unselectRadioButtons(radioButtonsId) {
	var list = document.getElementById(radioButtonsId);
	if (list) {
		var current = list.getAttribute("data-current");
		if (current) {
			var mark = document.getElementById(current + "mark");
			if (mark) {
				mark.style.visibility = "hidden"
			}

			list.removeAttribute("data-current");
		}

		var message = "radioButtonUnselected{session=" + sessionID + ",id=" + list.id + "}"
		sendMessage(message);
		scanElementsSize();
	}
}

function radioButtonClickEvent(element, event) {
	event.stopPropagation();
	event.preventDefault();
	selectRadioButton(element.id)
}

function radioButtonKeyClickEvent(element, event) {
	if (enterOrSpaceKeyClickEvent(event)) {
		radioButtonClickEvent(element, event);
	}
}

function editViewInputEvent(element) {
	var text = element.value
	text = text.replace(/\\/g, "\\\\")
	text = text.replace(/\"/g, "\\\"")
	var message = "textChanged{session=" + sessionID + ",id=" + element.id + ",text=\"" + text + "\"}"
	sendMessage(message);
}

function setInputValue(elementId, text) {
	var element = document.getElementById(elementId);
	if (element) {
		element.value = text;
		scanElementsSize();
	}
}

function fileSelectedEvent(element) {
	var files = element.files;
	if (files) {
		var message = "fileSelected{session=" + sessionID + ",id=" + element.id + ",files=[";
		for(var i = 0; i < files.length; i++) {
			if (i > 0) {
				message += ",";
			}
			message += "_{name=\"" + files[i].name + 
				"\",last-modified=" + files[i].lastModified +
				",size=" + files[i].size +
				",mime-type=\"" + files[i].type + "\"}";
		}
		sendMessage(message + "]}");
	}
}

function loadSelectedFile(elementId, index) {
	var element = document.getElementById(elementId);
	if (element) {
		var files = element.files;
		if (files && index >= 0 && index < files.length) {
			const reader = new FileReader();
         	reader.onload = function() { 
				sendMessage("fileLoaded{session=" + sessionID + ",id=" + element.id + 
					",index=" + index + 
					",name=\"" + files[index].name + 
					"\",last-modified=" + files[index].lastModified +
					",size=" + files[index].size +
					",mime-type=\"" + files[index].type + 
					"\",data=`" + reader.result + "`}");
			}
         	reader.onerror = function(error) {
				sendMessage("fileLoadingError{session=" + sessionID + ",id=" + element.id + ",index=" + index + ",error=`" + error + "`}");
			}
			reader.readAsDataURL(files[index]);
		} else {
			sendMessage("fileLoadingError{session=" + sessionID + ",id=" + element.id + ",index=" + index + ",error=`File not found`}");
		}
	} else {
		sendMessage("fileLoadingError{session=" + sessionID + ",id=" + element.id + ",index=" + index + ",error=`Invalid FilePicker id`}");
	}
}

function startResize(element, mx, my, event) {
	var view = element.parentNode;
	if (!view) {
		return;
	}

	var startX = event.clientX;
	var startY = event.clientY;
	var startWidth = view.offsetWidth
	var startHeight = view.offsetHeight

	document.addEventListener("mousemove", moveHandler, true);
	document.addEventListener("mouseup", upHandler, true);
	
	event.stopPropagation();
	event.preventDefault();

	function moveHandler(e) {
		if (mx != 0) {
			var width = startWidth + (e.clientX - startX) * mx;
			if (width <= 0) {
				width = 1;
			}
			view.style.width = width + "px";
			sendMessage("widthChanged{session=" + sessionID + ",id=" + view.id + ",width=" + view.style.width + "}");
		}
		
		if (my != 0) {
			var height = startHeight + (e.clientY - startY) * my;
			if (height <= 0) {
				height = 1;
			}
			view.style.height = height + "px";
			sendMessage("heightChanged{session=" + sessionID + ",id=" + view.id + ",height=" + view.style.height + "}");
		}

		event.stopPropagation();
		event.preventDefault();
		scanElementsSize();
	}

	function upHandler (e) {
		document.removeEventListener("mouseup", upHandler, true);
		document.removeEventListener("mousemove", moveHandler, true);
		e.stopPropagation();
	}
}

function transitionStartEvent(element, event) {
	var message = "transition-start-event{session=" + sessionID + ",id=" + element.id; 
	if (event.propertyName) {
		message += ",property=" + event.propertyName
	}
	sendMessage(message + "}");
	event.stopPropagation();
}

function transitionRunEvent(element, event) {
	var message = "transition-run-event{session=" + sessionID + ",id=" + element.id; 
	if (event.propertyName) {
		message += ",property=" + event.propertyName
	}
	sendMessage(message + "}");
	event.stopPropagation();
}

function transitionEndEvent(element, event) {
	var message = "transition-end-event{session=" + sessionID + ",id=" + element.id; 
	if (event.propertyName) {
		message += ",property=" + event.propertyName
	}
	sendMessage(message + "}");
	event.stopPropagation();
}

function transitionCancelEvent(element, event) {
	var message = "transition-cancel-event{session=" + sessionID + ",id=" + element.id; 
	if (event.propertyName) {
		message += ",property=" + event.propertyName
	}
	sendMessage(message + "}");
	event.stopPropagation();
}

function animationStartEvent(element, event) {
	var message = "animation-start-event{session=" + sessionID + ",id=" + element.id; 
	if (event.animationName) {
		message += ",name=" + event.animationName
	}
	sendMessage(message + "}");
	event.stopPropagation();
}

function animationEndEvent(element, event) {
	var message = "animation-end-event{session=" + sessionID + ",id=" + element.id; 
	if (event.animationName) {
		message += ",name=" + event.animationName
	}
	sendMessage(message + "}");
	event.stopPropagation();
}

function animationCancelEvent(element, event) {
	var message = "animation-cancel-event{session=" + sessionID + ",id=" + element.id; 
	if (event.animationName) {
		message += ",name=" + event.animationName
	}
	sendMessage(message + "}");
	event.stopPropagation();
}

function animationIterationEvent(element, event) {
	var message = "animation-iteration-event{session=" + sessionID + ",id=" + element.id; 
	if (event.animationName) {
		message += ",name=" + event.animationName
	}
	sendMessage(message + "}");
	event.stopPropagation();
}

function stackTransitionEndEvent(stackId, propertyName, event) {
	sendMessage("transition-end-event{session=" + sessionID + ",id=" + stackId + ",property=" + propertyName + "}");
	event.stopPropagation();
}

function loadImage(url) {
	var img = images.get(url);
	if (img != undefined) {
		return
	}
	
	img = new Image(); 
	img.addEventListener("load", function() {
		images.set(url, img)
		var message = "imageLoaded{session=" + sessionID + ",url=\"" + url + "\""; 
		if (img.naturalWidth) {
			message += ",width=" + img.naturalWidth
		}
		if (img.naturalHeight) {
			message += ",height=" + img.naturalHeight
		}
		sendMessage(message + "}")
	}, false);

	img.addEventListener("error", function(event) {
		var message = "imageError{session=" + sessionID + ",url=\"" + url + "\""; 
		if (event && event.message) {
			var text = event.message.replace(new RegExp("\"", 'g'), "\\\"")
			message += ",message=\"" + text + "\""; 
		}
		sendMessage(message + "}")
	}, false);

	img.src = url;
}

function clickOutsidePopup(e) {
	sendMessage("clickOutsidePopup{session=" + sessionID + "}")
	e.stopPropagation();
}

function clickClosePopup(element, e) {
	var popupId = element.getAttribute("data-popupId");
	sendMessage("clickClosePopup{session=" + sessionID + ",id=" + popupId + "}")
	e.stopPropagation();
}

function scrollTo(elementId, x, y) {
	var element = document.getElementById(elementId);
	if (element) {
		element.scrollTo(x, y);
	}
}

function scrollToStart(elementId) {
	var element = document.getElementById(elementId);
	if (element) {
		element.scrollTo(0, 0);
	}
}

function scrollToEnd(elementId) {
	var element = document.getElementById(elementId);
	if (element) {
		element.scrollTo(0, element.scrollHeight - element.offsetHeight);
	}
}

function focus(elementId) {
	var element = document.getElementById(elementId);
	if (element) {
		element.focus();
	}
}

function playerEvent(element, tag) {
	//event.stopPropagation();
	sendMessage(tag + "{session=" + sessionID + ",id=" + element.id + "}");
}

function playerTimeUpdatedEvent(element) {
	var message = "time-update-event{session=" + sessionID + ",id=" + element.id + ",value=";
	if (element.currentTime) {
		message += element.currentTime;
	} else {
		message += "0";
	}
	sendMessage(message + "}");
}

function playerDurationChangedEvent(element) {
	var message = "duration-changed-event{session=" + sessionID + ",id=" + element.id + ",value=";
	if (element.duration) {
		message += element.duration;
	} else {
		message += "0";
	}
	sendMessage(message + "}");
}

function playerVolumeChangedEvent(element) {
	var message = "volume-changed-event{session=" + sessionID + ",id=" + element.id + ",value=";
	if (element.volume && !element.muted) {
		message += element.volume;
	} else {
		message += "0";
	}
	sendMessage(message + "}");
}

function playerRateChangedEvent(element) {
	var message = "rate-changed-event{session=" + sessionID + ",id=" + element.id + ",value=";
	if (element.playbackRate) {
		message += element.playbackRate;
	} else {
		message += "0";
	}
	sendMessage(message + "}");
}

function playerErrorEvent(element) {
	var message = "player-error-event{session=" + sessionID + ",id=" + element.id;
	if (element.error) {
		if (element.error.code) {
			message += ",code=" + element.error.code;
		}
		if (element.error.message) {
			message += ",message=`" + element.error.message + "`";
		}
	}
	sendMessage(message + "}");
}

function setMediaMuted(elementId, value) {
	var element = document.getElementById(elementId);
	if (element) {
		element.muted = value
	}
}

function mediaPlay(elementId) {
	var element = document.getElementById(elementId);
	if (element && element.play) {
		element.play()
	}
}

function mediaPause(elementId) {
	var element = document.getElementById(elementId);
	if (element && element.pause) {
		element.pause()
	}
}

function mediaSetSetCurrentTime(elementId, time) {
	var element = document.getElementById(elementId);
	if (element) {
		element.currentTime = time
	}
}

function mediaSetPlaybackRate(elementId, time) {
	var element = document.getElementById(elementId);
	if (element) {
		element.playbackRate = time
	}
}

function mediaSetVolume(elementId, volume) {
	var element = document.getElementById(elementId);
	if (element) {
		element.volume = volume
	}
}