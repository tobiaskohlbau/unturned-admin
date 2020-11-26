<template>
  <div :class="['container', className]">
    <rcon></rcon>
    <manager></manager>
  </div>
</template>

<style scoped lang="postcss">
.container {
  width: 100%;
  display: flex;
  flex-wrap: wrap;
  justify-content: space-around;
}

.flex-gap {
  column-gap: 24px;
  row-gap: 24px;
}

.flex-no-gap {
  margin-top: -24px;
  margin-right: -24px;

  & > * {
    margin-top: 24px;
    margin-right: 24px;
  }
}
</style>

<script>
import RCON from '../components/RCON.vue'
import Manager from '../components/Manager.vue'
import { defineComponent, ref } from 'vue'
import { useRouter } from 'vue-router';
import { checkFlexGap } from '../utils';

export default defineComponent({
  name: 'Home',
  components: {
    'rcon': RCON,
    'manager': Manager,
  },
  setup() {
    const className = ref("flex-gap");

    // Workaround to support safari.
    // Drop if safari supports gap on flex items,
    // currently in technical preview.
    if (!checkFlexGap()) {
      className.value = "flex-no-gap";
    }

    return {
      className,
    }
  }
});
</script>
