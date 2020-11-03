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
    <div class="bg-white z-20 text-center p-2 my-2 shadow-md rounded content">
      <div class="absolute cursor-pointer m-2" @click="$emit('close')">
        <svg
          class="fill-current text-black"
          xmlns="http://www.w3.org/2000/svg"
          width="18"
          height="18"
          viewBox="0 0 18 18"
        >
          <path
            d="M14.53 4.53l-1.06-1.06L9 7.94 4.53 3.47 3.47 4.53 7.94 9l-4.47 4.47 1.06 1.06L9 10.06l4.47 4.47 1.06-1.06L10.06 9z"
          ></path>
        </svg>
      </div>
      <h3 class="text-2xl font-medium bg">Login</h3>
      <div
        class="flex flex-col items-center align-middle justify-between text-center"
      >
        <span class="label">
          <label for="username">Username:</label>
          <input id="username" type="text" class="input-box" />
        </span>
        <span class="flex flex-col label">
          <label for="password">Password:</label>
          <input id="password" type="password" class="input-box" />
        </span>
        <button class="button w-1/3" @click="$emit('close')">Login</button>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  name: "Login"
});
</script>

<style scoped lang="postcss">
.label {
  @apply flex flex-col text-gray-700 font-bold p-2 m-1 w-full;
}

.button {
  @apply flex-shrink-0 bg-primary text-sm text-white py-2 px-3 rounded;
}

.button:hover {
  @apply bg-primarydark;
}

.input-box {
  @apply bg-gray-200 appearance-none border-2 border-gray-200 rounded py-2 px-4 text-gray-700 leading-tight;
}

.input-box:focus {
  @apply outline-none bg-white border-indigo-500;
}
</style>
