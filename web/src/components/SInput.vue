<template>
  <div class="container">
    <input
      :type="type"
      name="input"
      placeholder=" "
      :value="modelValue"
      @input="$emit('update:modelValue', $event.target.value)"
    />
    <label for="input" class="label">
      <span class="content">{{ placeholder }}</span>
    </label>
  </div>
</template>

<style scoped>
.container {
  width: 100%;
  height: 42px;
  position: relative;
  overflow: auto;
}

.container input {
  width: 100%;
  height: 100%;
  color: #595f6e;
  padding-top: 20px;
  border: none;
  box-sizing: border-box;
  outline: none;
}

.container label {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  border-bottom: 1px solid black;
}

.container label::after {
  content: "";
  height: 100%;
  width: 100%;
  position: absolute;
  left: 0;
  bottom: -1px;
  border-bottom: 2px solid rgb(64, 158, 255);
  transform: scaleX(0);
  transition: .3s cubic-bezier(.25,.8,.5,1);
}

.content {
  position: absolute;
  bottom: 5px;
  left: 0;
  transition: .3s cubic-bezier(.25,.8,.5,1);
}

.container input:focus + .label .content, .container input:not(:placeholder-shown) + .label .content {
  transform: translateY(-100%);
  font-size: 14px;
  color: rgb(64, 158, 255);
}

.container input:focus + label::after, .container input:not(:placeholder-shown) + label::after {
  transform: scaleX(1);
}
</style>

<script lang="ts">
import { computed, defineComponent, ref, watchEffect } from "vue";
export default defineComponent({
  name: "SInput",
  props: {
    modelValue: String,
    placeholder: String,
    type: {
      type: String,
      required: true,
    }
  },
  emits: ["update:modelValue"],
});
</script>
