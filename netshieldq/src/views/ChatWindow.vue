<template>
  <div class="chat-window">
    <div class="header-bar">
      <span class="user-name">{{ currentFriend.username }}</span>
      <div class="icon-container">
      <span class="user-file-icon">
        <i class="fas fa-folder"></i>
      </span>
        <span class="delete-friend-icon" @click="handleDeleteFriend()">
            <i class="fas fa-user-minus"></i>
          </span>
          </div>


    </div>
    <div class="message-area">
      <div class="message"
           v-for="message in decryptedMessages"
           :key="message.id"
           :class="{ 'message-sent': message.senderId !== currentFriend.id, 'message-received': message.senderId === currentFriend.id }">
        <div v-if="message.messageType !== 'file'" class="bubble">{{ message.text }}</div>
        <div v-else-if="message.messageType === 'file'" class="file-message">
          <i class="fas fa-file"></i>
          <span>{{ message.text || '未知文件' }}</span>
          <a @click.prevent="download(message.text)">下载</a>
        </div>
      </div>
    </div>
    <div class="input-area">
      <input type="text" v-model="newMessage" @keyup.enter="sendMessage" placeholder="输入消息">
      <button @click="sendMessage">加密发送</button>
      <input type="file" id="file" ref="fileInput" style="display: none;" @change="sendFile">
      <button @click="triggerFileInput">发送文件</button>
    </div>
  </div>
</template>

<script>
import { mapState } from "vuex";
import axios from "axios";
import Module from "../../public/js/test_kyber512";

