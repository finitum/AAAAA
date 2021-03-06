import { reactive } from "vue";
import { GetPackages } from "@/api/API";
import { Package } from "@/api/Models";

export const packages = reactive<Package[]>([]);

export async function loadPackages() {
  packages.splice(0, packages.length);
  packages.push(...(await GetPackages()));
}
