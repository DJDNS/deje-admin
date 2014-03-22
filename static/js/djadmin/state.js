define(['jquery', 'djadmin/io'], function($, socket) {

var test_display_value = {
    "paxos": {
        "acceptors" : [
            "rspode@placeholder.me",
            "aglossop@placeholder.me"
        ],
        "learners" : [],
    },
    "handlers": {},
    "misc": {
        "true": true,
        "false": false,
        "null": null,
        "a float": -98.5,
        "an int": 17,
    }
};

// Expects a <div>
function StateDisplay(display_selector) {
    this.element = $(display_selector);
    this.value = {};

    socket.on("primitive", this.on_primitive.bind(this));
    this.render_root();
}

StateDisplay.prototype.on_primitive = function(primitive) {
    this.value = primitive.args.Value;
    this.render_root();
}

StateDisplay.prototype.render_root = function() {
    this.render(this.value, this.element);
}

StateDisplay.prototype.render = function(obj, element) {
    var renderer = get_renderer(obj);
    element.html(
        renderer.render_body()
     || '<code>' + renderer.render_header() + '</code>'
    );
}

// Select a renderer based on object type
function get_renderer(obj) {
    if (typeof(obj) == "string") {
        return new StringRenderer(obj);
    } else if (typeof(obj) == "number") {
        return new NumberRenderer(obj);
    } else if (typeof(obj) == "boolean") {
        return new BooleanRenderer(obj);
    } else if (obj == null || obj == undefined) {
        return new NullRenderer(obj);
    } else if (Array.isArray(obj)) {
        return new ArrayRenderer(obj);
    } else {
        return new MapRenderer(obj);
    }
}

function StringRenderer(obj) {
    this.obj = obj;
}
StringRenderer.prototype.render_header = function() {
    return [
        '"',
        this.obj,
        '"'
    ].join('');
}
StringRenderer.prototype.render_body = function() {
    return '';
}

function NumberRenderer(obj) {
    this.obj = obj;
}
NumberRenderer.prototype.render_header = function() {
    return this.obj.toString();
}
NumberRenderer.prototype.render_body = function() {
    return '';
}

function BooleanRenderer(obj) {
    this.obj = obj;
}
BooleanRenderer.prototype.render_header = function() {
    return this.obj.toString();
}
BooleanRenderer.prototype.render_body = function() {
    return '';
}

function NullRenderer() {}
NullRenderer.prototype.render_header = function() {
    return "null";
}
NullRenderer.prototype.render_body = function() {
    return '';
}

function ArrayRenderer(obj) {
    this.obj = obj;
}
ArrayRenderer.prototype.render_header = function() {
    return (this.obj.length > 0) ? '[...]' : '[]';
}
ArrayRenderer.prototype.render_body = function() {
    var element, i, renderer, child;
    element = $('<ul>');
    for (i = 0; i < this.obj.length; i++) {
        renderer = get_renderer(this.obj[i]);
        child = $([
            '<li>',
                '<code>', i, '</code>',
                ' : <code>', renderer.render_header(), '</code>',
            '</li>'
        ].join(''));
        child.append(renderer.render_body());
        element.append(child);
    }
    return (this.obj.length > 0) ? element : '';
}

function MapRenderer(obj) {
    this.obj = obj;
}
MapRenderer.prototype.render_header = function() {
    var property;
    for (property in this.obj) {
        if (this.obj.hasOwnProperty(property)) {
            return "{...}"
        }
    }
    return "{}"
}
MapRenderer.prototype.render_body = function() {
    var element, property, hasany, renderer, child;
    element = $('<ul>');
    for (property in this.obj) {
        if (this.obj.hasOwnProperty(property)) {
            hasany = true;
            renderer = get_renderer(this.obj[property]);
            child = $([
                '<li>',
                    '<code>"', property, '"</code>',
                    ' : <code>', renderer.render_header(), '</code>',
                '</li>'
            ].join(''));
            child.append(renderer.render_body());
            element.append(child);
        }
    }
    return hasany ? element : '';
}

return StateDisplay;

});
