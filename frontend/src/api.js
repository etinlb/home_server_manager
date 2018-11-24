// import _ from 'lodash'

class API {
  constructor (path) {
    let that = this
    this.callbackIdMap = {};

    this.connected = new Promise((resolve, reject) => {
      that.sock = new WebSocket(API.qualifyWebsocketURL(path))
      that.sock.onopen = function () {
        resolve('Success')
      }
      that.sock.onerror = function () {
        reject(new Error('Oh no!'))
      }
    })
    this.sock.onmessage = function (e) {
      let msg = JSON.parse(e.data);
      console.log(msg)
    }

    this.sock.onclose = function (e) {
      console.log(e)
    }

    // Thanks stack overflow
    // https://stackoverflow.com/questions/105034/create-guid-uuid-in-javascript
    this.getUUID = function() {
      return ([1e7]+-1e3+-4e3+-8e3+-1e11).replace(/[018]/g, c =>
        (c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> c / 4).toString(16)
      )
    }
  }

  sendMessage(msg) {
    this.sock.send(msg);
  }

  onOpen (callback) {
    this.connected.then(callback)
  }

  onClose (callback) {
    console.log(callback)
  }

  startCleanup(callback) {
    let callbackId = this.getUUID();
    let msg = this.createMessage(callbackId, "cleanup", {});
    this.callbackIdMap[callbackId] = callback;

    this.sock.send(msg);
  }

  createMessage(callbackId, action, parameters) {
    return JSON.stringify({
      "id" : callbackId,
      "action": action,
      "parameters": parameters,
    });
  }

  static qualifyWebsocketURL (path) {
    let protocol
    if (window.location.protocol === 'http:') {
      protocol = 'ws:'
    } else {
      protocol = 'wss:'
    }
    // Note that we need to handle the case where we may be running not at the root
    // so any path should get added on to the base url for the api. This is
    // especially the case for reversex proxies and such
    let apiBase = protocol + window.location.host + window.location.pathname
    return new URL(path, apiBase).toString()
  }
}

export default new API('/api/')
