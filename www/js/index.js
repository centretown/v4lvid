const blank = "blank";
var currentAction = blank;
var hideActions = true;
function doAction(action) {
    if (currentAction === action) {
        htmx.swap("#slot", "", {swapStyle: 'innerHTML'});
        currentAction = blank;
    } else {
        htmx.trigger("#"+action, "click");
        currentAction = action;
    }
}
function toggleMenu(id) {
    if (hideActions) {
        hideActions = false;
        htmx.removeClass("#"+id,"hide")
    } else {
        hideActions = true;
        htmx.addClass("#"+id,"hide")
    }
}

function clickSource() {
    if (currentAction !== blank) {
        doAction(currentAction)
    }
}
