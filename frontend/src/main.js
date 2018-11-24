import 'bootstrap'
import 'bootstrap/dist/css/bootstrap.min.css'

import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import api from './api'

Vue.config.productionTip = false

let vueInstance = new Vue({
  router,
  store,
  render: h => h(App),
  methods: {
    onWsClose: function () {
      console.log('websocket closed, setting disconnected')
    },

    onWsOpen: function () {
      console.log('websocket open, setting connected');
    }
  }
}).$mount('#app')

api.onOpen(vueInstance.onWsOpen)
api.onClose(vueInstance.onWsClose)
