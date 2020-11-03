<template>
  <div class="inline-block shadow rounded-lg overflow-hidden mt-2 min-w-2/3">
    <table class="border-collapse min-w-full">
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
import { defineComponent, reactive, ref } from "vue";
import { Package } from "@/api/Models";
import { frequencyToDuration } from "@/api/Models";

async function fetchPackages(packages: Package[]) {
  const receivedPackage: Package = {
    Name: "test",
    RepoURL: "github.com/test/test",
    KeepLastN: 2,
    RepoBranch: "main",
    LastHash: "AAAAAA+/refs/main",
    UpdateFrequency: 3600000000000
  };

  packages.push(receivedPackage);
  packages.push(receivedPackage);
  packages.push(receivedPackage);
}

export default defineComponent({
  name: "PackageTable",
  setup() {
    const packages = reactive<Package[]>([]);
    const simple = ref(false);

    fetchPackages(packages);

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
