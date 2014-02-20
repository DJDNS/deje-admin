// Admin interface code
(function(){

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
}

IRCLogger.prototype.print = function(data) {
    var logger, old_content, new_content;
    logger = $('#irc_logger');
    old_content = logger.text();

    if (old_content.length == "") {
        new_content = data;
    } else {
        new_content = old_content + '\n' + data;
    }
    logger.text(new_content);
}

IRCLogger.prototype.on_connect = function() {
    this.print("Connected to DEJE Admin server");
    this.subscribe(this.url);
}

IRCLogger.prototype.on_disconnect = function() {
    this.print("Disconnected from DEJE Admin server");
}

IRCLogger.prototype.on_output = function(line) {
    this.print(line);
}
IRCLogger.prototype.subscribe = function(url) {
    this.url = url;
    this.print("Subscribing to: " + url);
    this.socket.emit("subscribe", url);
}


window.djadmin = {
    'get_sio': get_sio,
    'irc': IRCLogger
}

}());
