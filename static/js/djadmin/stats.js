define(['jquery','djadmin/io'], function ($,socket) {

function StatsUpdater(selector) {
    this.element = $(selector);
    socket.on("stats", this.on_update.bind(this));
    this.interval = setInterval(
        this.request_update.bind(this),
        2000 //ms
    );
}

StatsUpdater.prototype.on_update = function(data) {
    var counts = ['events', 'quorums', 'ts'];
    for (var i = 0; i<counts.length; i++) {
        var name = 'num_' + counts[i];
        this.element.find('.' + name).text(data[name]);
    }
}
StatsUpdater.prototype.request_update = function() {
    socket.emit("stats_request")
}


return StatsUpdater;

});
