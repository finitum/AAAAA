<template>
  <Header />
  <router-view />
</template>

<script lang="ts">
import { defineComponent } from "vue";
import Header from "@/components/Header.vue";

export default defineComponent({
  name: "App",
  components: { Header }
});
</script>

<style lang="postcss">
@tailwind base;
@tailwind components;
@tailwind utilities;
</style>
