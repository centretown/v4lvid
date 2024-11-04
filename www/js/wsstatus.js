
let wsstat = document.getElementById('ws-status');

// htmx:wsConnecting
// htmx:wsError

let socket;
let elt;

document.addEventListener("visibilitychange", function(evt) {
    console.log('visibilitychange', document.visibilityState);
    if (socket) {
        socket.send(document.visibilityState, elt);
    }
});

document.body.addEventListener('htmx:wsOpen', function(evt) {
    console.log('connected');
    socket = evt.detail.socketWrapper;
    elt = evt.detail.elt;
    wsstat.innerText = 'chat';
    wsstat.setAttribute('ws-status', 'connected');
});
document.body.addEventListener('htmx:wsError', function(evt) {
    console.log('error');
    wsstat.innerText = 'chat_error';
    wsstat.setAttribute('ws-status', 'error');
});
document.body.addEventListener('htmx:wsClose', function(evt) {
    console.log('disconnected');
    wsstat.innerText = 'chat_error';
    wsstat.setAttribute('ws-status', 'disconnected');
});

