<template>
    <div class="relative h-[100dvh] pt-[4rem] max-w-[30rem] w-full mx-auto">
        <div style="height:4rem;" class="absolute border px-3 items-center justify-between gap-x-1 top-0 left-0 w-full p-2 flex">
            <div class="bg-slate-900 rounded-3xl px-2 text-sm font-medium py-0.5 text-white">room id: {{ roomID }}</div>

            <div class="italic text-slate-800">
                You
                <span v-if="activeMembers" id="members-list flex" >
                    <span v-for="(member, index) in activeMembers" :key="index">
                        , {{ member }}
                    </span>
                </span>
            </div>

        </div>
    <div ref="chatBox" id="chat-box" style="height:calc(100dvh - 8rem);">


  
        <div v-if="!messages.length" class="mx-auto my-5 text-sm italic text-center">No messages yet</div>

        <div v-else v-for="(msg, index) in messages" :key="index"
        class="block w-full"
        >
        <div class="max-w-[90%] w-fit" :class="msg.username == username ? 'ms-auto' : 'me-auto'">
        <div 
        :class="msg.username == username ? ' bg-cyan-500  text-white' : 'bg-gray-100 '"
        class="text-base break-words px-2 py-1  rounded-full mb-1 justify-start items-center gap-3 inline-flex">
        <!-- <h5 class="text-gray-900 text-sm font-normal leading-snug">Let me know</h5> -->
            {{ msg.text }}
    </div>
    <div v-if="index == messages.length -1 || messages[index+1].username !== msg.username" 
    class="mb-2 username italic text-blue mt-0.5 text-xs">
        {{ msg.username }}
    </div>
</div>
</div>
        <div v-if="typingUsers.length > 0" id="typing-indicator" class="mt-2" style="color: #555;">
            <div class="py-2 px-2  bg-gray-100 rounded-full w-fit">
                <div class='flex space-x-[0.2rem] justify-center items-center  h-3 dark:invert'>
                    <span class='sr-only'>Loading...</span>
                    <div class='h-1 w-1 bg-cyan-500 rounded-full animate-bounce [animation-delay:-0.3s]'></div>
                    <div class='h-1 w-1 bg-cyan-500 rounded-full animate-bounce [animation-delay:-0.15s]'></div>
                    <div class='h-1 w-1 bg-cyan-500 rounded-full animate-bounce'></div>
                </div>
            </div>
            <div class="mb-2 username italic text-cyan-500 mt-0.5 text-xs">
                {{ typingUsers.join(", ") }}
            </div>
        </div>
    </div>
    
    <div style="height:4rem;" class="absolute border px-3 items-center justify-between gap-x-1 bottom-0 left-0 w-full p-2 flex border-gray-200 bg-white">
    <!-- Return Button -->
    <button id="leave-btn" @click="leaveChat()" class="bg-slate-900 h-8 rounded-full px-3 py-2 flex items-center gap-2">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 16 16" fill="none">
            <path d="M10.7071 2.29289C10.3166 1.90237 9.68342 1.90237 9.29289 2.29289L5.29289 6.29289C4.90237 6.68342 4.90237 7.31658 5.29289 7.70711L9.29289 11.7071C9.68342 12.0976 10.3166 12.0976 10.7071 11.7071C11.0976 11.3166 11.0976 10.6834 10.7071 10.2929L7.70711 7.29289H13C13.5523 7.29289 14 6.84518 14 6.29289C14 5.74059 13.5523 5.29289 13 5.29289H7.70711L10.7071 2.29289Z" fill="white"/>
        </svg>
        <span class="text-white text-xs font-semibold">Return</span>
    </button>

    <!-- Message Input Section -->
     <form class="flex items-cenetr justify-center gap-x-1" @submit.prevent="sendMessage()">
        <input
            v-model="message"
            id="message"
            class="grow shrink basis-0 h-8 text-black text-xs font-medium leading-4 focus:outline-none bg-slate-200 px-3 py-2 rounded-full"
            type="text"
            ref="msgInput"
            placeholder="Type a message"
            @input="onMessageInput"
        />

    <!-- Send Button -->
    <button id="send-btn" class="bg-cyan-500 h-8 rounded-full px-3 py-2 flex items-center gap-2 shadow">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 16 16" fill="none">
            <path d="M9.04071 6.959L6.54227 9.45744M6.89902 10.0724L7.03391 10.3054C8.31034 12.5102 8.94855 13.6125 9.80584 13.5252C10.6631 13.4379 11.0659 12.2295 11.8715 9.81261L13.0272 6.34566C13.7631 4.13794 14.1311 3.03408 13.5484 2.45139C12.9657 1.8687 11.8618 2.23666 9.65409 2.97257L6.18714 4.12822C3.77029 4.93383 2.56187 5.33664 2.47454 6.19392C2.38721 7.0512 3.48957 7.68941 5.69431 8.96584L5.92731 9.10074C6.23326 9.27786 6.38623 9.36643 6.50978 9.48998C6.63333 9.61352 6.72189 9.7665 6.89902 10.0724Z" stroke="white" stroke-width="1.6" stroke-linecap="round" />
        </svg>
        <span class="text-white text-xs font-semibold leading-4">Send</span>
    </button>
    </form>
