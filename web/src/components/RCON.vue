<template>
  <s-card width="600px" height="400px">
    <s-card-title>RCON</s-card-title>
    <s-card-content>
      <s-scroller :autoscroll="autoscroll">
        <div v-for="(log, index) in logs" :key="index" class="msg">
          {{ log }}
        </div>
      </s-scroller>
    </s-card-content>
    <s-card-actions>
      <div class="actions">
        <div class="settings">
          <s-switch v-model="autoscroll">Autoscroll:</s-switch>
        </div>
        <div>
          <s-input
            type="text"
            placeholder="Command"
            v-model="command"
            v-on:keyup.enter="submit"
          />
          <s-button @click.stop="submit">Send</s-button>
        </div>
      </div>
    </s-card-actions>
  </s-card>
</template>

<style scoped>
.msg {
  margin-top: 3px;
  padding: 5px;
  border-radius: 3px;
  background: #f5f5f5;
}

.actions {
  width: 100%;
}

.actions > div {
  display: flex;
}

.actions > .settings {
  justify-content: flex-end;
  margin: 8px;
}
</style>

<script lang="ts">
import { defineComponent, onMounted, ref } from "vue";
import { Notify } from "./SNotification";

var backendURL = "ws://localhost:8080/api/rcon";

if (import.meta.env.MODE == "production") {
  backendURL = "wss://" + document.location.host + "/api/rcon";
}

export default defineComponent({
  name: "rcon",
  setup({}, { attrs }) {
    const logs = ref([]);
    const command = ref("");
    const autoscroll = ref(true);
    var conn;

    onMounted(() => {
      conn = new WebSocket(backendURL);
      conn.onerror = (event: ErrorEvent) => {
        Notify(
          `Could not establish WebSocket conneciton. Try reloading the page. ${JSON.stringify(
            event,
            ["message", "arguments", "type", "name"]
          )}`
        );
        conn = null;
      };
      conn.onclose = function (evt) {
        if (conn == null) {
          return;
        }
        Notify("Connection closed.");
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
        Notify("Need to specify command!");
        return;
      }
      conn.send(command.value);
      command.value = "";
    };
    return {
      logs,
      command,
      autoscroll,
      submit,
    };
  },
});
</script>
