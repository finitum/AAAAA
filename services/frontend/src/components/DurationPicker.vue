<template>
  <div class="flex flex-row items-center justify-between w-100">
    <label for="d">
      <span>days:</span>
      <input
        id="d"
        type="number"
        v-model="day"
        @change="update()"
        min="0"
        @focus="$event.target.select()"
        required
      />
    </label>
    <label for="h">
      <span>hours:</span>
      <input
        id="h"
        type="number"
        v-model="hour"
        @change="update()"
        min="0"
        @focus="$event.target.select()"
        required
      />
    </label>
    <label for="m">
      <span>minutes:</span>
      <input
        id="m"
        type="number"
        v-model="min"
        @change="update()"
        min="0"
        @focus="$event.target.select()"
        required
      />
    </label>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, watch } from "vue";

export default defineComponent({
  name: "DurationPicker",
  props: ["modelValue"],

  setup(props, { emit }) {
    const min = ref(0);
    const hour = ref(0);
    const day = ref(0);

    function normalizeTime() {
      if (min.value >= 60) {
        hour.value += Math.floor(min.value / 60);
        min.value %= 60;
      }

      if (hour.value >= 24) {
        day.value += Math.floor(hour.value / 24);
        hour.value %= 24;
      }
    }

    function timeAsSeconds(): number {
      let seconds = 0;

      seconds += min.value * 60;
      seconds += hour.value * 60 * 60;
      seconds += day.value * 60 * 60 * 24;

      return seconds;
    }

    watch(
      () => props.modelValue,
      (value, prev) => {
        const seconds = props.modelValue / 1000 / 1000 / 1000;
        if (timeAsSeconds() !== seconds) {
          min.value = Math.floor(seconds / 60);
          normalizeTime();
        }
      }
    );

    function update() {
      if (((min.value as unknown) as string) === "") {
        min.value = 0;
      }
      if (((hour.value as unknown) as string) === "") {
        hour.value = 0;
      }
      if (((day.value as unknown) as string) === "") {
        day.value = 0;
      }

      normalizeTime();

      const seconds = timeAsSeconds();
      emit("update:modelValue", seconds * 1000 * 1000 * 1000);
    }

    return {
      focusOut: update,
      min,
      day,
      hour
    };
  }
});
</script>

<style lang="postcss" scoped>
input {
  @apply w-24 bg-gray-200 appearance-none border-2 border-gray-200 rounded py-2 px-4 text-gray-700 leading-tight;

  &:focus {
    @apply outline-none bg-white border-secondary;
  }
}

label span {
  @apply mr-3;
}
</style>
