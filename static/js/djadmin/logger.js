// Logs to page IRC log
define(['jquery'], function($) {

function Logger(selector) {
    this.element = $(selector);
}

Logger.prototype.scroll_position = function() {
    return this.element.scrollTop();
}
Logger.prototype.scroll_max = function() {
    var el = this.element;
    return el.prop('scrollHeight') - el.innerHeight();
}
Logger.prototype.is_following = function() {
    return (this.scroll_position() === this.scroll_max());
}

Logger.prototype.print = function(data) {
    var old_content, new_content, following;
    old_content = this.element.text();
    following = this.is_following();

    if (old_content.length == "") {
        new_content = data;
    } else {
        new_content = old_content + '\n' + data;
    }
    this.element.text(new_content);

    if (following) {
        this.element.scrollTop(this.scroll_max());
    }
}

return Logger;
});
