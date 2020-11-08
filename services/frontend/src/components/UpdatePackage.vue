<template>
  <div
    class="fixed z-10 inset-0 overflow-y-auto flex justify-center items-center"
  >
    <div class="fixed inset-0 transition-opacity">
      <div
        class="absolute inset-0 bg-gray-500 opacity-75"
        @click="$emit('close')"
      ></div>
    </div>
    <div
      class="bg-white z-20 text-center p-2 my-5 shadow-md rounded grid grid-cols-4 auto-cols-fr gap-8 lg:w-1/2 w-full"
    >
      <h1 class="col-span-full text-3xl font-semibold">{{ pkg.Name }}</h1>

      <label for="name" class="label">Package Name:</label>
      <input
        id="name"
        v-model="pkg.Name"
        :disabled="!externalPackage"
        class="input"
        required
      />

      <label for="repourl" class="label">Repository URL:</label>
      <input
        id="repourl"
        v-model="pkg.RepoURL"
        :disabled="!externalPackage"
        class="input"
        required
      />

      <label for="branch" class="label">Repository Branch:</label>
      <input
        id="branch"
        v-model="pkg.RepoBranch"
        :disabled="!externalPackage"
        class="input"
        required
      />

      <label for="duration" class="label">Update frequency: </label>
      <DurationPicker
        id="duration"
        class="col-span-3"
        v-model="pkg.UpdateFrequency"
      />

      <label for="lastn" class="label">Keep last:</label>
      <input
        id="lastn"
        type="number"
        v-model="pkg.KeepLastN"
        @focusout="focusOut"
        class="input"
        min="1"
        @focus="$event.target.select()"
        required
      />

      <button class="col-span-full" @click="addPackage" v-if="mode==='add'">Add Package</button>
      <button class="col-span-full" @click="updatePackage" v-else>Update Package</button>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType, ref, onMounted } from "vue";
import DurationPicker from "@/components/DurationPicker.vue";
import { NewPackage, Package } from "@/api/Models";
import { AddPackage, UpdatePackage } from "@/api/API";

export default defineComponent({
  name: "UpdatePackage",
  components: { DurationPicker },

  props: {
    pkgprop: {
      type: Object as PropType<Package>,
      required: true
    },
    mode: {
      type: Object as PropType<"update" | "add">,
      required: true
    }
  },

  setup(props, { emit }) {
    const pkg = ref(NewPackage());
    const externalPackage = ref(false);

    onMounted(() => {
      pkg.value = props.pkgprop;
    });

    function focusOut(e: Event) {
      if ((e.target as HTMLInputElement).value === "") {
        (e.target as HTMLInputElement).value = "1";
      }
    }

    function addPackage() {
      AddPackage(pkg.value).then(() => emit("close"));
    }

    function updatePackage() {
      UpdatePackage(pkg.value).then(() => emit("close"));
    }

    return {
      pkg,
      focusOut,
      externalPackage,
      addPackage,
      updatePackage,
    };
  }
});
</script>

<style lang="postcss" scoped>
button {
  @apply flex-shrink-0 bg-secondary text-sm text-white py-2 px-3 rounded;
}

.input {
  @apply bg-gray-200 appearance-none border-2 border-gray-200 rounded py-2 px-4 text-gray-700 leading-tight col-span-3;

  &:focus {
    @apply outline-none bg-white border-secondary;
  }
}

.label {
  @apply col-span-1 text-center block pt-1;
}
</style>
