import { reactive } from "vue";

export const notificationState = reactive({
  enabled: false,
  message: "default message",
  color: "#fff"
});
