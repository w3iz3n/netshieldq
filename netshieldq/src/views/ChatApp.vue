<template>
  <div class="chat-app">
    <NavigationBar/>
    <FriendList :friends="friends" @select-friend="openChat" :friendMessages="messages"/>
    <ChatWindow class="chat-window" :messages="messages" @send-message="sendMessagePostoBack" :currentFriend="currentFriend"/>
  </div>
</template>

<script>
import ChatWindow from './ChatWindow.vue';
import NavigationBar from './NavigationBar.vue';
import FriendList from './FriendList.vue';
import Module from "../../public/js/test_kyber512";
import axios from "axios";

export default {
  components: {
    ChatWindow,
    NavigationBar,
    FriendList
  },
  data() {
    return {
      name: 'ChatApp',
      messages: [],
      friends: [],
      currentFriend: null,
      friendMessages: {}
    };
  },
  methods: {
    async sendMessagePostoBack(message) {
      try {
        const jwt = localStorage.getItem('jwt');
        let timestamp = new Date().toISOString(); // "2021-01-02T15:04:05.000Z"
        const response = await fetch('/api/message/send', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${jwt}`
          },
          body: JSON.stringify({
            id: Date.now(),
            senderId: Number(message.senderId),
            receiverId: Number(message.receiverId),
            text: message.text,
            timestamp: timestamp,
            isHistorical: message.isHistorical,
            messageType: message.messageType,
          })
        });

        if (!response.ok) {
          const errorDetails = await response.text();
          throw new Error(`Failed to send message: ${response.status} ${errorDetails}`);
        }

        const data = await response.json();
        console.log('Message sent successfully:', data);
      } catch (error) {
        console.error('Error sending message:', error);
      }
    },



    openChat(friend) {
      this.currentFriend = friend;
      this.$store.dispatch('initchat', this.currentFriend.id);
      this.messages = this.friendMessages[friend["id"]] || [];
    },
    async getFriends() {
      try {
        const jwt = localStorage.getItem('jwt');
        const response = await fetch('/api/getfriends', {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${jwt}`
          }
        });

        if (!response.ok) {
          const errorDetails = await response.text();
          throw new Error(`Failed to fetch friends: ${response.status} ${errorDetails}`);
        }

        let data = await response.json();
        this.friends = data;
        // 检查每个好友是否已经有共享密钥文件
        for (const friend of this.friends) {
          const sharedSecretFileName = `ssA_${localStorage.getItem('userid')}_${friend.id}.txt`;
          const sharedSecretExists = await window.electron.checkFileExists(sharedSecretFileName);

          // 如果共享密钥文件不存在，那么生成并保存共享密钥
          if (!sharedSecretExists) {
            await this.retrieveAndSaveSharedSecret(friend.id);
          }
        }
      } catch (error) {
        console.error(error);
      }
    },
    async retrieveAndSaveSharedSecret(friendId) {
      const jwt = localStorage.getItem('jwt');
      const userId = localStorage.getItem('userid');
      try {
        const response = await axios.get(`/api/userkeys/retrieve?id=${friendId}`, {
          headers: { Authorization: `Bearer ${jwt}` }
        });

        const friendPublicKeyBase64 = response.data;
        const friendPublicKey = Uint8Array.from(atob(friendPublicKeyBase64), c => c.charCodeAt(0));
        const sharedSecretSize = 32;
        const sharedSecretPtr = Module._malloc(sharedSecretSize);
        const friendPublicKeyPtr = Module._malloc(friendPublicKey.length);
        Module.HEAPU8.set(friendPublicKey, friendPublicKeyPtr);
        const ciphertextSize = 768;
        const ciphertextPtr = Module._malloc(ciphertextSize);
        Module._Kyberencrypt(ciphertextPtr, sharedSecretPtr, friendPublicKeyPtr);
        const ciphertext = new Uint8Array(Module.HEAPU8.buffer, ciphertextPtr, ciphertextSize);
        const ciphertextBase64 = window.btoa(String.fromCharCode(...ciphertext));
        const sharedtext = new Uint8Array(Module.HEAPU8.buffer, sharedSecretPtr, sharedSecretSize);
        const sharedtextBase64 = window.btoa(String.fromCharCode(...sharedtext));
        // 用后面用户的公钥签的ssa
        window.electron.saveSharedSecret(sharedtextBase64, `ssA_${userId}_${friendId}.txt`);
        let postData = JSON.stringify({
          UserID: userId,
          FriendID: friendId,
          C: ciphertextBase64
        });

        console.log("Sending the following data to the backend:", postData);

        await axios.post('/api/friendc/add', postData, {
          headers: { Authorization: `Bearer ${jwt}` }
        });

        Module._free(sharedSecretPtr);
        Module._free(friendPublicKeyPtr);
      } catch (error) {
        console.error('Failed to retrieve the public key and save the shared secret:', error);
      }
    }



  },
  mounted() {
    this.getFriends();
  }
};
</script>

<style scoped>
.chat-app {
  display: flex;
  height: 100vh;
}

.chat-window {
  flex: 1;
}
</style>
