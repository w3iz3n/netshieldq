<template>
  <div>
    <div v-if="loading">Processing...</div>
    <div v-else-if="loggedIn">Redirecting...</div>
    <div v-else>Error or waiting for action</div>
  </div>
</template>

<script>
import axios from "axios";

export default {
  data() {
    return {
      loading: false,
      loggedIn: false,
      username: '',
    };
  },
  created() {
    this.fetchIPAddress()
  },
  methods: {
    async checkAndHandleUser() {
      this.loading = true;
      try {
        const userExists = await this.checkUserExists();
        if (userExists) {
          await this.login(this.username);
        } else {
          await this.register(this.username);
          await this.login(this.username);
        }
        this.loggedIn = true;
        this.$router.push('/chat');
      } catch (error) {
        console.error('Error during user handling:', error);
        alert('Operation failed: ' + error.message); // Provide user feedback
      } finally {
        this.loading = false;
      }
    },
    async checkUserExists() {
      const response = await fetch('/api/users');
      const users = await response.json();
      return users.some(user => user.username === this.username);
    },
    async register(ip) {
      const response = await fetch('/api/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          username: ip,
          password: ip
        })
      });
      if (!response.ok) {
        throw new Error('Registration failed');
      }
    },
    async login(ip) {
      const response = await fetch('/api/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          username: ip,
          password: ip
        })
      });
      if (!response.ok) {
        throw new Error('Login failed');
      }
      const data = await response.json();
      localStorage.setItem('jwt', data.token);
    },
    fetchIPAddress() {
      axios.get('/api/get-ip/')
          .then(response => {
            this.username = response.data;
            this.checkAndHandleUser(); // Initiate the login process once the IP is fetched
          })
          .catch(error => {
            console.error('Error fetching IP:', error);
            this.username = 'Failed to fetch IP';
            alert('Failed to fetch IP: ' + error.message); // Inform the user if IP fetch fails
          });
    }
  }
};
</script>

<style>
/* Add your styles here */
</style>
