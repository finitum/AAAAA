<template>
  <header class="flex justify-between flex-wrap bg-primary items-stretch">
    <h1 class="font-semibold text-white m-3">AAAAA</h1>

    <div class="flex flex-row justify-center">
      <div
              class="bg-primarylight cursor-pointer flex flex-col items-center justify-center h-full"
              @click="router.push('/')"
              v-if="router.currentRoute.value.path !== '/'"
      >
        <span class="button">Home</span>
      </div>
      <div
              class="bg-primarylight cursor-pointer flex flex-col items-center justify-center h-full"
              @click="router.push('/users')"
              v-if="loggedIn && router.currentRoute.value.path !== '/users'"
      >
        <span class="button">Users</span>
      </div>

      <div
              class="bg-primarylight cursor-pointer flex flex-col items-center justify-center h-full"
              @click="loginButton()"
      >
        <span v-if="!loggedIn" class="button">Login</span>
        <span v-else class="button">Log out</span>
      </div>
    </div>

  </header>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { loggedIn, logOut } from "@/api/API";
import router from "@/router";

export default defineComponent({
  name: "Header",
  setup(_, { emit }) {
    function loginButton() {
      if (loggedIn.value) {
        router.push("/")
        logOut();
      } else {
        emit("login");
      }
    }

    return {
      loginButton,
      loggedIn,
      router
    };
  }
});
</script>

<style lang="postcss" scoped>
.button {
  @apply text-sm text-white px-5;
}
</style>
