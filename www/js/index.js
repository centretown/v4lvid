const blank = "blank";
var leftAction = blank;
function doAction(action) {
    if (leftAction === action) {
        htmx.swap("#slot-left", "", {swapStyle: 'innerHTML'});
        leftAction = blank;
    } else {
        htmx.trigger("#"+action, "click");
        leftAction = action;
    }
}

var hideLeft = true;
function toggleMenu(id) {
    if (hideLeft) {
        hideLeft = false;
        htmx.removeClass("#"+id,"hide")
    } else {
        hideLeft = true;
        htmx.addClass("#"+id,"hide")
    }
}

var hideChat = true;
function toggleChat(id) {
    if (hideChat) {
        hideChat = false;
        htmx.removeClass("#"+id,"hide")
    } else {
        hideChat = true;
        htmx.addClass("#"+id,"hide")
    }
}

function clickSource() {
    if (leftAction !== blank) {
        doAction(leftAction)
    }
}
