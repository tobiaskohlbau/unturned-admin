<template>
  <div class="container">
    <s-card>
      <s-card-title>Login</s-card-title>
      <template v-if="activated === false">
        <s-card-content>
          An administrator needs to activate your account. <br />
          Please try again later!
        </s-card-content>
        </template>
      <template v-else>
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
          <a href="/api/login/steam"
            ><img
              src="https://steamcdn-a.akamaihd.net/steamcommunity/public/images/steamworks_docs/english/sits_large_noborder.png"
          /></a>
        </s-card-actions>
      </template>
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
import { defineComponent, onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { Notify } from "../components/SNotification";
import { useStore } from "../store";
import { getToken, isActivated } from "../utils";
import { httpPost, http } from "../utils/http";
export default defineComponent({
  name: "Login",
  setup() {
    const router = useRouter();
    const store = useStore();
    const username = ref("");
    const password = ref("");
    const dummy = ref(false);
    const activated = ref(false);

    onMounted(async () => {
      activated.value = isActivated();
    })

    const login = async () => {
      try {
        const response = await httpPost(
          "/api/login",
          JSON.stringify({ username: username.value, password: password.value })
        );
        store.setToken(getToken());
        router.push({ name: "Home" });
      } catch (e) {
        Notify("failed to login");
      }
    };

    return {
      login,
      username,
      password,
      dummy,
      activated,
    };
  },
});
</script>
