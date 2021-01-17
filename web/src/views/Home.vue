<template>
  <div class="site">
    <div class="sidebar">
      <p>Welcome {{ store.state.token.username }}!</p>
      <s-button @click.stop="logout">Logout</s-button>
    </div>
    <div :class="['container', className]">
      <rcon v-if="hasPermission('MODERATOR')"></rcon>
      <manager v-if="hasPermission('MODERATOR')"></manager>
      <unturned-map></unturned-map>
      <users v-if="hasPermission('ADMIN')"></users>
    </div>
  </div>
</template>

<style scoped lang="postcss">
.site {
  display: flex;
}

.side {
  min-width: 200px;
}

.container {
  flex-grow: 1;
  width: 100%;
  display: flex;
  flex-wrap: wrap;
  justify-content: space-around;
}

.flex-gap {
  column-gap: 24px;
  row-gap: 24px;
}

.flex-no-gap {
  margin-top: -24px;
  margin-right: -24px;

  & > * {
    margin-top: 24px;
    margin-right: 24px;
  }
}
</style>

<script>
import RCON from "../components/RCON.vue";
import Manager from "../components/Manager.vue";
import UnturnedMap from "../components/UnturnedMap.vue";
import Users from "../components/Users.vue";
import { defineComponent, ref } from "vue";
import { useRouter } from "vue-router";
import { useStore } from "../store";
import { checkFlexGap } from "../utils";
import { hasPermission, isAuthenticated } from '../utils';

export default defineComponent({
  name: "Home",
  components: {
    rcon: RCON,
    manager: Manager,
    'unturned-map': UnturnedMap,
    users: Users,
  },
  setup() {
    const className = ref("flex-gap");
    const router = useRouter();
    const store = useStore();

    // Workaround to support safari.
    // Drop if safari supports gap on flex items,
    // currently in technical preview.
    if (!checkFlexGap()) {
      className.value = "flex-no-gap";
    }

    const logout = () => {
      document.cookie = "tid=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
      router.push('/login');
    };

    return {
      className,
      store,
      logout,
      hasPermission,
    };
  },
});
</script>
