<template>
  <button class="button">
    <slot></slot>
  </button>
</template>

<style scoped>
.button {
  margin-left: 24px;
  background-color: rgb(64, 158, 255);
  color: rgb(255, 255, 255);
  font-weight: 500;
  padding: 12px 20px;
  appearance: none;
  cursor: pointer;
  font-size: 14px;
  border: 1px solid rgb(220, 223, 230);
  border-radius: 4px;
  transition: all 0.1s ease 0s;
}

.button:hover {
  background-color: rgb(102, 177, 255);
  border-color: rgb(102, 177, 255);
  outline: none;
}

.button:active {
  background-color: rgb(64, 158, 255);
}

.button:foucs {
  outline: none;
}
</style>

<script lang="ts">
import { defineComponent } from "vue";
export default defineComponent({
  name: "SButton",
});
</script>
