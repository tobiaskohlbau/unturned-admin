import { createApp } from 'vue'
import App from './App.vue'
import './index.css'
import router from './router';

import { SCard, SCardTitle, SCardContent, SCardActions } from './components/SCard';
import SButton from './components/SButton.vue';
import SInput from './components/SInput.vue';
import SSwitch from './components/SSwitch.vue';
import SScroller from './components/SScroller.vue';
import SPopup from './components/SPopup.vue';
import SSelect from './components/SSelect.vue';
import { storeSymbol, createStore } from './store';
import { ClickOutside } from './directives';

const app = createApp(App);

app.component('SCard', SCard);
app.component('SCardTitle', SCardTitle);
app.component('SCardContent', SCardContent);
app.component('SCardActions', SCardActions);
app.component('SButton', SButton);
app.component('SInput', SInput);
app.component('SSwitch', SSwitch);
app.component('SScroller', SScroller);
app.component('SPopup', SPopup);
app.component('SSelect', SSelect);

app.provide(storeSymbol, createStore());

app.use(router);

app.directive('click-outside', ClickOutside);

if ('serviceWorker' in navigator) {
  window.addEventListener('load', function() {
    navigator.serviceWorker.register('/sw.js').then(function(registration) {
    }, function(err) {
      // registration failed :(
      console.log('ServiceWorker registration failed: ', err);
    });
  });
}

app.mount('#app')
