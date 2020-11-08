import { reactive } from "vue";
import * as API from "@/api/API";
import { User } from "@/api/Models";

export const users = reactive<User[]>([]);

export async function loadUsers() {
  users.splice(0, users.length);
  users.push(...(await API.GetAllUsers()));
}
