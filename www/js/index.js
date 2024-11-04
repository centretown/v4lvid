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

function startTime() {
    const today = new Date();
    let h = today.getHours();
    let m = today.getMinutes();
    m = (m < 10) ? "0" + m : m;
    document.getElementById('clock').innerHTML =  h + ":" + m;
    setTimeout(startTime, 1000*60);
}

// const chatId = "chat";
let drag_data = {};
let chat_data = {};

function dragstartHandler(ev) {
    drag_data.offsetX = ev.offsetX; 
    drag_data.offsetY = ev.offsetY;
}

function dragendHandler(ev) {
    const target = document.getElementById(ev.target.id);
    chat_data.X = ev.clientX - drag_data.offsetX;
    chat_data.Y = ev.clientY - drag_data.offsetY;
    target.style.left = chat_data.X +'px';
    target.style.top = chat_data.Y +'px';
    setdraggable(ev.target.id, false);
}

function addDragHandlers(id) {
    const target = document.getElementById(id);
    target.addEventListener("dragstart", dragstartHandler);
    target.addEventListener("dragend", dragendHandler);
}

function setdraggable(id, dragabble) {
    document.getElementById(id).setAttribute('draggable', dragabble);
}

window.addEventListener("DOMContentLoaded", () => {
    addDragHandlers("chat");
    addDragHandlers("slot-left");
});
