<template>
  <transition name="fade">
    <div
      class="notification"
      v-show="visible"
      ref="el"
      @mouseenter="clear"
      @mouseleave="start"
    >
      <div class="container">
        <div class="message">
          {{ message }}
        </div>
        <div class="close">
          <i class="material-icons button" @click="close">highlight_off</i>
        </div>
      </div>
    </div>
  </transition>
</template>

<style scoped>
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 1s ease;
}

.notification {
  position: absolute;
  width: 400px;
  left: 50%;
  margin-left: -200px;
  bottom: 24px;
  border: 1px solid rgba(255, 0, 0, 0.4);
  border-radius: 4px;
  background-color: rgba(255, 0, 0, 0.1);
  transition: all 1s;
}

.container {
  display: flex;
  margin: 5px;
}

.message {
  flex-grow: 1;
}

.close {
  justify-self: flex-end;
  margin-left: 24px;
  padding: 5px;
}

.button {
  cursor: pointer;
}
</style>

<script lang="ts">
import {
  createVNode,
  defineComponent,
  render,
  PropType,
  onMounted,
  ref,
  Ref,
} from "vue";

export default defineComponent({
  name: "SNotification",
  props: {
    message: String,
    onClose: {
      type: Function as PropType<() => void>,
      required: true,
    },
    duration: {
      type: Number,
      default: 3000,
    },
    id: String,
  },
  setup(props, { attrs }) {
    const el: Ref<HTMLDivElement> = ref(null);
    const visible: Ref<boolean> = ref(false);
    const timer: Ref<number> = ref(null);

    const destroy = () => {
      visible.value = false;
      el.value.removeEventListener("transitionend", destroy);
      props.onClose();
    };

    const close = () => {
      visible.value = false;
      el.value.addEventListener("transitionend", destroy);
    };

    const start = () => {
      if (props.duration > 0) {
        timer.value = setTimeout(() => {
          close();
        }, props.duration);
      }
    };

    const clear = () => {
      clearTimeout(timer.value);
      timer.value = null;
    };

    onMounted(() => {
      visible.value = true;
      start();
    });

    return {
      el,
      visible,
      close,
      destroy,
      timer,
      start,
      clear,
    };
  },
});
</script>
