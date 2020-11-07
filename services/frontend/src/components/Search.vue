<template>
  <div>
    <input
      @focusin="showResults = true"
      @focusout="doFocusOut"
      @input="onInput"
      @keydown="onKeySearch"
      class="w-full"
      id="search"
      name="search"
      placeholder="search for packages here"
      type="text"
      v-model="term"
    />
    <div class="relative" v-show="showResults">
      <div
        class="absolute flex flex-col w-inherit bg-gray-200 rounded-b border-t-2 shadow-xl -mt-1 w-full py-1"
      >
        <button
          @mouseover="selected = result"
          class="flex flex-row w-full px-2 cursor-pointer block dropdown"
          v-bind:class="{ active: selected === result }"
          v-bind:key="result.ID"
          v-for="result of results"
          @click="addPackage"
        >
          <span class="font-bold mr-1 flex-none">
            {{ result.Name }}
          </span>
          <span
            class="opacity-50 min-w-0 overflow-hidden inline-block overflow-ellipsis whitespace-no-wrap mr-3"
          >
            {{ result.Description }}
          </span>
          <span class="ml-auto flex-none">
            {{ result.Version }}
          </span>
        </button>
      </div>
    </div>

    <PackageBuildSelection v-if="showPackageBuildSelection" :pkgprop="ToPackage(selected)" @close="showPackageBuildSelection=false" />
  </div>
</template>

<script lang="ts">
import { defineComponent, reactive, ref } from "vue";
import { ToPackage, Result, search, NewResult } from "@/api/AUR";
import PackageBuildSelection from "@/components/PackageBuildSelection.vue";

export default defineComponent({
  name: "Search",
  components: { PackageBuildSelection },
  setup() {
    const term = ref("");
    const results: Result[] = reactive([]);
    const selected = ref<Result>(NewResult());
    const showResults = ref(false);
    const showPackageBuildSelection = ref(false);

    function onKeySearch(event: KeyboardEvent) {
      let index = -1;

      for (let i = 0; i < results.length - 1; i++) {
        if (results[i].ID === selected.value.ID) {
          index = i;
          break;
        }
      }

      if (index === -1) {
        return
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

      selected.value = results[index];
    }

    function addPackage() {
      showPackageBuildSelection.value = true;
      showResults.value = false
    }

    function onInput() {
      search(term.value).then(resp => {
        results.splice(0, results.length);
        results.push(...resp);
      });
    }

    function doFocusOut(e: FocusEvent) {
      if (e.relatedTarget !== null && (e.relatedTarget as HTMLElement).classList.contains("dropdown")) {
          return
      }
      showResults.value = false
    }

    return {
      term,
      onInput,
      results,
      selected,
      onKeySearch,
      addPackage,
      doFocusOut,
      showResults,
      showPackageBuildSelection,
      ToPackage,
    };
  }
});
</script>

<style lang="postcss" scoped>
#search {
  @apply bg-gray-200 appearance-none border-2 border-gray-300 rounded py-2 px-4 text-gray-700 leading-tight;
}

.active {
  @apply bg-gray-400;
}
</style>
