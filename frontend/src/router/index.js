import { createRouter, createMemoryHistory } from "vue-router";
import Login from '../views/Login.vue';
const routes = [
{
    path: "/",
    component: Login,
},
];
const router = createRouter({
  history: createMemoryHistory(),
  routes,
});

export default router;