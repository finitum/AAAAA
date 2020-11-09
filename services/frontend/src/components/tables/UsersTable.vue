<template>
  <div class="inline-block shadow rounded-lg overflow-hidden mt-2 min-w-full">
    <table class="border-collapse min-w-full bg-white">
      <tr class="bg-secondary text-white">
        <th>Name</th>
        <th v-if="users.length > 0"></th>
      </tr>
      <tr v-for="user in users" v-bind:key="user.Name" class="row">
        <td>{{ user.Username }}</td>

        <td v-if="users.length > 0" class="lastcolthinner">
          <button @click="editUser = user" class="mr-3">
            <font-awesome-icon icon="pen" />
          </button>
          <button @click="deleteUser = user" class="mr-3">
            <font-awesome-icon class="text-red-600" icon="times" />
          </button>
        </td>
      </tr>

      <tr class="caption">
        <button class="w-full bg-white" @click="newUser = true">
          <font-awesome-icon class="text-green-600 text-xl" icon="plus" />
        </button>
      </tr>

      <tr v-if="users.length === 0" class="text-center caption py-4">
        There are no users yet.
      </tr>
    </table>

    <Login mode="new-user" v-if="newUser" @close="newUser = false" />
    <Login
      mode="edit-password"
      :original-user-name="editUser.Username"
      v-if="editUser !== null"
      @close="editUser = null"
    />

    <Dialog
      v-if="deleteUser !== null && users.length > 1"
      @close="deleteUser = null"
      mode="Confirm"
      @accept="doDeleteUser()"
    >
      <template v-slot:header> Delete user {{ deleteUser.Username }} </template>
      <template v-slot:message>
        Are you sure you want to remove user
        <span class="font-mono font-semibold">{{ deleteUser.Username }}</span
        >?
      </template>
    </Dialog>

    <Dialog
      v-if="deleteUser !== null && users.length <= 1"
      @close="deleteUser = null"
      mode="Info"
    >
      <template v-slot:header> Delete user {{ deleteUser.Username }} </template>
      <template v-slot:message>
        You can't delete the last user
      </template>
    </Dialog>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from "vue";
import { User } from "@/api/Models";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import Dialog from "@/components/modals/Dialog.vue";
import { DeleteUser } from "@/api/API";
import Login from "@/components/modals/Login.vue";
import { loadUsers, users } from "@/api/users";

export default defineComponent({
  name: "UsersTable",
  components: {
    FontAwesomeIcon,
    Dialog,
    Login
  },
  async setup() {
    await loadUsers();

    const newUser = ref(false);
    const editUser = ref(null);
    const deleteUser = ref<User | null>(null);

    function doDeleteUser() {
      if (deleteUser.value !== null) {
        const index = users.indexOf(deleteUser.value);

        DeleteUser(deleteUser.value.Username).then(() => {
          users.splice(index, 1);
          deleteUser.value = null;
        });
      }
    }

    return {
      users,
      deleteUser,
      editUser,
      doDeleteUser,
      newUser
    };
  }
});
</script>

<style scoped lang="postcss">
th,
td {
  @apply px-5 text-center border-collapse py-2 table-cell border-b-2 border-gray-100 border-opacity-25;
}
td {
  @apply bg-white;
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
