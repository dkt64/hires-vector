import Vue from 'vue'
import App from './App.vue'
// import axios from "axios";
import router from './router'
import store from './store'
import vuetify from './plugins/vuetify';

Vue.config.productionTip = false

// if (process.env.NODE_ENV == "development") {
//   axios.defaults.baseURL = "http://localhost:8080";
//   // axios.defaults.baseURL = "http://raspberrypi-jost-01:1057"
// }

new Vue({
  router,
  store,
  vuetify,
  render: h => h(App)
}).$mount('#app')
