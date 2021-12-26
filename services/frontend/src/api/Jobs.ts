import axios from "axios";
import { notificationState } from "@/components/modals/NotificationState";
import { Job } from "@/api/Models";

const client = axios.create({
  baseURL: " http://localhost:5002/",
  headers: {
    "Content-type": "application/json"
  }
});

client.interceptors.response.use(undefined, error => {
  notificationState.message = error.message;
  notificationState.color = "#feb2b2";
  notificationState.enabled = true;

  return Promise.reject(error);
});

export async function GetJobs(): Promise<Job[]> {
  return client.get("/jobs").then(resp => {
    if (typeof resp !== "undefined") {
      return resp.data;
    } else {
      return [];
    }
  });
}
