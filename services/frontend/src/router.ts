import Homepage from "./views/homepage.vue"
import {createRouter, createWebHistory} from 'vue-router'

const router = createRouter({
    history: createWebHistory(process.env.BASE_URL),
    routes: [
        {
            path: "/",
            name: "homepage",
            component: Homepage,
        },
    ]
});

export default router;