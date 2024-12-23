<script setup>
import {useRouter} from "vue-router";
import {ref, onMounted} from "vue";
const router = useRouter();

const goToRoom = (id) => {
    
}
const logout = ()=>{
    localStorage.clear();
    router.push('/');
}
const rooms = ref([])

class RoomCard{
    constructor(data){
        this.msg = data.msg;
        this.msgAuthor = data.msgAuthor;
        this.id = data.id;
        this.name = data.name;
        this.nonReadCounter = data.nonReadCounter;
    }
}

const eventSource = ref(null)
onMounted(()=>{

    // do the fetch
    const data = [{
        msg : "hello world",
        msgAuthor: "zenku",
        id : 1,
        name : "chat room 1",
        nonReadCounter: 4
    },

    {
        msg : "how are you bro",
        msgAuthor: "zenddku",
        id : 1,
        name : "chat room 2",
        nonReadCounter: 4
    },

    {
        msg : "kikokok",
        msgAuthor: "zzz",
        id : 1,
        name : "chat room 9",
        nonReadCounter: 4
    },
]

    rooms.value = data.map(el=> new RoomCard(el))
    eventSource.value = new EventSource('/events');
    eventSource.onmessage = function(event) {

    const ev = event.data;
    rooms.value.map((el)=>{
        if(el.id == ev.id){
            el.nonReadCounter += 1;
            el.msgAuthor = ev.msgAuthor;
            el.msg = ev.msg;
        }
    })
};

})




</script>
<template>

    <header class="px-4 h-14 shadow py-2 flex justify-between items-center">
        <a href="#" class="flex items-center justify-center">
          <img class="w-9 h-9 rounded-lg mr-2" src="./../assets/images/logo.png" alt="logo">
          ZegoChat    
      </a>
        <button @click="logout" class="text-white font-medium text-base px-2 py-0.5 bg-slate-900 rounded-md">logout</button>
    </header>

    <div style="height:calc(100vh -  3.5rem);" class="max-w-80 w-full py-10  overflow-y-auto mx-auto">
        <div v-if="!rooms.length">You havent joined any rooms yet</div>
        <div v-else v-for="el in rooms" :key="el" >
            <div @click="goToRoom(el.id)" class="w-full cursor-pointer hover:bg-slate-100 p-3 bg-slate-50 border">
            <div class="flex items-center justify-between">
                <div class="p-2 rounded-lg me-1 bg-slate-900 text-white">
                        <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                            <path d="M7.45648 3.08984C4.21754 4.74468 2 8.1136 2 12.0004C2 13.6001 2.37562 15.1121 3.04346 16.4529C3.22094 16.8092 3.28001 17.2165 3.17712 17.6011L2.58151 19.8271C2.32295 20.7934 3.20701 21.6775 4.17335 21.4189L6.39939 20.8233C6.78393 20.7204 7.19121 20.7795 7.54753 20.957C8.88836 21.6248 10.4003 22.0005 12 22.0005C16.8853 22.0005 20.9524 18.4973 21.8263 13.866C20.1758 15.7851 17.7298 17.0004 15 17.0004C10.0294 17.0004 6 12.971 6 8.00045C6 6.18869 6.53534 4.50197 7.45648 3.08984Z" fill="white"/>
                            <path d="M21.8263 13.8655C21.9403 13.2611 22 12.6375 22 12C22 6.47715 17.5228 2 12 2C10.4467 2 8.97611 2.35415 7.66459 2.98611C7.59476 3.01975 7.52539 3.05419 7.45648 3.08939C6.53534 4.50152 6 6.18824 6 8C6 12.9706 10.0294 17 15 17C17.7298 17 20.1758 15.7847 21.8263 13.8655Z" fill="white"/>
                            </svg>
                    </div>
                <div class="w-9/12">

                    <h2 class="text-blue text-sm  font-semibold">#{{ el.name }}</h2>
            <div class="flex justify-between items-center">
                <div class="text-xs">
                    <span class="text-blue font-medium italic">{{el.msgAuthor}}: </span>
                    <span class="font-semibold">{{ el.msg }}</span>

                </div>
            </div>
                </div>

                <div v-if="el.nonReadCounter" class="w-5 text-xs font-bold h-5 flex items-center justify-center rounded-full bg-slate-900 text-white">{{ el.nonReadCounter }}</div>
                <div v-else></div>
            </div>
           
            </div>
        </div>
    </div>
</template>