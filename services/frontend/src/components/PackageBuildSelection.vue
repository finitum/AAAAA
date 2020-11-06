<template>
  <div
          class="fixed z-10 inset-0 overflow-y-auto flex justify-center items-center"
  >
    <div class="fixed inset-0 transition-opacity">
      <div
              class="absolute inset-0 bg-gray-500 opacity-75"
              @click="log()"
      ></div>
    </div>
    <div class="bg-white z-20 text-center p-2 my-2 shadow-md rounded content flex flex-col items-center">
      <h1>{{ package.Name }}</h1>

      <input value="{{pkg.Name}}" disabled>
      <input value="{{pkg.RepoURL}}" disabled>
      <input value="{{pkg.}}" disabled>

      hello world!
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType, ref } from "vue";
import {Result, ToPackage} from "@/api/AUR";

export default defineComponent({
  name: "PackageBuildSelection",

  props: {
    result: Object as PropType<Result>
  },

  setup(props) {
    const pkg = ref(ToPackage(props.result))

    return {
      pkg
    };
  }
});
</script>

<style lang="postcss" scoped>
.button {
  @apply flex-shrink-0 bg-secondary text-sm text-white py-2 px-3 rounded;
}
</style>
