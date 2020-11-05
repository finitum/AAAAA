<template>
  <header class="flex items-center justify-between flex-wrap bg-primary p-3 ">
    <h1 class="font-semibold text-white">AAAAA</h1>
    <button class="button" @click="$emit('login')">Login</button>
  </header>
</template>

<script lang="ts">
import { defineComponent } from "vue";
export default defineComponent({
  name: "Header"
});
</script>

<style lang="postcss" scoped>
.button {
  @apply flex-shrink-0 bg-secondary text-sm text-white py-2 px-3 rounded;
}
</style>
