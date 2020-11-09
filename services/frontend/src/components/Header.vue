<template>
  <header class="flex justify-between flex-wrap bg-primary items-stretch">
    <h1 class="font-semibold text-white m-3">AAAAA</h1>

    <div class="flex flex-row justify-center">
      <div
        class="cursor-pointer flex flex-col items-center justify-center h-full"
        @click="router.push('/')"
        v-bind:class="[
          router.currentRoute.value.path !== '/'
            ? 'bg-primary'
            : 'bg-primarylight'
        ]"
      >
        <span class="button">Home</span>
      </div>
      <div
        class="cursor-pointer flex flex-col items-center justify-center h-full"
        @click="router.push('/users')"
        v-bind:class="[
          router.currentRoute.value.path !== '/users'
            ? 'bg-primary'
            : 'bg-primarylight'
        ]"
        v-if="loggedIn"
      >
        <span class="button">Users</span>
      </div>

      <div
        class="bg-primary cursor-pointer flex flex-col items-center justify-center h-full"
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
        router.push("/");
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
