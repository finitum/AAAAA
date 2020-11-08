import axios from "axios";
import { Package, User } from "@/api/Models";
import { notificationState } from "@/components/NotificationState";
import { ref, watch } from "vue";

const client = axios.create({
  baseURL: "http://localhost:5000",
  headers: {
    "Content-type": "application/json"
  }
});

export const loggedIn = ref(false);
let token: string | null = null;

// Authorization interceptor
client.interceptors.request.use(conf => {
  if (token != null) {
    conf.headers.Authorization = `Bearer ${token}`;
  }

  return conf;
});

client.interceptors.response.use(undefined, error => {
  notificationState.message = error.message;
  notificationState.color = "#feb2b2";
  notificationState.enabled = true;

  console.log(JSON.parse(JSON.stringify(error)));

  return Promise.reject(error);
});

watch(loggedIn, value => {
  if (value && token !== null) {
    localStorage.setItem("token", token);
  }
});

export function CheckLoggedIn() {
  const result = localStorage.getItem("token");
  if (result !== null) {
    token = result;
    loggedIn.value = true;
  }
}

// Functions
export async function GetPackages(): Promise<Package[]> {
  return client.get("/package").then(resp => resp.data);
}

export async function Login(user: User): Promise<string | null> {
  return client.post("/login", user).then(resp => {
    token = resp.data["token"];
    loggedIn.value = true;
    return token;
  });
}

export async function NewUser(user: User): Promise<string | null> {
  return client.post("/user", user).then(resp => {
    token = resp.data["token"];
    loggedIn.value = true;
    return token;
  });
}

export function logOut() {
  token = null;
  loggedIn.value = false;

  localStorage.removeItem("token");
}

export async function GetAllUsers(
    localToken?: string
): Promise<User[]> {
  const originalToken = token;
  if (typeof localToken !== "undefined") {
    token = localToken;
  }

  if (token == null) {
    return Promise.reject("null token");
  }

  return client.get("/users").then(resp => {
    token = originalToken;
    return resp.data
  });
}

export async function DeleteUser(username: string, localToken?: string): Promise<void> {
  const originalToken = token;
  if (typeof localToken !== "undefined") {
    token = localToken;
  }

  if (token == null) {
    return Promise.reject("null token");
  }

  return client.delete("/user/" + username).then(() => {
    token = originalToken;
  });
}

export async function AddPackage(
  pkg: Package,
  localToken?: string
): Promise<void> {
  const originalToken = token;
  if (typeof localToken !== "undefined") {
    token = localToken;
  }

  if (token == null) {
    return Promise.reject("null token");
  }

  return client.post("/package", pkg).then(() => {
    token = originalToken;
  });
}

export async function UpdatePackage(
  pkg: Package,
  localToken?: string
): Promise<void> {
  const originalToken = token;
  if (typeof localToken !== "undefined") {
    token = localToken;
  }

  if (token == null) {
    return Promise.reject("null token");
  }

  return client.put("/package/" + pkg.Name, pkg).then(() => {
    token = originalToken;
  });
}

export async function DeletePackage(
  pkgname: string,
  localToken?: string
): Promise<void> {
  const originalToken = token;
  if (typeof localToken !== "undefined") {
    token = localToken;
  }

  if (token == null) {
    return Promise.reject("null token");
  }

  return client.delete("/package/" + pkgname).then(() => {
    token = originalToken;
  });
}
