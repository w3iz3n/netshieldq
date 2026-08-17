<template>
  <div>
    <h1>Friends</h1>
    <div v-for="friend in friends" :key="friend.Username">
      {{ friend.Username }}
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      friends: []
    };
  },
  mounted() {
    this.getFriends();
  },
  methods: {
    getFriends() {
      fetch('http://localhost:3000/api/friend/get', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + localStorage.getItem('jwt'),
        },
      })
          .then(response => response.json())
          .then(data => {
            this.friends = data;
          });
    }
  }
};
</script>