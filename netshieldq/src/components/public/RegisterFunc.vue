<template>
  <div class="register">
    <h2>注册</h2>
    <form @submit.prevent="register">
      <div>
        <label for="username">用户名</label>
        <input type="text" id="username" v-model="username" required>
      </div>
      <div>
        <label for="password">密码</label>
        <input type="password" id="password" v-model="password" required>
      </div>
      <div>
        <label for="email">邮箱</label>
        <input type="email" id="email" v-model="email" required>
      </div>
      <button type="submit">注册</button>
      <p v-if="error" style="color: red;">{{ error }}</p>
    </form>
  </div>
</template>

<script>
export default {
  data() {
    return {
      username: '',
      password: '',
      email: '',
      error: null,
      name: 'RegisterFunc',
    };
  },
  methods: {
    async register() {
      try {
        const response = await fetch('/api/register', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            username: this.username,
            password: this.password,
            email: this.email
          })
        });

        if (!response.ok) {
          const errorData = await response.json();
          throw new Error(errorData.message || 'Registration failed');
        }

        this.$router.push('/plogin');
      } catch (error) {
        this.error = error.message;
      }
    },
  }
};
</script>
<style scoped>
body {
  background-image: url('../../assets/bg.jpg'); /* Consistent background image */
  background-size: cover;
  background-position: center;
}

.register {
  width: 100%;
  max-width: 320px; /* Standard width for both forms */
  padding: 20px;
  margin: 50px auto; /* Central alignment */
  box-shadow: 0 4px 8px rgba(0,0,0,0.1); /* Consistent shadow */
  border-radius: 8px;
  background: rgba(245, 245, 245, 0.8); /* Semi-transparent Morandi-inspired background */
  color: #676767; /* Consistent text color */
}

.register form {
  display: flex;
  flex-direction: column; /* Consistent form structure */
}

.register form div {
  margin-bottom: 20px; /* Consistent spacing */
}

.register form label {
  display: block;
  margin-bottom: 8px;
  color: #5c5c5c; /* Morandi gray for label text */
  font-size: 14px;
  font-weight: bold; /* Bold labels */
}

.register form input[type="text"],
.register form input[type="password"],
.register form input[type="email"] {
  padding: 10px;
  border: 1px solid #ddd; /* Consistent borders */
  border-radius: 4px;
  background: rgba(255, 255, 255, 0.9); /* Lighter background for input fields */
  font-size: 16px;
  color: #5c5c5c; /* Matching label color for text input */
  transition: border 0.3s, box-shadow 0.3s; /* Smooth transitions */
}

.register form input:focus {
  border-color: #A3C4F3; /* Soft blue focus color */
  box-shadow: 0 0 8px rgba(163, 196, 243, 0.6); /* Light blue shadow */
}

.register form button {
  background-color: #88a2c9; /* Soft blue from the Morandi palette */
  color: #fff;
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  margin-top: 10px;
  cursor: pointer;
  transition: background-color 0.3s, transform 0.2s; /* Consistent animations */
}

.register form button:hover {
  background-color: #7291b4; /* A slightly darker blue for hover effect */
  transform: scale(1.05); /* Consistent hover effect */
}

p {
  margin-top: 10px;
  font-size: 14px;
  color: #676767; /* Consistent paragraph color */
}
</style>

