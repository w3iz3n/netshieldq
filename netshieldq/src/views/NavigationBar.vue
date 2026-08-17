<template>
  <div class="navigation-bar">
    <div class="user-profile">
      <img :src="userAvatar" alt="User Avatar">
    </div>
    <div class="navigation-item" v-for="(item, index) in navigationItems" :key="index" @click="navigate(item.route)">
      <i :class="['icon', item.icon]"></i>
      <span>{{ item.label }}</span>
    </div>
    <div class="logout-button" @click="logout">
      <i class="icon fa fa-sign-out-alt"></i>
      <span>Logout</span>
    </div>
  </div>
</template>

<script>
import axios from "axios";

export default {
  data() {
    return {
      userAvatar: 'https://gra1n-cn.oss-cn-hangzhou.aliyuncs.com/default_avatar.png',
      navigationItems: [
        { route: '/friendship', icon: 'fa fa-envelope', label: '' },
        { route: '/chat', icon: 'fa fa-comments', label: '' },
        { route: '/files', icon: 'fa fa-folder', label: '' },
        { route: '/kyber', icon: 'fa fa-cogs', label: '' },
      ]
    };
  },
  methods: {
    navigate(route) {
      this.$router.push(route);
    },
    logout() {
      localStorage.removeItem('jwt');
      localStorage.removeItem('userid');
      localStorage.removeItem('username');
      this.$router.push('/home');
    },
    async fetchAvatar() {
      try {
        const jwt = localStorage.getItem('jwt');
        const response = await axios.get('/api/get-avatar', {
          headers: {
            'Authorization': `Bearer ${jwt}`,
            'Accept': 'application/json'
          }
        });

        // 获取Base64编码的头像数据
        const base64Avatar = response.data.avatar;

        this.userAvatar = `data:image/png;base64,${base64Avatar}`;

      } catch (error) {
        console.error('Failed to load avatar:', error);
      }
    },

  },
  mounted() {
    this.fetchAvatar(); // 在组件挂载时调用以加载头像
  }
};
</script>

<style scoped>
.navigation-bar {
  width: 80px;
  background-color: #dae5ff;
  color: #4a4a4a;
  box-sizing: border-box;
  padding: 10px;
  height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  border-top: 1px  ;

}

.user-profile img {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  margin-bottom: 20px;
}

.navigation-item {
  width: 100%;
  display: flex;
  align-items: center;
  padding: 10px 0;
  cursor: pointer;
  transition: background-color 0.3s;
}

.navigation-item:hover {
  background-color: #dae5ff;
}

.navigation-item .icon {
  font-size: 24px;
  color: #4a4a4a;
  margin-left: 20px;
}

.navigation-item span {
  margin-left: 10px;
  flex-grow: 1;
  text-align: left;
  color: #4a4a4a;
}

.logout-button {
  margin-top: auto;
  color: #4a4a4a;
  cursor: pointer;
  display: flex;
  align-items: center;
  padding: 10px 0;
}
</style>
