<template>
    <div class="relative h-screen pt-[4rem] max-w-[30rem] w-full mx-auto">
        <div style="height:4rem;" class="absolute border px-3 items-center justify-between gap-x-1 top-0 left-0 w-full p-2 flex">
            <div class="italic text-slate-800">
                You, 
                <span v-if="activeMembers" id="members-list flex" >
                    <span v-for="(member, index) in activeMembers" :key="index">
                        {{ member }}
                        <span v-if="index<activeMembers.length-1">,</span>
                    </span>
                </span>
            </div>

            <div class="bg-slate-900 rounded-md px-1 py-0.5 text-white">id: {{ roomID }}</div>
        </div>
    <div ref="chatBox" id="chat-box" style="height:calc(100vh - 9rem);">


  
        <div v-if="!messages.length" class="mx-auto text-center">No messages yet</div>

        <div v-else v-for="(msg, index) in messages" :key="index">
        <div 
        :class="msg.username == username ? 'bg-blue text-white' : 'bg-gray-100 '"
        class="text-sm mt-1 break-words px-2 py-1 rounded justify-start items-center gap-3 inline-flex">
        <!-- <h5 class="text-gray-900 text-sm font-normal leading-snug">Let me know</h5> -->
            {{ msg.text }}
    </div>
    <div v-if="index == messages.length -1 || messages[index+1].username !== msg.username" class="username italic text-blue mt-0.5 text-xs">
        {{ msg.username }}
    </div>

</div>
        <div v-if="typingUsers.length > 0" id="typing-indicator" style="color: #555;">
        {{ typingUsers.join(", ") }} typing...
        </div>
    </div>
    <div style="height:5rem;" class="absolute border px-3 items-center justify-between gap-x-1 bottom-0 left-0 w-full p-2 flex">
        <button id="leave-btn" @click="leaveChat()" class="bg-slate-900 rounded-md">Return </button>


        <input
        v-model="message"
        id="message"
        class=" bg-slate-200 px-2 py-0.5 rounded-md h-[2.2rem]"
        type="text"
        placeholder="Type a message"
        @input="onMessageInput"
        />

        <button id="send-btn" class="bg-slate-900 rounded-md px-1 text-white" @click="sendMessage">Send</button>

    </div>

    </div>



</template>

<script setup>
import { ref, onMounted, onUnmounted } from "vue";
import {useRoute, useRouter} from "vue-router";

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
console.log(    `${protocol}://${location.host}/ws?room_id=${roomID}&password=${password}&username=${username}`)
ws.onmessage = (event) => {
    const msg = JSON.parse(event.data);
    console.log(msg)

    if (msg.type === "message") {
    messages.value.push({ username: msg.username, text: msg.text });
    chatBox.value.scrollTop = chatBox.value.scrollHeight;
    typingUsers.value = typingUsers.value.filter((user) => user !== msg.username);
    } else if (msg.type === "active-members" && msg.active_members) {
    activeMembers.value = msg.active_members.filter((user) => user !== username);
    } else if (msg.type === "typing" && msg.username) {
    if (!typingUsers.value.includes(msg.username) && msg.username !== username) {
        typingUsers.value.push(msg.username);
    }
    } else if (msg.type === "stop-typing" && msg.username) {
    typingUsers.value = typingUsers.value.filter((user) => user !== msg.username);
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
</style>
