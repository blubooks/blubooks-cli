import { createApp } from 'vue'
import { registerPlugins } from './plugins'
import App from './App.vue'


//if (import.meta.env.VITE_APP_PRODUKT) {
//    import(/* @vite-ignore */"./styles/App." + import.meta.env.VITE_APP_PRODUKT + ".scss");
//}else {
//    import(/* @vite-ignore */"./styles/App.pruefsoftware.scss");
//}
import(/* @vite-ignore */"./styles/default/App.default.scss");



const app = createApp(App)
registerPlugins(app)

app.mount('#app')
