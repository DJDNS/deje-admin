// Logs to page IRC log
define(['jquery'], function($) {

function Logger(selector) {
    this.element = $(selector);
}

Logger.prototype.print = function(data) {
    var old_content, new_content;
    old_content = this.element.text();

    if (old_content.length == "") {
        new_content = data;
    } else {
        new_content = old_content + '\n' + data;
    }
    this.element.text(new_content);
}

return Logger;
});
