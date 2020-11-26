<template>
  <div class="container" ref="container">
    <slot></slot>
  </div>
</template>

<style scoped>
.container {
  height: 100%;
  overflow: auto;
}
</style>

<script lang="ts">
import {
  defineComponent, onUpdated, ref,
} from "vue";

export default defineComponent({
  name: "SScroller",
  props: {
    autoscroll: {
      type: Boolean,
      default: false,
    }
  },
  setup(props, { attrs }) {
    const container = ref(null);

    onUpdated(() => {
       if (
        container.value.scrollTop < container.value.scrollHeight &&
        props.autoscroll
      ) {
        container.value.scrollTop = container.value.scrollHeight;
      }
    });

    return {
      container,
    };
  },
});
</script>
