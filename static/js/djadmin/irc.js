define(['djadmin/logger', 'djadmin/io'], function(Logger, socket) {

function IRCLogger(logging_selector) {
    this.logger = new Logger(logging_selector);

    socket.on("connect",    this.on_connect.bind(this));
    socket.on("disconnect", this.on_disconnect.bind(this));
    socket.on("output",     this.on_output.bind(this));
    socket.on("error",      this.on_error.bind(this));
}

IRCLogger.prototype.print = function(data) {
    this.logger.print(data);
}

IRCLogger.prototype.on_connect = function() {
    this.print("Connected to DEJE Admin server");
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

return IRCLogger;

});
