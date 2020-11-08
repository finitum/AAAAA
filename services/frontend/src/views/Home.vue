<template>
  <div class="flex flex-col w-100 items-center align-middle">
    <div class="lg:w-2/3 w-full p-2">
      <Search v-show="loggedIn" />
      <Suspense>
        <template #default>
          <PackageTable />
        </template>
        <template #fallback>
          <div>Loading...</div>
        </template>
      </Suspense>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import PackageTable from "@/components/PackageTable.vue";
import Search from "@/components/Search.vue";
import { loggedIn } from "@/api/API";

export default defineComponent({
  name: "Home",
  components: { PackageTable, Search },
  setup() {

    return { loggedIn };
  }
});
</script>
