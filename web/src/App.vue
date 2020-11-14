<template>
  <div class="app">
    <rcon width="600px" height="400px"></rcon>
    <screenshot width="600px" height="400px"></screenshot>
  </div>
</template>

<style lang="scss" scoped>
.app {
  display: flex;
  justify-content: space-around;
  flex-wrap: wrap;
  margin-top: -24px;

  > * {
    margin-top: 24px;
  }
}
</style>

<script>
import RCON from './components/RCON.vue'
import Screenshot from './components/Screenshot.vue'

export default {
  name: 'App',
  components: {
    'rcon': RCON,
    'screenshot': Screenshot,
  }
}
</script>
