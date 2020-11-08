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
        <th v-if="loggedIn"></th>
      </tr>
      <tr v-for="pkg in packages" v-bind:key="pkg.Name" class="row">
        <td>{{ pkg.Name }}</td>
        <td>{{ pkg.RepoURL }}</td>
        <td v-if="!simple">{{ pkg.RepoBranch }}</td>
        <td v-if="!simple">{{ pkg.LastHash }}</td>
        <td v-if="!simple">{{ pkg.KeepLastN }}</td>
        <td>{{ frequencyToDuration(pkg.UpdateFrequency) }}</td>

        <td v-if="loggedIn">
          <button @click="editPackage = pkg">
            <font-awesome-icon icon="pen" />
          </button>
        </td>
      </tr>
      <tr v-if="packages.length === 0" class="text-center caption py-4">
        There are no packages yet.
      </tr>
    </table>

    <UpdatePackage
      v-if="editPackage !== null"
      :pkgprop="editPackage"
      mode="update"
      @close="editPackage = null"
    ></UpdatePackage>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from "vue";
import * as API from "@/api/API";
import { frequencyToDuration } from "@/api/Models";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import UpdatePackage from "@/components/UpdatePackage.vue";
import {loggedIn} from '@/api/API';

export default defineComponent({
  name: "PackageTable",
  components: {
    FontAwesomeIcon,
    UpdatePackage
  },
  async setup() {
    const simple = ref(false);

    const packages = await API.GetPackages();

    const editPackage = ref(null);

    return {
      simple,
      packages,
      frequencyToDuration,
      editPackage,
      loggedIn
    };
  }
});
</script>

<style scoped lang="postcss">
th,
td {
  @apply px-5 text-center border-collapse py-2 table-cell border-b-2 border-gray-100 border-opacity-25;
}

.caption {
  /*
  Spans a tr the entire width of the table, without using colspan=0. With colspan 0 there seems to be a
  bug where text isn't centered. Only a positive non-zero integer allows text to be centered between columns,
  which doesn't work as well for us because the columns can change (depending on simple vs non-simple layout)
  */
  display: table-caption;
  caption-side: bottom;
}
</style>