</div>


    </div>



</template>

<script setup>
import { ref, onMounted, onUnmounted } from "vue";
import {useRoute, useRouter} from "vue-router";

const msgInput = ref()
onMounted(()=>{
    msgInput.value.focus();
})
const route = useRoute()
const router = useRouter()
// Reactive variables
const chatBox = ref(null);
const username = localStorage.getItem('username')
const password = localStorage.getItem('password')
const roomID = route.params.id;

const message = ref("");
const messages = ref([]);
const activeMembers = ref([]);
const typingUsers = ref([]);

let typingTimeout = null;

// WebSocket setup
const protocol = location.protocol === "https:" ? "wss" : "ws";

const ws = new WebSocket(
    `${protocol}://${location.host}/ws?room_id=${roomID}&password=${password}&username=${username}`
);


// Create audio elements for different sound effects
const newMessageAudio = new Audio('/sounds/new-msg.mp3');
newMessageAudio.volume = 0.5;

const typingAudio = new Audio('/sounds/typing.mp3');
typingAudio.volume = 0.5;

const typingPlaying = ref(false)
const joinedAudio = new Audio('/sounds/joined.mp3');
joinedAudio.volume = 0.5;

// Initialize the WebSocket handler
ws.onmessage = (event) => {
    const msg = JSON.parse(event.data);
    console.log(msg);

    // Play new message audio
    if (msg.type === "message") {
        // Play new message sound once
        if(msg.username !== username){
            newMessageAudio.play();
        }

        // Add the message to the list
        messages.value.push({ username: msg.username, text: msg.text });

        // Scroll to the bottom of the chat
        chatBox.value.scrollTop = chatBox.value.scrollHeight;

        // Remove user from typing list after receiving the message
        typingUsers.value = typingUsers.value.filter((user) => user !== msg.username);

    } else if (msg.type === "active-members" && msg.active_members) {
        // Detect if there are new active members
        const newActiveMembers = new Set(msg.active_members);
        const oldActiveMembers = new Set(activeMembers.value);

        // Only play the 'joined' sound if there are new members who weren't in the previous list
        if (newActiveMembers.size > oldActiveMembers.size || [...newActiveMembers].some(user => !oldActiveMembers.has(user))) {
            joinedAudio.play();
        }

        // Update the active members list
        activeMembers.value = [...newActiveMembers].filter((user) => user !== username);

    } else if (msg.type === "typing" && msg.username) {
        // Start playing the typing sound if the user is typing
        if (!typingUsers.value.includes(msg.username) && msg.username !== username) {
            typingUsers.value.push(msg.username);
            if(!typingPlaying.value){
                // Play typing sound infinitely
                typingAudio.loop = true;
                typingPlaying.value = true;
                typingAudio.play();
            }
        }

    } else if (msg.type === "stop-typing" && msg.username) {
        // Stop playing the typing sound if the user stops typing
        typingUsers.value = typingUsers.value.filter((user) => user !== msg.username);
        // If no users are typing, stop the typing sound
        if (typingUsers.value.length === 0) {
            if(typingPlaying.value){
                typingAudio.pause();
                typingAudio.currentTime = 0; // Reset the sound to the beginning
                typingPlaying.value = false;

            }

        }
    }
};


