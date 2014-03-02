define(['jquery'], function($) {

function Inputter(form_selector) {
    this.form_element = $(form_selector);
    this.latency = 200; // milliseconds
    this.timeout = null;
    this.getTextarea().on(
        'keyup change input propertychange',
        this.delay_evaluate.bind(this)
    );
}

Inputter.prototype.getTextarea = function() {
    return this.form_element.find('textarea');
}
Inputter.prototype.getMsgarea = function() {
    return this.form_element.find('.help-block');
}

Inputter.prototype.enable = function() {
    this.form_element.children('.form-group')
        .removeClass('has-error')
        .addClass('has-success');
    this.form_element.children('fieldset').removeAttr('disabled');
    this.getMsgarea().text("That's valid. Good job.");
}
Inputter.prototype.disable = function(msg) {
    this.form_element.children('.form-group')
        .addClass('has-error')
        .removeClass('has-success');
    this.form_element.children('fieldset').attr('disabled', true);
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


return {
    "Inputter": Inputter
}

});
