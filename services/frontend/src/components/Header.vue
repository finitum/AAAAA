<template>
  <header class="flex justify-between flex-wrap bg-primary items-stretch">
    <h1 class="font-semibold text-white m-3">AAAAA</h1>

    <div
      class="bg-primarylight flex flex-col items-center justify-center cursor-pointer"
      @click="loginButton()"
    >
      <span v-if="!loggedIn" class="button">Login</span>
      <span v-else class="button">Log out</span>
    </div>
  </header>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { loggedIn, logOut } from "@/api/API";

export default defineComponent({
  name: "Header",
  setup(_, { emit }) {
    function loginButton() {
      if (loggedIn.value) {
        logOut();
      } else {
        emit("login");
      }
    }

    return {
      loginButton,
      loggedIn
    };
  }
});
</script>

<style lang="postcss" scoped>
.button {
  @apply text-sm text-white px-5;
}
</style>
