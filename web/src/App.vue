<template>
  <div class="container">
    <router-view />
  </div>
</template>

<style scoped>
.container {
  display: flex;
  align-content: center;
  min-height: 100vh;
}
</style>

<script>
import RCON from "./components/RCON.vue";
import Manager from "./components/Manager.vue";
import { onMounted } from 'vue';
import { useStore } from './store';
import { getToken } from './utils';

export default {
  name: "App",
  components: {
    rcon: RCON,
    manager: Manager,
  },
  setup() {
    const store = useStore();
    onMounted(() => {
      const token = getToken();
      if (token !== null) {
        store.setToken(token);
      }
    });
  },
};
</script>
