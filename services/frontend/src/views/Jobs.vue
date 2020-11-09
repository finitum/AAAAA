<template>
  <div class="flex flex-col w-100 items-center align-middle">
    <div class="lg:w-1/4 w-full p-2">
      <Suspense>
        <template #default>
          <JobsTable />
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
import JobsTable from "@/components/tables/JobsTable.vue";

export default defineComponent({
  name: "Jobs",
  components: { JobsTable },
  setup() {
    return {};
  }
});
</script>
