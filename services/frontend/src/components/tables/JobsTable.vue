<template>
  <div class="inline-block shadow rounded-lg overflow-hidden mt-2 min-w-full">
    <table class="border-collapse min-w-full bg-white">
      <tr class="bg-secondary text-white">
        <th>Started</th>
        <th>Name</th>
        <th>Status</th>
      </tr>
      <tr v-for="job in jobs" v-bind:key="job.Name" class="row">
        <td>{{ new Date(job.Time).toLocaleString() }}</td>
        <td>{{ job.PackageName }}</td>
        <td>{{ BuildStatusToString(job.Status) }}</td>
      </tr>
      <tr v-if="jobs.length === 0" class="text-center caption py-4">
        There are no jobs currently scheduled.
      </tr>
    </table>
  </div>
</template>

<script lang="ts">
import { defineComponent, reactive } from "vue";
import { GetJobs } from "@/api/Jobs";
import { BuildStatusToString, Job } from "@/api/Models";

export default defineComponent({
  name: "JobsTable",
  async setup() {
    const jobs = reactive<Job[]>([]);

    jobs.push(...(await GetJobs()));

    return {
      jobs,
      BuildStatusToString
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
