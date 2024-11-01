const blank = "blank";
var leftAction = blank;
var hideLeft = true;
function doAction(action) {
    if (leftAction === action) {
        htmx.swap("#slot-left", "", {swapStyle: 'innerHTML'});
        leftAction = blank;
    } else {
        htmx.trigger("#"+action, "click");
        leftAction = action;
    }
}
function toggleMenu(id) {
    if (hideLeft) {
        hideLeft = false;
        htmx.removeClass("#"+id,"hide")
    } else {
        hideLeft = true;
        htmx.addClass("#"+id,"hide")
    }
}

function clickSource() {
    if (leftAction !== blank) {
        doAction(leftAction)
    }
}
