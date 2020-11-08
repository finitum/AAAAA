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
    <div class="bg-white z-20 text-center p-2 my-2 shadow-md rounded">
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
      <h3 v-if="mode === 'login'" class="text-2xl font-medium bg">Login</h3>
      <h3 v-if="mode === 'new-user'" class="text-2xl font-medium bg">
        New user
      </h3>
      <h3 v-if="mode === 'edit-password'" class="text-2xl font-medium bg">
        Edit password
      </h3>
      <form
        class="flex flex-col items-center align-middle justify-between text-center"
        @submit.prevent="ClickLogin"
      >
        <span class="label">
          <label for="username">Username:</label>
          <input
            v-model="user.Username"
            id="username"
            type="text"
            class="input-box"
            :disabled="mode === 'edit-password'"
          />
        </span>
        <span class="flex flex-col label">
          <label for="password">Password:</label>
          <input
            v-model="user.Password"
            id="password"
            type="password"
            class="input-box"
          />
        </span>
        <button v-if="mode === 'login'" class="button w-1/3" type="submit">
          Login
        </button>
        <button v-if="mode === 'new-user'" class="button w-1/3" type="submit">
          Add
        </button>
        <button
          v-if="mode === 'edit-password'"
          class="button w-1/3"
          type="submit"
        >
          Update
        </button>
      </form>
    </div>
  </div>
</template>

<script lang="ts">
import {
  defineComponent,
  onMounted,
  onUnmounted,
  reactive,
  PropType
} from "vue";
import { User } from "@/api/Models";
import { Login, NewUser, UpdateUser } from "@/api/API";
import { users } from "@/api/users";

export default defineComponent({
  name: "Login",
  props: {
    mode: {
      type: String as PropType<"login" | "new-user" | "edit-password">,
      default: "login"
    },

    originalUserName: {
      type: String
    }
  },
  setup(props, { emit }) {
    const user: User = reactive({ Username: "", Password: "" });

    onMounted(() => {
      if (
        props.mode === "edit-password" &&
        typeof props.originalUserName !== "undefined"
      ) {
        user.Username = props.originalUserName;
      }
    });

    function userExists(username: string): boolean {
      for (const u of users) {
        if (username === u.Username) {
          return true;
        }
      }
      return false;
    }

    function ClickLogin() {
      if (props.mode === "login") {
        Login(user).then(() => {
          emit("login");
          emit("close");
        });
      } else if (props.mode === "new-user") {
        if (userExists(user.Username)) {
          console.log("user exists");
          return true;
        }

        NewUser(user).then(() => {
          users.push(user);
          emit("login");
          emit("close");
        });
      } else if (props.mode === "edit-password") {
        UpdateUser(user).then(() => {
          for (let i = 0; i < users.length; i++) {
            if (users[i].Username == user.Username) {
              users[i] = user;
              break;
            }
          }

          emit("login");
          emit("close");
        });
      }
    }

    function escapeHandler(e: KeyboardEvent) {
      if (e.key === "Escape") {
        emit("close");
      }
    }

    onMounted(() => {
      window.addEventListener("keydown", escapeHandler);
    });

    onUnmounted(() => {
      window.removeEventListener("keydown", escapeHandler);
    });

    return {
      user,
      ClickLogin,
      userExists
    };
  }
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
