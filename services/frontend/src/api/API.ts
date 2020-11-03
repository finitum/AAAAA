import axios from "axios";
import { Package, User } from "@/api/Models";
import { notificationState } from "@/components/NotificationState";

const client = axios.create({
  baseURL: "http://localhost:5000",
  headers: {
    "Content-type": "application/json"
  }
});

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

export async function GetPackages(): Promise<Package[]> {
  return client.get("/package").then(resp => resp.data);
}

export async function Login(user: User): Promise<string | null> {
  return client.post("/login", user).then(resp => {
    token = resp.data["token"];
    return token;
  });
}
