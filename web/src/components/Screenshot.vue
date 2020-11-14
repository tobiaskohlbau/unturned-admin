<template>
  <div v-bind:style="{ width: width, height: height }" class="card">
    <div class="card_content">
      <div v-for="screenshot in screenshots" :key="screenshot">
        {{ screenshot }}
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.card {
  box-shadow: rgba(0, 0, 0, 0.1) 0px 2px 12px 0px;
  padding: 20px;
  border-radius: 4px;
  border: 1px solid rgb(235, 238, 245);
  background-color: rgb(255, 255, 255);
  color: rgb(48, 49, 51);
  transition: all 0.3s ease 0s;
}

.card_content {
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  overflow: auto;
}
</style>

<script lang="ts">
import {
  defineComponent,
  onBeforeMount,
  onBeforeUpdate,
  onMounted,
  onUpdated,
  reactive,
  ref,
  Ref,
} from "vue";

import { Message } from "element-plus/lib/message";
import { http, HttpResponse } from "../utils";

export default defineComponent({
  name: "Screenshot",
  props: {
    width: String,
    height: String,
  },
  setup({}, { attrs }) {
    const screenshots = ref([]);

    onMounted(async () => {
      let response: HttpResponse<string[]>;
      try {
        response = await http<string[]>("/api/screenshots");
        screenshots.value.push(...response.parsedBody)
      } catch (err) {
        screenshots.value.push("Failed to fetch screenshots!", err);
      }
    });

    return {
      screenshots,
    };
  },
});
</script>
