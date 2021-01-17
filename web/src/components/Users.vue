<template>
  <s-card width="600px" height="400px">
    <s-card-title>Users</s-card-title>
    <s-card-content>
      <s-scroller>
        <div class="line">
          <div class="headline username">Username</div>
          <div class="headline">Permissions</div>
          <div class="headline activated">Activated</div>
        </div>
        <div class="line" v-for="user in users" :key="user.username">
          <span class="username">{{ user.username }}</span>
          <s-select v-model="user.permissions" :items="permissions" @click.stop="saveUser(user)" placeholder="None"></s-select>
          <s-switch class="activated" v-model="user.activated" @click.stop="saveUser(user)"></s-switch>
        </div>
      </s-scroller>
    </s-card-content>
    <s-card-actions>
    </s-card-actions>
  </s-card>
</template>

<style scoped>
.line {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.headline {
  font-weight: 600;
  flex-grow: 1;
}

.username, .activated {
  width: 100px;
  flex: 0 0 auto;
}
</style>

<script lang="ts">
import {
  computed,
  defineComponent,
  onMounted,
  Ref,
  ref,
  watchEffect,
} from "vue";
import { http, httpPost } from "../utils";
import { User } from "../models";
import { Notify } from "./SNotification";

export default defineComponent({
  name: "Users",
  props: {},
  setup(props, ctx) {
    const users: Ref<User[]> = ref([]);

    const permissions = ref([
      "MODERATOR",
      "ADMIN",
    ]);

    onMounted(async () => {
      try {
        const res = await http<User[]>("/api/users");
        users.value = res.parsedBody as User[];
      } catch (err) {
        Notify("Failed to fetch users!" + err);
      }
    });

  const saveUser = async (user: User) => {
    try {
    const response = await httpPost(`/api/users/${user.username}`, JSON.stringify(user));
    } catch (err) {
      Notify(`Failed to save user: ${user.username}` + err);
    }
  };

    return {
      users,
      saveUser,
      permissions,
    };
  },
});
</script>
