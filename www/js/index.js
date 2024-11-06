const blank = "blank";
let swapID;

function doAction(action) {
    const slot = "slot-" + action;
    let target = document.getElementById(slot);
    if (target) {
        htmx.swap("#" + slot, "", { swapStyle: 'delete' });
        return
    }

    htmx.trigger("#" + action, "click");
}

var hideLeft = true;
function toggleMenu(id) {
    if (hideLeft) {
        hideLeft = false;
        htmx.removeClass("#" + id, "hide")
    } else {
        hideLeft = true;
        htmx.addClass("#" + id, "hide")
    }
}


var hideChat = true;
function toggleChat(id) {
    if (hideChat) {
        hideChat = false;
        htmx.removeClass("#" + id, "hide")
    } else {
        hideChat = true;
        htmx.addClass("#" + id, "hide")
    }
}

function startTime() {
    const today = new Date();
    const h = today.getHours();
    let m = today.getMinutes();
    m = (m < 10) ? "0" + m : m;
    document.getElementById('clock').innerHTML = h + ":" + m;
    setTimeout(startTime, 1000 * (60 - today.getSeconds()));
}

// const chatId = "chat";
let drag_data = {};
let chat_data = {};
let slots = new Map();

function dragstartHandler(ev) {
    drag_data.offsetX = ev.offsetX;
    drag_data.offsetY = ev.offsetY;
}

function dragendHandler(ev) {
    const target = ev.target;
    const id = target.id;
    let data = {};
    data.X = ev.clientX - drag_data.offsetX;
    data.Y = ev.clientY - drag_data.offsetY;
    slots.set(id, data);
    target.style.left = data.X + 'px';
    target.style.top = data.Y + 'px';
    setdraggable(ev.target.id, false);
}

function addDragHandlers(id) {
    const target = document.getElementById(id);
    console.log("addDragHandlers", id, target);
    target.addEventListener("dragstart", dragstartHandler);
    target.addEventListener("dragend", dragendHandler);
}

function removeDragHandlers(id) {
    const target = document.getElementById(id);
    target.removeEventListener("dragstart", dragstartHandler);
    target.removeEventListener("dragend", dragendHandler);
}

function setdraggable(id, dragabble) {
    document.getElementById(id).setAttribute('draggable', dragabble);
}

window.addEventListener("DOMContentLoaded", () => {
    addDragHandlers("chat");
});

window.addEventListener('htmx:load', function (evt) {
    const target = evt.detail.elt;
    let id = target.id;
    if (!id.startsWith('slot-')) {
        return;
    }
    if (slots.has(id)) {
        data = slots.get(id);
        target.style.left = data.X + 'px';
        target.style.top = data.Y + 'px';
    }
    addDragHandlers(id);
})

window.addEventListener("gamepadconnected", (e) => {
    console.log(
      "Gamepad connected at index %d: %s. %d buttons, %d axes.",
      e.gamepad.index,
      e.gamepad.id,
      e.gamepad.buttons.length,
      e.gamepad.axes.length,
    );
});

window.addEventListener("gamepaddisconnected", (e) => {
    console.log(
      "Gamepad disconnected from index %d: %s",
      e.gamepad.index,
      e.gamepad.id,
    );
});
  