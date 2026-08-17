<template>
  <div class="page-container">
    <div class="title-container">
      <h2><i class="fas fa-shield-alt"></i> Netshieldq</h2>
      <h3>基于后量子Kyber的即时通信</h3>
    </div>
    <div class="login">
      <h4>登录</h4>
      <form @submit.prevent="login">
        <div>
          <label for="username">用户名</label>
          <input type="text" id="username" v-model="username" required>
        </div>
        <div>
          <label for="password">密码</label>
          <input type="password" id="password" v-model="password" required>
        </div>
        <button type="submit">登录</button>
        <p v-if="error" style="color: red;">{{ error }}</p>
      </form>
      <router-link to="/register" class="register-link">注册</router-link>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      username: '',
      password: '',
      error: null,
      name: 'LogIn',
    };
  },
  methods: {
    async login() {
      try {
        const response = await fetch('/api/login', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            username: this.username,
            password: this.password
          })
        });

        if (!response.ok) {
          const errorData = await response.json();
          throw new Error("登录失败" || 'Login failed');
        }
        const data = await response.json();
        localStorage.setItem('jwt', data.token);
        localStorage.setItem('userid', data.userid);
        localStorage.setItem('username', data.username);
        this.$router.push('/chat');
        await this.$store.dispatch('createSocket');
      } catch (error) {
        this.error = "登录失败";
      }
    },
  }
};
</script>

<style scoped>
html, body {
  height: 100%;
  margin: 0;
}

body {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh; /* Full viewport height */
  margin: 0;
}

.page-container {
  height: 100vh;
  width: 100vw;
  //max-width: 360px; /* Adjust if necessary */
  margin: auto;
  background-color: #ecf2ff; /* Soft blue background color */
  border-radius: 10px; /* Rounded corners */
  //min-height: 80vh; /* Ensure it takes up enough space */
}


.title-container {
  top: 10%;
  position: relative;
  text-align: center;
  margin-bottom: 20px;
}
.title-container h2, .title-container h3 {
  font-family: 'Nunito', sans-serif;
  font-size: 24px;
  margin: 0;
}

.title-container h2 {
  color: #4a4a4a;
  font-size: 24px;
  font-weight: bold;
  line-height: 1.5;
  margin: 0;
}

.title-container h3 {
  color: #4a4a4a; /* Adjusted color */
  font-size: 16px; /* Adjusted size */
  font-weight: normal;
  line-height: 1.4; /* Adjust according to your needs */
  margin: 0; /* Removed default margin */
}

.login {
  top: 10%;
  position: relative;
  padding: 20px;
  max-width: 400px; /* Set the maximum width */
  margin: 0 auto;   /* Center the login box */
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.9); /* Slightly translucent background */
  color: #676767;
}

.login h4 {
  text-align: center;
  color: #4a4a4a;
  font-size: 18px; /* Adjusted size */
  margin-bottom: 20px;
}

.login form {
  display: flex;
  flex-direction: column;
}

.login form div {
  margin-bottom: 20px;
}

.login form label {
  display: block;
  margin-bottom: 8px;
  color: #4a4a4a;
  font-size: 14px;
  font-weight: bold;
}

.login form input {
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background: rgba(255, 255, 255, 0.9);
  font-size: 16px;
  color: #4a4a4a;
  transition: border 0.3s, box-shadow 0.3s;
}

.login form input:focus {
  border-color: #A3C4F3;
  box-shadow: 0 0 8px rgba(163, 196, 243, 0.6);
}

.login form button {
  background-color: #88a2c9;
  color: #fff;
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  margin-top: 10px;
  cursor: pointer;
  transition: background-color 0.3s, transform 0.2s;
}

.login form button:hover {
  background-color: #7291b4; /* A slightly darker blue for hover effect */
  transform: scale(1.05);
}

p {
  margin-top: 10px;
  font-size: 14px;
  color: #676767; /* Ensuring readability */
}

.register-link {
  display: block;
  text-align: center;
  margin-top: 10px;
  color: #88a2c9; /* A soft blue from the Morandi palette */
  font-size: 14px;
  text-decoration: none;
}

.register-link:hover {
  text-decoration: underline;
}
</style>