const sendMessage = () => {
    if (username.trim() && message.value.trim()) {
    ws.send(
        JSON.stringify({
        type: "message",
        username: username.trim(),
        text: message.value.trim(),
        })
    );
    message.value = "";
    ws.send(
        JSON.stringify({ type: "stop-typing", username: username.trim() })
    );
    chatBox.value.scrollTop = chatBox.value.scrollHeight;
    msgInput.value.focus();
    }
};

const onMessageInput = () => {
    if (username.trim() && message.value.trim().length > 0) {
    ws.send(
        JSON.stringify({ type: "typing", username: username.trim() })
    );
    } else {
    ws.send(
        JSON.stringify({ type: "stop-typing", username: username.trim() })
    );
    }

    clearTimeout(typingTimeout);
    typingTimeout = setTimeout(() => {
    if (message.value.trim().length === 0) {
        ws.send(
        JSON.stringify({ type: "stop-typing", username: username.trim() })
        );
    }
    }, 1000);
};

const leaveChat = () => {
    // ws.send(JSON.stringify({ type: "leave", username: username.trim() }));
    localStorage.removeItem("roomID");
    router.back();
};

// Cleanup WebSocket on component unmount
onUnmounted(() => {
    ws.close();
});
</script>

<style scoped>
body {
    font-family: Arial, sans-serif;
    margin: 0;
    padding: 0;
}
#chat-box {
    height: 300px;
    border: 1px solid #ccc;
    overflow-y: scroll;
    padding: 10px;
}
.msg {
    margin: 5px 0;
}
.msg .username {
    font-weight: bold;
    color: #007BFF;
}
#input-box {
    display: flex;
    padding: 10px;
    border-top: 1px solid #ccc;
}

#message {
    flex: 1;
}
#send-btn {
    padding: 5px 15px;
}
#active-members {
    margin-top: 20px;
}
#active-members h3 {
    margin: 0;
    font-size: 18px;
}
#active-members ul {
    list-style: none;
    padding: 0;
}
#leave-btn {
    padding: 5px 15px;
    color: white;
    border: none;
    cursor: pointer;
}

/* HTML: <div class="loader"></div> */
.loader {
  width: 30px;
  height: 6px;
  --_g: no-repeat radial-gradient(farthest-side,#007BFF 94%,#007BFF);
  background:
    var(--_g) 25% 0,
    var(--_g) 75% 0;
  background-size: 6px 6px;
  position: relative;
  animation: l24-0 1s linear infinite;
}
.loader:before {
  content: "";
  position: absolute;
  height: 6px;
  aspect-ratio: 1;
  border-radius: 50%;
  background: #007BFF;
  inset: 0;
  margin: auto;
  animation: l24-1 1s cubic-bezier(0.5,300,0.5,-300) infinite;
}
@keyframes l24-0 {
  0%,24%  {background-position: 25% 0,75% 0}
  40%     {background-position: 25% 0,85% 0}
  50%,72% {background-position: 25% 0,75% 0}
  90%     {background-position: 15% 0,75% 0}
  100%    {background-position: 25% 0,75% 0}
}
@keyframes l24-1 {
  100% {transform:translate(0.1px)}
}
</style>
