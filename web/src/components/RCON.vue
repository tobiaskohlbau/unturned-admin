<template>
  <div v-bind:style="{ width: width, height: height }" class="card">
    <div class="card_content">
      <div ref="logContainer" class="log_container">
        <div v-for="log in logs" :key="log" class="msg">
          {{ log }}
        </div>
      </div>
      <div class="settings">
        <input type="checkbox" v-model="autoscroll" id="autoscroll" />
        <label for="autoscroll">Autoscroll</label>
      </div>
      <div class="cli_container" v-on:keyup.enter="submit">
        <input v-model="command" placeholder="Please input command" class="input"/>
        <button @click.stop="submit" class="button">Send</button>
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

.settings {
  display: flex;
  justify-content: flex-end;
  margin: 8px;
}

.log_container {
  height: 100%;
  overflow: auto;
}

.msg {
  margin-top: 3px;
  padding: 5px;
  border-radius: 3px;
  background: #f5f5f5;
}

.card_content {
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.cli_container {
  display: flex;
}

#autoscroll {
  border-color: rgb(64, 158, 255);
  border-color: rgb(64, 158, 255);
}

.input {
  border-radius: 4px;
  border: 1px solid rgb(220, 223, 230);
  color: rgb(96, 98, 102);
  display: inline-block;
  height: 40px;
  line-height: 40px;
  outline: 0px;
  appearance: none;
  padding: 0px 15px;
  transition: border-color 0.2s cubic-bezier(0.645, 0.045, 0.355, 1) 0s;
  width: 100%;

  &::placeholder {
    color: rgb(192, 196, 204)
  }
}

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

  &:hover {
    background-color: rgb(102, 177, 255);
    border-color: rgb(102, 177, 255);
    outline: none;
  }

  &:active {
    background-color: rgb(64, 158, 255);
  }

  &:foucs {
    outline: none;
  }
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

var backendURL = "ws://localhost:8080/ws";

if (import.meta.env.MODE == "production") {
  backendURL = "wss://" + document.location.host + "/ws";
}

export default defineComponent({
  name: "rcon",
  props: {
    width: String,
    height: String,
  },
  setup({}, { attrs }) {
    const logContainer: Ref<HTMLDivElement> = ref(null);
    const logs = ref([]);
    const command = ref("");
    const autoscroll = ref(true);
    var conn;

    onUpdated(() => {
      if (
        logContainer.value.scrollTop < logContainer.value.scrollHeight &&
        autoscroll.value
      ) {
        logContainer.value.scrollTop = logContainer.value.scrollHeight;
      }
    });
    onMounted(() => {
      console.log("connecting to ", backendURL);
      conn = new WebSocket(backendURL);
      conn.onerror = (event: ErrorEvent) => {
        Message({
          type: "error",
          message:
            "Could not establish WebSocket connection. Try reloading the page.",
          center: true,
        });
        conn = null;
      };
      conn.onclose = function (evt) {
        if (conn == null) {
          return;
        }
        logs.value.push("Connection closed.");
      };
      conn.onmessage = (event: MessageEvent<string>) => {
        logs.value.push(...event.data.split("\n"));
      };
    });

    const submit = () => {
      if (!conn) {
        return;
      }
      if (command.value === "") {
        logs.value.push("Need to specify command");
        return;
      }
      conn.send(command.value);
      command.value = "";
    };
    return {
      logContainer,
      conn,
      logs,
      command,
      submit,
      autoscroll,
    };
  },
});
</script>
