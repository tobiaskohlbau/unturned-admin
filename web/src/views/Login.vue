<template>
  <div class="container">
    <s-card>
      <s-card-title>Login</s-card-title>
      <s-card-content>
        <s-input
          type="text"
          placeholder="Username"
          v-model="username"
          @keyup.enter="login"
        />
        <s-input
          type="password"
          placeholder="Password"
          v-model="password"
          @keyup.enter="login"
        />
      </s-card-content>
      <s-card-actions>
        <s-button @click.stop="login">Login</s-button>
      </s-card-actions>
    </s-card>
  </div>
</template>

<style scoped>
.container {
  width: 100%;
  display: flex;
  flex-wrap: wrap;
  justify-content: space-around;
  row-gap: 24px;
}
</style>

<script lang="ts">
import { defineComponent, ref } from "vue";
import { useRouter } from "vue-router";
import { Notify } from "../components/SNotification";
import { httpPost } from "../utils/http";
export default defineComponent({
  name: "Login",
  setup() {
    const router = useRouter();
    const username = ref("");
    const password = ref("");

    const login = async () => {
      try {
        const response = await httpPost(
          "/api/login",
          JSON.stringify({ username: username.value, password: password.value })
        );
        router.push({ name: "Home" });
      } catch (e) {
        Notify("failed to login");
      }
    };

    return {
      login,
      username,
      password,
    };
  },
});
</script>
