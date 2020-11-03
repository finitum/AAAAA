<template>
  <div class="flex flex-col w-100 items-center">
    <Suspense>
      <template #default>
        <PackageTable />
      </template>
      <template #fallback>
        <div>Loading...</div>
      </template>
    </Suspense>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import PackageTable from "@/components/PackageTable.vue";

export default defineComponent({
  name: "Home",
  components: { PackageTable }
});
</script>
