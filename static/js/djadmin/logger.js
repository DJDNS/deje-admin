// Logs to page IRC log
define(['jquery'], function($) {

var logger = {
    "element": null
};

logger.setup_element = function() {
    if (logger.element === null) {
        logger.element = $('#irc_logger');
    }
}

logger.print = function(data) {
    var element, old_content, new_content;
    logger.setup_element();

    element = logger.element;
    old_content = element.text();

    if (old_content.length == "") {
        new_content = data;
    } else {
        new_content = old_content + '\n' + data;
    }
    element.text(new_content);
}

$(document).ready(logger.setup_element);
return logger;
});
