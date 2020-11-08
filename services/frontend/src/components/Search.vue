<template>
  <div>
    <input
      @focusin="doFocusIn"
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
          @click="addPackage"
          @mouseover="selected = result"
          class="flex flex-row w-full px-2 cursor-pointer block dropdown"
          v-bind:class="{ active: selected === result }"
          v-bind:key="result.ID"
          v-for="result of results"
        >
          <span v-if="result.ID !== -1" class="font-bold mr-1 flex-none">
            {{ result.Name }}
          </span>
          <span
            v-if="result.ID !== -1"
            class="opacity-50 min-w-0 overflow-hidden inline-block overflow-ellipsis whitespace-no-wrap mr-3"
          >
            {{ result.Description }}
          </span>
          <span v-if="result.ID !== -1" class="ml-auto flex-none">
            {{ result.Version }}
          </span>

          <span v-if="result.ID === -1" class="mr-auto">
            Add package by url:
            {{ result.URL }}
          </span>
        </button>
      </div>
    </div>

    <UpdatePackage
      :pkgprop="ToPackage(selected, selected.ID === -1)"
      @close="showPackageBuildSelection = false"
      :external="selected.ID === -1"
      mode="add"
      v-if="showPackageBuildSelection"
    />
  </div>
</template>

<script lang="ts">
import { defineComponent, reactive, ref } from "vue";
import { NewResult, Result, search, ToPackage } from "@/api/AUR";
import UpdatePackage from "@/components/UpdatePackage.vue";
import { packages } from "@/api/packages";

const gitUrlRegex = /(?:(?:git|ssh|https?)|(?:git@[\w.]+))(?::(\/\/)?)([\w.@:/\-~]+)(?:\.git)(\/)?/;

export default defineComponent({
  name: "Search",
  components: { UpdatePackage },
  setup() {
    const term = ref("");
    const results: Result[] = reactive([]);
    const selected = ref<Result>(NewResult());
    const showResults = ref(false);
    const showPackageBuildSelection = ref(false);

    function addPackage() {
      showPackageBuildSelection.value = true;
      showResults.value = false;
    }

    function onKeySearch(event: KeyboardEvent) {
      let index = -1;

      for (let i = 0; i < results.length - 1; i++) {
        if (results[i].ID === selected.value.ID) {
          index = i;
          break;
        }
      }

      if (index === -1 && results.length == 0) {
        return;
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
        case "Enter": {
          addPackage();
          break;
        }
        case "Escape": {
          if (event.target !== null && event.target instanceof HTMLElement) {
            event.target.blur();
          }
          break;
        }
      }

      if (!(index > results.length || index < 0)) {
        selected.value = results[index];
      }
    }

    function onInput() {
      const regexmatch = gitUrlRegex.exec(term.value);
      if (regexmatch !== null) {
        const res = NewResult();

        res.ID = -1;
        res.URL = term.value;

        const urlpath = regexmatch[2].split("/");
        res.Name = urlpath[urlpath.length - 1];

        results.splice(0, results.length);
        results.push(res);
      } else if (term.value.includes("/")) {
        results.splice(0, results.length);
        return;
      } else {
        search(term.value).then(resp => {
          results.splice(0, results.length);

          for (const res of resp) {
            let found = false;
            for (const pkg of packages) {
              if (pkg.Name === res.Name) {
                found = true;
                break;
              }
            }

            if (!found) {
              results.push(res);
            }
          }
        });
      }
    }

    function doFocusOut(e: FocusEvent) {
      if (
        e.relatedTarget !== null &&
        (e.relatedTarget as HTMLElement).classList.contains("dropdown")
      ) {
        return;
      }
      showResults.value = false;
      results.splice(0, results.length);
    }

    function doFocusIn() {
      showResults.value = true;
      onInput();
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
      doFocusIn
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
