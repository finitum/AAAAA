<template>
  <div class="inline-block shadow rounded-lg overflow-hidden mt-2 min-w-full">
    <table class="border-collapse min-w-full bg-white">
      <tr class="bg-secondary text-white">
        <th>Name</th>
        <th>Repo url</th>
        <th v-if="!simple">Branch</th>
        <th v-if="!simple">Keep last</th>
        <th>Update frequency</th>
        <th v-if="loggedIn"></th>
      </tr>
      <tr v-for="pkg in packages" v-bind:key="pkg.Name" class="row">
        <td>{{ pkg.Name }}</td>
        <td>{{ pkg.RepoURL }}</td>
        <td v-if="!simple">{{ pkg.RepoBranch }}</td>
        <td v-if="!simple">{{ pkg.KeepLastN }}</td>
        <td>{{ frequencyToDuration(pkg.UpdateFrequency) }}</td>

        <td v-if="loggedIn" class="lastcolthinner">
          <button @click="editPackage = pkg" class="mr-3">
            <font-awesome-icon icon="pen" />
          </button>
          <button @click="deletePackage = pkg" class="mr-3">
            <font-awesome-icon class="text-red-600" icon="times" />
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

    <Dialog
      v-if="deletePackage !== null"
      @close="deletePackage = null"
      mode="Confirm"
      @accept="doDeletePackage()"
    >
      <template v-slot:header>
        Delete package {{ deletePackage.Name }}
      </template>
      <template v-slot:message>
        Are you sure you want to delete
        <span class="font-mono font-semibold">{{ deletePackage.Name }}</span
        >?
      </template>
    </Dialog>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from "vue";
import { frequencyToDuration, Package } from "@/api/Models";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import UpdatePackage from "@/components/UpdatePackage.vue";
import { DeletePackage, loggedIn } from "@/api/API";
import Dialog from "@/components/Dialog.vue";
import { loadPackages, packages } from "@/api/packages";

export default defineComponent({
  name: "PackageTable",
  components: {
    FontAwesomeIcon,
    UpdatePackage,
    Dialog
  },
  async setup() {
    const simple = ref(false);

    loadPackages();

    const editPackage = ref(null);
    const deletePackage = ref<Package | null>(null);

    function doDeletePackage() {
      if (deletePackage.value !== null) {
        const index = packages.indexOf(deletePackage.value);

        DeletePackage(deletePackage.value.Name).then(() => {
          packages.splice(index, 1);
          deletePackage.value = null;
        });
      }
    }

    return {
      simple,
      packages,
      frequencyToDuration,
      deletePackage,
      editPackage,
      loggedIn,
      doDeletePackage
    };
  }
});
</script>

<style scoped lang="postcss">
th,
td {
  @apply px-5 text-center border-collapse py-2 table-cell border-b-2 border-gray-100 border-opacity-25;
}

.lastcolthinner {
  width: 1%;
  white-space: nowrap;
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
