<template>
  <div class="popup" v-if="show">
    <div class="popup_content">
      <slot class="popup_content"></slot>
    </div>
  </div>
</template>

<style scoped>
.popup {
  position: absolute;
  left: 0;
  top: 0;
  background-color: rgba(22, 22, 22, 0.5);
  width: 100%;
  height: 100%;
  z-index: 100;
}

.popup_content {
  display: inline-block;
  position: relative;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  -webkit-box-shadow: 0 0 10px rgba(0, 0, 0, 0.4);
  -moz-box-shadow: 0 0 10px rgba(0, 0, 0, 0.4);
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.4);
}
</style>

<script lang="ts">
import {
  defineComponent, onUpdated, ref,
} from "vue";

export default defineComponent({
  name: "SPopup",
  props: {
    show: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { attrs }) {
    return {};
  },
});
</script>
