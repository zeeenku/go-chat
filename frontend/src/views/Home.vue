<script setup>
import {ref, onBeforeMount} from "vue";
import {useRouter} from "vue-router";
const router = useRouter()


const roomId = ref("")


const submit = async () => {
    localStorage.setItem('roomID', roomId.value);
    router.push({name: 'Room', params: { id: roomId.value } });
};

</script> 
<template>
<section class="bg-gray-50 dark:bg-gray-900">
  <div class="flex flex-col items-center justify-center px-6 py-8 mx-auto md:h-screen lg:py-0">
      <a href="#" class="flex items-center mb-6 text-2xl font-semibold text-gray-900 dark:text-white">
          <img class="w-9 h-9 rounded-lg mr-2" src="./../assets/images/logo.png" alt="logo">
          ZegoChat    
      </a>
      <div class="w-full bg-white rounded-lg shadow dark:border md:mt-0 sm:max-w-md xl:p-0 dark:bg-gray-800 dark:border-gray-700">
          <div class="p-6 space-y-4 md:space-y-6 sm:p-8">
              <h1 class="text-xl text-center font-bold leading-tight tracking-tight text-gray-900 md:text-2xl dark:text-white">
                Join a room
            </h1>

              <form @submit.prevent="submit()" class="space-y-4 md:space-y-6" action="#">
                  <div>
                      <label for="roomId" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">roomId</label>
                      <input v-model="roomId" type="roomId" name="roomId" id="roomId" class="bg-gray-50 border border-gray-300 text-gray-900 rounded-lg focus:ring-slate-900 focus:border-slate-900 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-400 dark:placeholder-gray-400 dark:text-white dark:focus:ring-slate-900 dark:focus:border-slate-900" placeholder="roomId" required="">
                  </div>


                  <button type="submit" class="w-full text-white bg-slate-900 hover:bg-slate-900 focus:ring-4 focus:outline-none focus:ring-slate-900 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-slate-900 dark:hover:bg-slate-900 dark:focus:ring-slate-900">Join room</button>

              </form>
          </div>
      </div>
  </div>
</section>
</template>