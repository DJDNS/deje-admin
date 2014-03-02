define(['djadmin/logger', 'socket.io'], function(Logger, io) {

function get_sio(resource) {
    return [
        window.location.protocol,
       '//',
        window.location.hostname,
        ':3001',
        resource
    ].join('');
}

function IRCLogger(logging_selector) {
    this.url = "";
    this.logger = new Logger(logging_selector);
    this.socket = io.connect(get_sio('/irc'));
    this.socket.on("connect",    this.on_connect.bind(this));
    this.socket.on("disconnect", this.on_disconnect.bind(this));
    this.socket.on("output",     this.on_output.bind(this));
    this.socket.on("error",      this.on_error.bind(this));
}

IRCLogger.prototype.print = function(data) {
    this.logger.print(data);
}

IRCLogger.prototype.on_connect = function() {
    this.print("Connected to DEJE Admin server");
    if (this.url != "") {
        this.subscribe(this.url);
    }
}

IRCLogger.prototype.on_disconnect = function() {
    this.print("Disconnected from DEJE Admin server");
}

IRCLogger.prototype.on_output = function(line) {
    this.print(line);
}
IRCLogger.prototype.on_error = function(line) {
    this.print("ERROR: " + line);
}
IRCLogger.prototype.subscribe = function(url) {
    this.url = url;
    if (this.socket.socket.connected) {
        this.print("Subscribing to: " + url);
        this.socket.emit("subscribe", url);
    }
}

return IRCLogger;

});
