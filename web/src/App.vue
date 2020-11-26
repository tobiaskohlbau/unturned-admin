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
import RCON from './components/RCON.vue'
import Manager from './components/Manager.vue'

export default {
  name: 'App',
  components: {
    'rcon': RCON,
    'manager': Manager,
  }
}
</script>
