import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import { CheckLoggedIn } from "@/api/API";

CheckLoggedIn();

// Only add one icon to the library
import { library } from "@fortawesome/fontawesome-svg-core";
import { faPen } from "@fortawesome/free-solid-svg-icons/faPen";
import { faTimes } from "@fortawesome/free-solid-svg-icons/faTimes";
library.add(faPen);
library.add(faTimes);

// eslint-disable-next-line
createApp(App as any)
  .use(router)
  .mount("#app");
