<template>
  <div ref="container" class="container"></div>
</template>

<style scoped>
.container {
  height: 100%;
  overflow: hidden;
}
</style>

<script lang="ts">
import {
  defineComponent, onMounted, ref, watchEffect,
} from "vue";

import { http, httpPut, HttpResponse } from "../utils";

import * as monaco from "monaco-editor";
import EditorWorker from "monaco-editor/esm/vs/editor/editor.worker.js?worker";
import JsonWorker from "monaco-editor/esm/vs/language/json/json.worker.js?worker";

// @ts-ignore
self.MonacoEnvironment = {
  getWorker: function (moduleId, label) {
    if (label === "json") {
      return new JsonWorker();
    }
    return new EditorWorker();
  },
};

export default defineComponent({
  name: "Editor",
  props: {
    filepath: String,
  },
  emits: [
    'saved',
  ],
  setup(props, { attrs, emit }) {
    const container = ref(null);
    var editor = null;

    watchEffect(async () => {
      if (props.filepath == null) {
        return;
      }
      let response = await http(`/api/files?path=${props.filepath}`, "text");
      editor.setValue(response.parsedBody);
    });

    onMounted(() => {
      editor = monaco.editor.create(container.value, {
        value: "",
        language: "text",
      });
    });

    const save = async () => {
      const response = await httpPut(`/api/files?path=${props.filepath}`, editor.getValue());
    };

    return {
      container,
      save,
    };
  },
});
</script>
