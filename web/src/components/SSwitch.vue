<template>
  <div class="wrapper">
    <label for="input" class="label">
      <slot></slot>
    </label>
    <div class="container" @click="$emit('update:modelValue', !modelValue)">
      <input type="checkbox" name="input" :value="modelValue" />
      <div :class="['slider', modelValue ? 'on' : 'off']"></div>
      <div :class="['point', modelValue ? 'on' : 'off']"></div>
    </div>
  </div>
</template>

<style scoped lang="postcss">
.wrapper {
  display: flex;
  align-items: center;
}

.container {
  width: 38px;
  height: 24px;
  position: relative;
  overflow: hidden;
  display: inline-block;
}

.label {
  margin-right: 8px;
}

.container input {
  position: absolute;
  opacity: 0;
}

.slider {
  position: absolute;
  width: 100%;
  height: 14px;
  top: calc(50% - 7px);
  box-sizing: border-box;
  border-radius: 12px;

  &.on {
    background-color: #54B2FF;
  }

  &.off {
    background-color: #6D7382;
  }
}

.point {
  position: absolute;
  width: 20px;
  height: 20px;
  top: calc(50% - 10px);
  box-sizing: border-box;
  border-radius: 50%;
  transition: 0.3s cubic-bezier(0.25, 0.8, 0.5, 1);
}

.point.on {
  background-color: rgb(64, 158, 255);
  transform: translateX(18px);
}

.point.off {
  background-color: #595f6e;
}
</style>

<script lang="ts">
import { computed, defineComponent, Ref, ref, watchEffect } from "vue";
export default defineComponent({
  name: "SSwitch",
  props: {
    modelValue: Boolean,
  },
  emits: ["update:modelValue"],
});
</script>
