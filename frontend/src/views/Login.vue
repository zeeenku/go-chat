<script setup>
import {ref, onBeforeMount} from "vue";
import {useRouter} from "vue-router";
const router = useRouter()

const error = ref("")

const username = ref(localStorage.getItem('username'))
const password = ref(localStorage.getItem('password'))
onBeforeMount(()=>{
    if(password.value && username.value){
        const roomID = localStorage.getItem('roomID')
        if(
            roomID
        )
        {
            router.push({name: 'Room', params: { id: roomID } });
            return;
        }

        router.push('/home')
        return ;
    }
})

const submit = async () => {

        const credentials = {
        username: username.value,
        password: password.value
    };


    try {
        // Send the request to the server
        const response = await fetch('/verify-login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(credentials)
        });

        // Parse the response as JSON
        const data = await response.json();

        if (response.ok && data.success) {
            localStorage.setItem('username', username);
            localStorage.setItem('password', password);
            router.push('/home')

        } else {
            error.value = data.error;
        }
    } catch (error) {
        // Catch network or other errors
        console.error('Network error:', error);
        alert('Network error, please try again later.');
    }
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
                Developed By Zenku
            </h1>
            <div v-if="error" class="text-sm w-11/12 mx-auto text-center font-medium text-red-600 dark:text-red-600">
                {{ error }}
            </div>
              <form @submit.prevent="submit()" class="space-y-4 md:space-y-6" action="#">
                  <div>
                      <label for="username" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Username</label>
                      <input v-model="username" type="username" name="username" id="username" class="bg-gray-50 border border-gray-300 text-gray-900 rounded-lg focus:ring-slate-900 focus:border-slate-900 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-400 dark:placeholder-gray-400 dark:text-white dark:focus:ring-slate-900 dark:focus:border-slate-900" placeholder="username" required="">
                  </div>
                  <div>
                      <label for="password" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Password</label>
                      <input v-model="password" type="text" name="password" id="password" placeholder="••••••••" class="bg-gray-50 border border-gray-300 text-gray-900 rounded-lg focus:ring-slate-900 focus:border-slate-900 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-400 dark:placeholder-gray-400 dark:text-white dark:focus:ring-slate-900 dark:focus:border-slate-900" required="">
                  </div>

                  <button type="submit" class="w-full text-white bg-slate-900 hover:bg-slate-900 focus:ring-4 focus:outline-none focus:ring-slate-900 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-slate-900 dark:hover:bg-slate-900 dark:focus:ring-slate-900">Sign in</button>

                  <p class="text-sm text-center text-slate-700 dark:text-slate-700">
                    Enter your last connection data to keep your conversations <br/> or Create new one
                </p>
              </form>
          </div>
      </div>
  </div>
</section>
</template>