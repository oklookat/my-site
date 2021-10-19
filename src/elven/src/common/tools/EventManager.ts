// https://stackoverflow.com/a/58090101/16762009

class EventManager {

    private events: object

    constructor() {
        this.events = {};
    }

    public add(event, callback) {
        const map = this.events[event] = this.events[event] || new Map;
        if (typeof callback === 'function') {
            map.set(callback);
        }
    }

    public remove(event, callback) {
        const map = this.events[event];
        if (map) {
            map.delete(callback);
        }
    }

    public fire(event, data) {
        const map = this.events[event];
        if (map) {
            [...map].forEach(([cb]) => cb(data));
        }
    }
}