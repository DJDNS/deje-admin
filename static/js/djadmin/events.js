define(['jquery','djadmin/io','djadmin/goto'],
function($,       socket,      goto_event) {

function Inputter(form_selector) {
    this.form_element = $(form_selector);
    this.latency = 200; // milliseconds
    this.timeout = null;
    this.template = {
        "parent": "123456789",
        "handler": "SET",
        "args": {
            "path": [],
            "value": { "hello": "world"}
        }
    }
    this.getTextarea().on(
        'keyup change input propertychange focus blur',
        this.delay_evaluate.bind(this)
    );
    this.form_element.on("submit", this.on_submit.bind(this));
    this.form_element.find('.btn-secondary').click(
        this.fill_from_template.bind(this)
    );
    socket.on("event_registered", this.on_register.bind(this));
    socket.on("event_error", this.on_error.bind(this));
}

Inputter.prototype.getTextarea = function() {
    return this.form_element.find('textarea');
}
Inputter.prototype.getMsgarea = function() {
    return this.form_element.find('.help-block');
}

Inputter.prototype.enable = function() {
    var form_group = this.form_element.children('.form-group');
    if (form_group.hasClass('has-success')) {
        return;
    }
    form_group.removeClass('has-error').addClass('has-success');
    this.form_element.find('fieldset .btn-primary').removeAttr('disabled');
    this.getMsgarea().text("That's valid. Good job.");
}
Inputter.prototype.disable = function(msg) {
    var form_group = this.form_element.children('.form-group');
    if (form_group.hasClass('has-error')) {
        return;
    }
    form_group.addClass('has-error').removeClass('has-success');
    this.form_element.find('fieldset .btn-primary').attr('disabled', "disabled");
    this.getMsgarea().text(msg || "Must be valid JSON.");
}

Inputter.prototype.evaluate = function() {
    var value;
    value = this.getTextarea().val();
    try {
        JSON.parse(value);
        this.enable();
    } catch(e) {
        this.disable();
    }
}
Inputter.prototype.delay_evaluate = function() {
    if (this.timeout !== null) {
        clearTimeout(this.timeout);
    }
    this.timeout = setTimeout(
        this.evaluate.bind(this),
        this.latency
    )
}

Inputter.prototype.on_submit = function(e) {
    var value;
    e.preventDefault();

    value = this.getTextarea().val();
    this.getMsgarea().text("Publishing event...");
    socket.emit("event", value);
}

Inputter.prototype.on_register = function(hash) {
    var msgarea = this.getMsgarea();
    msgarea.html([
        'Event ',
         '<a href="#" class="js-goto-hash">', hash, '</a>',
         ' published!'
    ].join(''));
    msgarea.find('a.js-goto-hash').click(function(e) {
        e.preventDefault();
        goto_event(hash);
    });
}
Inputter.prototype.on_error = function(msg) {
    this.getMsgarea().text(msg);
    //this.disable()
}

Inputter.prototype.fill_from_template = function() {
    this.getTextarea().text(
        JSON.stringify(this.template, undefined, "  ")
    );
}

return {
    "Inputter": Inputter
}

});
