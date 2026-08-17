<template>
  <div class="container">
    <NavigationBar />
    <div class="content">
      <h1>生成密钥对</h1>
      <div class="key-generation-section">
        <textarea v-model="publicKeyBase64" readonly placeholder="公钥区域"></textarea>
        <textarea v-model="privateKeyBase64" readonly placeholder="私钥区域"></textarea>
        <button @click="generateKeyPair">生成密钥对</button>
      </div>
    </div>
    <div class="right-container">
      <img :src="avatar" alt="User Avatar" class="user-avatar" />
      <button @click="triggerFileInput" class="change-avatar-button">修改头像</button>
      <input
          type="file"
          ref="fileInput"
          @change="updateAvatar"
          accept="image/*"
          class="upload-avatar"
          style="display: none;"
      />
      <h3>{{ username }}</h3>
      <p>{{ email }}</p>
    </div>
  </div>
</template>

<script>
import Module from "../../public/js/test_kyber512.js";
import NavigationBar from "@/views/NavigationBar.vue";
import axios from "axios";
import avatarImg from "@/assets/default_avatar.png";

export default {
  components: { NavigationBar },
  data() {
    return {
      publicKeyBase64: '',
      privateKeyBase64: '',
      username: localStorage.getItem('username') || 'N/A',
      email: localStorage.getItem('email') || 'N/A',
      avatar: avatarImg
    };
  },
  methods: {
    async generateKeyPair() {
      const publicKeySize = 800;
      const privateKeySize = 1632;
      const publicKeyPtr = Module._malloc(publicKeySize);
      const privateKeyPtr = Module._malloc(privateKeySize);

      Module._generateKeyPair(publicKeyPtr, privateKeyPtr);

      const publicKey = new Uint8Array(Module.HEAPU8.buffer, publicKeyPtr, publicKeySize);
      const privateKey = new Uint8Array(Module.HEAPU8.buffer, privateKeyPtr, privateKeySize);

      this.publicKeyBase64 = window.btoa(String.fromCharCode(...publicKey));
      this.privateKeyBase64 = window.btoa(String.fromCharCode(...privateKey));

      const ciphertextSize = 768;
      const sharedSecretSize = 32;
      const ciphertextPtr = Module._malloc(ciphertextSize);
      const sharedSecretPtr = Module._malloc(sharedSecretSize);

      Module._Kyberencrypt(ciphertextPtr, sharedSecretPtr, publicKeyPtr);
      Module._free(publicKeyPtr);
      Module._free(privateKeyPtr);
      Module._free(ciphertextPtr);
      Module._free(sharedSecretPtr);
      window.electron.savePrivateKey(this.privateKeyBase64, localStorage.getItem('userid'));
      await this.sendPublicKey();
    },

    async sendPublicKey() {
      const jwt = localStorage.getItem('jwt');

      try {
        const response = await axios.post('/api/user/sendpk', {
          publicKey: this.publicKeyBase64
        }, {
          headers: { Authorization: `Bearer ${jwt}` }
        });
      } catch (error) {
        console.error('Failed to send the public key:', error);
      }
    },

    triggerFileInput() {
      this.$refs.fileInput.click();
    },

    async updateAvatar(event) {
      const file = event.target.files[0];
      if (file) {
        const formData = new FormData();
        formData.append('avatar', file);

        try {
          const jwt = localStorage.getItem('jwt');
          const response = await axios.post('/api/update-avatar', formData, {
            headers: {
              'Authorization': `Bearer ${jwt}`,
              'Content-Type': 'image/png',
            },
          });
          this.avatar = URL.createObjectURL(file);
          this.$message.success('头像更新成功');
        } catch (error) {
          console.error('Failed to update avatar:', error);
          this.$message.error('头像更新失败');
        }
      }
    },
    async fetchAvatar() {
      try {
        const jwt = localStorage.getItem('jwt');
        const response = await axios.get('/api/get-avatar', {
          headers: {
            'Authorization': `Bearer ${jwt}`,
            'Accept': 'application/json' // 请求的是 JSON 数据
          }
        });

        // 从响应中获取 Base64 编码的头像字符串
        const base64Avatar = response.data.avatar;
        // 设置 img 元素的 src 属性为 Base64 编码的字符串
        this.avatar = `data:image/png;base64,${base64Avatar}`;
      } catch (error) {
        console.error('Failed to load avatar:', error);
      }
    },


    initWasm() {
      return new Promise((resolve) => {
        if (Module.onRuntimeInitialized) {
          resolve();
        } else {
          Module.onRuntimeInitialized = resolve;
        }
      });
    }
  },
  async created() {
    await this.initWasm();
  },
  mounted() {
    this.fetchAvatar();
  },

};
</script>

<style scoped>
.container {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  background: #f7f9fb;
  font-family: 'Nunito', sans-serif;
  //padding: 20px;
  //border-radius: 10px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  max-height: 100vh;
  //overflow: hidden;
}

.content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding-right: 20px;
}

.key-generation-section {
  display: flex;
  flex-direction: column;
  gap: 20px;
  background: #dae5ff; /* Soft blue background */
  border-radius: 10px;
  padding: 20px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

textarea {
  background: #f0f5ff; /* Light background for textarea */
  border: 1px solid #789afd; /* Matching border color */
  border-radius: 8px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.05);
  padding: 10px;
  min-height: 150px;
  resize: none;
  color: #4a4a4a; /* Text color matching the theme */
  width: 100%;
  transition: border-color 0.3s;
}

textarea:hover {
  border-color: #5b7f97; /* Darker border on hover */
}

button {
  background: #789afd; /* Soft blue button */
  color: white;
  border: none;
  border-radius: 5px;
  padding: 10px 15px;
  cursor: pointer;
  width: 100%;
  transition: background-color 0.3s, transform 0.2s;
}

button:hover {
  background: #5674d6; /* Darker hover effect */
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
}

.right-container {
  width: 250px;
  padding: 20px;
  border-left: 1px solid #e0e0e0;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
  background: #ffffff;
  border-radius: 10px;
}

.user-avatar {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  border: 2px solid #ccc;
  margin-bottom: 10px;
}

.change-avatar-button {
  background: #789afd; /* Matching soft blue color */
  color: white;
  border: none;
  border-radius: 5px;
  padding: 10px 15px;
  cursor: pointer;
  width: 100%;
  transition: background-color 0.3s, transform 0.2s;
  text-align: center;
}

.change-avatar-button:hover {
  background: #5674d6; /* Darker hover effect */
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
}

.upload-avatar {
  margin-bottom: 10px;
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 8px;
  background-color: #ffffff;
  cursor: pointer;
  display: none;
}

h1 {
  font-size: 24px;
  color: #4a4a4a;
  margin-bottom: 15px;
}

/* Media Query for Responsive Adjustments */
@media (max-width: 768px) {
  .container {
    flex-direction: column;
    align-items: center;
  }

  .right-container {
    width: 90%;
    border-left: none;
    border-top: 1px solid #e0e0e0;
    margin-top: 20px;
  }

  .content, .right-container {
    width: 90%;
    padding: 15px;
  }
}

</style>