export default {
  props: {
    messages: {
      type: Array,
      required: true
    },
    currentFriend: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      newMessage: '',
      decryptedMessages: []
    };
  },
  methods: {
    async sendMessage() {
      if (this.newMessage.trim() !== '') {
        let date = new Date();
        let timestamp = date.toLocaleString('en-GB', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', second: '2-digit' }).replace(/\//g, '-');
        let filename = `ssA_${this.currentUserID}_${this.currentFriend.id}.txt`;
        let key = window.electron.readSharedSecret(filename);
        const encryptedMessage = window.electron.aesEncrypt(this.newMessage, key);

        const message = {
          id: Date.now(),
          senderId: this.currentUserID,
          receiverId: this.currentFriend.id,
          text: encryptedMessage,
          timestamp: timestamp,
          isHistorical: false,
          messageType: 'text'
        };

        console.log("发送的实际数据：", message.text);
        this.$emit('send-message', message);
        this.$store.commit('ADD_MESSAGE', message);
        this.newMessage = '';
      }
    },
    triggerFileInput() {
      this.$refs.fileInput.click();
    },
    async handleDeleteFriend() {
      try {
        const userId = parseInt(this.currentUserID, 10);  // Ensure this is an integer
        const friendId = parseInt(this.currentFriend.id, 10);  // Ensure this is an integer

        const response = await axios.post('/api/remove-friend', {
          user_id1: userId,
          user_id2: friendId,
        });

        if (response.status === 200) {
          this.$message.success('好友已成功删除');
          // 从好友列表中移除该好友
          this.friends = this.friends.filter(friend => friend.id !== this.currentFriend.id);
        }
      } catch (error) {
        console.error('删除好友时出错:', error);
        // this.$message.error('无法删除好友');
      }
    },


    download(fileName) {
      axios({
        url: `/api/file/download?fileName=${fileName}`,
        method: 'GET',
        responseType: 'blob',
      })
          .then(response => {
            const url = window.URL.createObjectURL(new Blob([response.data]));
            const link = document.createElement('a');
            link.href = url;
            link.setAttribute('download', fileName);
            document.body.appendChild(link);
            link.click();
            link.remove();
            window.URL.revokeObjectURL(url);
          })
          .catch(error => {
            console.error('Error downloading file:', error);
            this.$message.error('文件已删除');
          });
    },
    async fetchCValueFromBackend(friendId) {
      const userId = this.currentUserID;
      try {
        const response = await axios.get(`/api/friendc/get?userID=${friendId}&friendID=${userId}`);
        if (response.status === 200) {
          return response.data;
        } else {
          throw new Error('Failed to fetch C value from backend');
        }
      } catch (error) {
        console.error('Error fetching C value from backend:', error);
        return null;
      }
    },
    async sendFile() {
      const file = this.$refs.fileInput.files[0];
      if (!file) return;
      const formData = new FormData();
      formData.append('file', file);
      formData.append('receiverId', this.currentFriend.id);

      try {
        const response = await axios.post('/api/file/upload', formData, {
          headers: {
            'Content-Type': 'multipart/form-data',
            'Authorization': `Bearer ${localStorage.getItem('jwt')}`,
          }
        });

        this.$refs.fileInput.value = null;

        let date = new Date();
        let timestamp = date.toLocaleString('en-GB', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', second: '2-digit' }).replace(/\//g, '-');
        let filename = `ssA_${this.currentUserID}_${this.currentFriend.id}.txt`;
        let key = window.electron.readSharedSecret(filename);
        const encryptedMessage = window.electron.aesEncrypt(file.name, key);
        const message = {
          id: Date.now(),
          senderId: this.currentUserID,
          receiverId: this.currentFriend.id,
          text: encryptedMessage,
          timestamp: timestamp,
          messageType: 'file',
          isHistorical: false,
        };

        this.$emit('send-message', message);
        this.$store.commit('ADD_MESSAGE', message);
      } catch (error) {
        console.error('Error uploading file:', error.response.data);
      }
    },
    async decryptMessages(messages) {
      if (!messages || !Array.isArray(messages) || messages.some(message => message == null)) {
        console.error("Invalid message array");
        return [];
      }

      if (!this.currentFriend || !this.currentUserID) {
        console.error("Current friend or user ID is not set");
        return messages;
      }

      return Promise.all(messages.map(async message => {
        let decryptedText = '';

        if (message.senderId !== this.currentFriend.id) {
          // 如果是我发送的
          const sharedSecretFileName = `ssA_${this.currentUserID}_${this.currentChat}.txt`;
          const sharedSecretBase64 = window.electron.readSharedSecret(sharedSecretFileName);
          decryptedText = window.electron.aesDecrypt(message.text, sharedSecretBase64);
        } else {
          // 如果是对方发送的
          const cValue = await this.fetchCValueFromBackend(this.currentChat);
          const privateKeyBase64 = window.electron.readPrivateKey(this.currentUserID);
          const privateKey = Uint8Array.from(atob(privateKeyBase64), c => c.charCodeAt(0));

          const sharedSecretSize = 32;
          const ciphertextSize = 768;
          const privateKeySize = 1632;
          const ciphertextPtr = Module._malloc(ciphertextSize);
          const sharedSecretPtr = Module._malloc(sharedSecretSize);
          const privateKeyPtr = Module._malloc(privateKeySize);

          const cValueUint8 = Uint8Array.from(atob(cValue), c => c.charCodeAt(0));
          Module.HEAPU8.set(privateKey, privateKeyPtr);
          Module.HEAPU8.set(cValueUint8, ciphertextPtr);
          Module._Kyberdecrypt(sharedSecretPtr, ciphertextPtr, privateKeyPtr);

          const sharedSecretKey = new Uint8Array(Module.HEAPU8.buffer, sharedSecretPtr, sharedSecretSize);
          const sharedSecretKeyBase64 = window.btoa(String.fromCharCode(...sharedSecretKey));
          decryptedText = window.electron.aesDecrypt(message.text, sharedSecretKeyBase64);
        }

        return {
          ...message,
          senderId: parseInt(message.senderId),
          receiverId: parseInt(message.receiverId),
          text: decryptedText,
        };
      }));
    }

  },
  beforeMount() {
    this.$store.commit('CLEAR_MESSAGES');
  },
  computed: {
    ...mapState({
      historicalMessages: state => state.historicalMessages,
      realTimeMessages: state => state.messages,
      currentUserID: state => state.currentUserID,
      currentChat: state => state.currentChat
    }),
    currentMessages() {
      const friendId = this.currentChat;
      let historical = [...(this.historicalMessages[friendId] || []), ...(this.historicalMessages[this.currentUserID] || [])];
      let realTime = [...(this.realTimeMessages[friendId] || []), ...(this.realTimeMessages[this.currentUserID] || [])];
      return [...historical, ...realTime];
    }
  },
  watch: {
    currentMessages: {
      immediate: true,
      handler(messages) {
        this.decryptMessages(messages).then(decryptedMessages => {
          this.decryptedMessages = decryptedMessages;
        });
      }
    }
  }
};
</script>


<style scoped>
.chat-window {
  background: #f7f9fc;
  display: flex;
  flex-direction: column;
  height: 100%;
  border-radius: 0px;
  //box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  border-top: 1px;

}

.header-bar {
  display: flex;
  justify-content: space-between;
  background-color: #dae5ff;
  padding: 12px 20px;
  //box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
  //border-top-left-radius: 10px;
  //border-top-right-radius: 10px;
}

.user-name {
  font-weight: bold;
  font-size: 1.2em;
  color: #4a4a4a;
}

.icon-container {
  display: flex;
  justify-content: flex-end; /* 使内容靠右对齐 */
  align-items: center; /* 垂直居中对齐 */
}

.user-file-icon {
  margin-right: 10px;
  color: #4a4a4a;
}

.delete-friend-icon {
  color: #4a4a4a; /* 红色，表示删除 */
  cursor: pointer; /* 鼠标悬停时显示为可点击 */
  margin-left: 10px; /* 与其他图标保持间距 */
}

.delete-friend-icon:hover {
  color: #c0392b; /* 悬停时颜色加深 */
}



.message-area {
  flex-grow: 1;
  padding: 10px 20px;
  background-color: #ffffff;
  overflow-y: auto;
  border-bottom: 1px solid #dae5ff;
}

.message {
  display: flex;
  margin-bottom: 15px;
}

.message-sent {
  justify-content: flex-end;
}

.message-received {
  justify-content: flex-start;
}

.bubble {
  padding: 12px 18px;
  border-radius: 20px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
  max-width: 65%;
  word-wrap: break-word;
}

.message-sent .bubble {
  background-color: rgba(179, 198, 247, 0.61);
  border-bottom-right-radius: 5px;
}

.message-received .bubble {
  background-color: #e9effb;
  border-bottom-left-radius: 5px;
}

.file-message {
  display: flex;
  align-items: center;
  padding: 12px 18px;
  background-color: #e0e0e0;
  border-radius: 20px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}

.file-message i {
  margin-right: 8px;
  color: #666;
}

.file-message a {
  color: #789afd;
  text-decoration: none;
  margin-left: auto;
  font-weight: bold;
}

.input-area {
  display: flex;
  align-items: center;
  padding: 10px 20px;
  background-color: #f7f9fc;
}

.input-area input {
  flex: 1;
  padding: 12px 15px;
  border-radius: 20px;
  border: 1px solid #dae5ff;
  box-shadow: inset 0 1px 3px rgba(0, 0, 0, 0.1);
}

.input-area button {
  background-color: #789afd;
  color: white;
  border: none;
  padding: 12px 20px;
  border-radius: 20px;
  margin-left: 10px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.input-area button:hover {
  background-color: #4a4a4a;
}
</style>
