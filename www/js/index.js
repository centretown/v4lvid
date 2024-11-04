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

const chatId = "chat";
let drag_offsetX = 0;
let drag_offsetY = 0;
  
function dragstartHandler(ev) {
    // Add the target element's id to the data transfer object
    console.log("dragstart", ev.target.id, ev.clientX, ev.clientY,
        ev.offsetX, ev.offsetY);
    drag_offsetX = ev.offsetX; 
    drag_offsetY = ev.offsetY;
    ev.dataTransfer.setData("text/plain", ev.target.id);
}

function dragendHandler(ev) {
    // Add the target element's id to the data transfer object
    console.log("dragend", ev.target.id, ev.clientX, ev.clientY,
        ev.offsetX, ev.offsetY);
    const element = document.getElementById(chatId);
    console.log("current-left and top", element.style.left, element.style.top)
    console.log("current-right and bottom", element.style.right, element.style.bottom)
    element.style.left = (ev.clientX-drag_offsetX)+'px';
    element.style.top = (ev.clientY-drag_offsetY)+'px';
}

function dropHandler(ev) {
    // Add the target element's id to the data transfer object
    console.log("drop", ev.clientX, ev.clientY,
        ev.offsetX, ev.offsetY);
}

window.addEventListener("DOMContentLoaded", () => {
    // Get the element by id
    const element = document.getElementById(chatId);
    // Add the ondragstart event listener
    element.addEventListener("dragstart", dragstartHandler);
    element.addEventListener("dragend", dragendHandler);
});