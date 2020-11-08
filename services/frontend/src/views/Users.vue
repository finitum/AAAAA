<template>
  <div class="flex flex-col w-100 items-center align-middle">
    <div class="lg:w-2/3 w-full p-2">

      <Suspense>
        <template #default>
          <UsersTable />
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
import UsersTable from "@/components/UsersTable.vue";

export default defineComponent({
  name: "Users",
  components: { UsersTable },
  setup() {

    return { };
  }
});
</script>
