<template>
  <input
    id="search"
    v-model="term"
    @input="onInput"
    @keydown="onKeySearch"
    name="search"
    placeholder="yay"
    type="text"
    class="w-full"
  />
  <div class="relative">
    <div
      class="z-10 absolute flex flex-col w-inherit bg-gray-200 rounded-b border-t-2 shadow-xl -mt-1 w-full py-1"
    >
      <div
        v-for="result in results"
        v-bind:key="result.ID"
        class="flex flex-row w-full px-2"
        @mouseover="selected = result.ID"
        v-bind:class="{ active: selected === result.ID }"
      >
        <span class="font-bold mr-1 flex-none">{{ result.Name }}</span
        ><span
          class="opacity-50 min-w-0 overflow-hidden inline-block overflow-ellipsis whitespace-no-wrap mr-3"
          >{{ result.Description }}</span
        >
        <span class="ml-auto flex-none">{{ result.Version }}</span>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, reactive } from "vue";
import { Result, search } from "@/api/AUR";
export default defineComponent({
  name: "Search",
  setup() {
    const term = ref("");
    const results: Result[] = reactive([]);
    const selected = ref(-1);

    function onKeySearch(event: KeyboardEvent) {
      console.log(event);

      let index = -1;

      for (let i = 0; i < results.length - 1; i++) {
        if (results[i].ID === selected.value) {
          index = i;
          break;
        }
      }

      if (index === -1) {
        index = 0;
      }

      switch (event.key) {
        case "ArrowUp": {
          event.preventDefault();

          index--;
          if (index < 0) {
            index = results.length - 1;
          }
          break;
        }
        case "ArrowDown": {
          event.preventDefault();

          index++;

          if (index > results.length) {
            index = 0;
          }
          break;
        }
      }

      selected.value = results[index].ID;
    }

    function onInput() {
      search(term.value).then(resp => {
        results.splice(0, results.length);
        results.push(...resp);
      });
    }

    return {
      term,
      onInput,
      results,
      selected,
      onKeySearch
    };
  }
});
</script>

<style lang="postcss" scoped>
#search {
  @apply bg-gray-200 appearance-none border-2 border-gray-200 rounded py-2 px-4 text-gray-700 leading-tight;
}

.active {
  @apply bg-gray-400;
}
</style>
