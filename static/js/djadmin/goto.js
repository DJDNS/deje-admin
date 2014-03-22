define(['djadmin/io'], function(socket) {

function goto_event(hash) {
    socket.emit("goto_request", hash);
}

return goto_event;

});
