<template>
  <div class="inline-block shadow rounded-lg overflow-hidden mt-2 min-w-full">
    <table class="border-collapse min-w-full bg-white">
      <tr class="bg-secondary text-white">
        <th>Name</th>
        <th>Repo url</th>
        <th v-if="!simple">Branch</th>
        <th v-if="!simple">Hash</th>
        <th v-if="!simple">Keep last</th>
        <th>Update frequency</th>
      </tr>
      <tr v-for="pkg in packages" v-bind:key="pkg.Name" class="row">
        <td>{{ pkg.Name }}</td>
        <td>{{ pkg.RepoURL }}</td>
        <td v-if="!simple">{{ pkg.RepoBranch }}</td>
        <td v-if="!simple">{{ pkg.LastHash }}</td>
        <td v-if="!simple">{{ pkg.KeepLastN }}</td>
        <td>{{ frequencyToDuration(pkg.UpdateFrequency) }}</td>
      </tr>
    </table>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from "vue";
import * as API from "@/api/API";
import { frequencyToDuration } from "@/api/Models";

export default defineComponent({
  name: "PackageTable",
  async setup() {
    const simple = ref(false);

    const packages = await API.GetPackages();

    return {
      simple,
      packages,
      frequencyToDuration
    };
  }
});
</script>

<style scoped lang="postcss">
th,
td {
  @apply px-5 text-center border-collapse py-2 table-cell border-b-2 border-gray-100 border-opacity-25;
}
</style>
