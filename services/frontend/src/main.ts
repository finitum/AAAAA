import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";

// Only add one icon to the library
import { library } from "@fortawesome/fontawesome-svg-core";
import { faPen } from "@fortawesome/free-solid-svg-icons/faPen";
library.add(faPen);

// eslint-disable-next-line
createApp(App as any)
  .use(router)
  .mount("#app");
