<template>
  <div v-bind:style="{ width: width, height: height }" :class="$style.card">
    <slot></slot>
  </div>
</template>

<style module>
.card {
  box-shadow: rgba(0, 0, 0, 0.1) 0px 2px 12px 0px;
  padding: 20px;
  border-radius: 4px;
  border: 1px solid rgb(235, 238, 245);
  background-color: rgb(255, 255, 255);
  color: rgb(48, 49, 51);
  transition: all 0.3s ease 0s;
  display: flex;
  flex-direction: column;
  min-width: 240px;
}
</style>

<script lang="ts">
import {
  defineComponent,
} from "vue";

export default defineComponent({
  name: "SCard",
  props: {
    width: String,
    height: String,
  },
});
</script>
