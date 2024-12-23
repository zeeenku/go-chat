import { createRouter, createMemoryHistory } from "vue-router";
import Login from '../views/Login.vue';
import Home from '../views/Home.vue';
import Room from '../views/Room.vue';

const routes = [
{
    path: "/",
    component: Login,
},




{
    path: "/home",
    component: Home,
},


{
    path: '/room/:id',
    name: 'Room',
    component: Room, // Replace with your component
  },

];
const router = createRouter({
  history: createMemoryHistory(),
  routes,
});

export default router;