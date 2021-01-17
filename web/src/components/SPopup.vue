<template>
  <div>
    <slot name="activator" class="popup_actor" :on="listener"></slot>
    <div class="popup" v-if="visible">
      <div class="popup_content" v-click-outside="closeConditional">
        <slot class="popup_content"></slot>
      </div>
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
import { defineComponent, onUpdated, ref, toRefs, watch } from "vue";

export default defineComponent({
  name: "SPopup",
  props: {
    persistent: {
      type: Boolean,
      default: false,
    },
  },
  emits: ["update:modelValue"],
  setup(props, { attrs, emit }) {
    const visible = ref(false);

    const listener = {
      click: (e: MouseEvent) => {
        visible.value = true;
      }
    };

    const open = () => {
      visible.value = true;
    };

    const close = () => {
      visible.value = false;
    }

    const closeConditional = () => {
      if (props.persistent !== true) {
        visible.value = false;
      }
    }

    return {
      visible,
      listener,
      open,
      close,
      closeConditional,
    };
  },
});
</script>
