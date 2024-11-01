
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

    wsstat.innerText = 'Connected';
    wsstat.setAttribute('ws-status', 'connected');
});
document.body.addEventListener('htmx:wsError', function(evt) {
    console.log('error');
    wsstat.innerText = 'Error';
    wsstat.setAttribute('ws-status', 'error');
});
document.body.addEventListener('htmx:wsClose', function(evt) {
    console.log('disconnected');
    wsstat.innerText = 'Disconnected';
    wsstat.setAttribute('ws-status', 'disconnected');
});
