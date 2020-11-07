<template>
  <Header @login="showLogin = true" />
  <router-view />
  <Login v-if="showLogin" @close="showLogin = false" />
  <Notification />
</template>

<script lang="ts">
import { defineComponent, ref } from "vue";
import Header from "@/components/Header.vue";
import Notification from "@/components/Notification.vue";
import Login from "@/components/Login.vue";

export default defineComponent({
  name: "App",
  components: { Header, Notification, Login },
  setup() {
    const showLogin = ref(false);

    return {
      showLogin
    };
  }
});
</script>

<style lang="postcss">
@tailwind base;
@tailwind components;
@tailwind utilities;

html {
  background: #f4f6f7;
}
</style>
