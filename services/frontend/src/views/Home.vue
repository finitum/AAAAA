<template>
  <div class="flex flex-col w-100 items-center">
    <PackageTable />
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
