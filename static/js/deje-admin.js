// Admin interface code
define(['djadmin/logger', 'socket.io'], function(logger, io) {

function get_sio(resource) {
    return [
        window.location.protocol,
       '//',
        window.location.hostname,
        ':3001',
        resource
    ].join('');
}

function IRCLogger(url) {
    this.url = url;
    this.socket = io.connect(get_sio('/irc'));
    this.socket.on("connect",    this.on_connect.bind(this));
    this.socket.on("disconnect", this.on_disconnect.bind(this));
    this.socket.on("output",     this.on_output.bind(this));
    this.socket.on("error",      this.on_error.bind(this));
}

IRCLogger.prototype.on_connect = function() {
    logger.print("Connected to DEJE Admin server");
    this.subscribe(this.url);
}

IRCLogger.prototype.on_disconnect = function() {
    logger.print("Disconnected from DEJE Admin server");
}

IRCLogger.prototype.on_output = function(line) {
    logger.print(line);
}
IRCLogger.prototype.on_error = function(line) {
    logger.print("ERROR: " + line);
}
IRCLogger.prototype.subscribe = function(url) {
    this.url = url;
    logger.print("Subscribing to: " + url);
    this.socket.emit("subscribe", url);
}

return {
    'irc': IRCLogger
}

});
