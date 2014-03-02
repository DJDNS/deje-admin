define(['socket.io'], function(io) {

var socket_url = [
    window.location.protocol,
   '//',
    window.location.hostname,
    ':3001/',
].join('');

var socket = io.connect(socket_url);
socket.subscription = ""

io.SocketNamespace.prototype.subscribe = function(url) {
    this.subscription = url;
    if (this.socket.connected) {
        this.emit("subscribe", url);
    }
}

socket.on('connect', function(){
    if (socket.subscription !== "") {
        socket.subscribe(socket.subscription);
    }
});

return socket;

});
