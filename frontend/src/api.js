import _ from 'lodash'

class API {
  constructor (path) {
    let that = this

    this.connected = new Promise((resolve, reject) => {
      that.sock = new WebSocket(API.qualifyWebsocketURL(path))
      that.sock.onopen = function () {
        resolve('Success')
      }
      that.sock.onerror = function () {
        reject(new Error('Oh no!'))
      }
    })
    this.notify_map = {}
    this.sock.onmessage = function (e) {
      console.log(e)
    }

    this.sock.onclose = function (e) {
      console.log(e)
    }
  }

  _notify (name, obj) {
    var cbs = this._subscribers[name]
    for (var ii = 0; ii < cbs.length; ii++) {
      cbs[ii](obj)
    }
  }

  onOpen (callback) {
    this.connected.then(callback)
  }

  onClose (callback) {
    console.log(callback)
  }

  onAbnormalModeTermination (callback) {
    this._subscribers.abnormalModeTermination.push(callback)
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
