<template>
  <div class="friend-list">
    <!-- 搜索框 -->
    <div class="search-box">
      <input type="text" v-model="searchQuery" placeholder="搜索好友..." @input="filterFriends">
    </div>
    <div class="friend-item"
         v-for="friend in filteredFriends"
         :key="friend.id"
         @click="selectFriend(friend)">
      <!-- 用户图标 -->
      <img :src="friend.avatar || '../assets/default_avatar.png'" alt="User Icon" class="user-icon">
      <div class="friend-info">
        <div class="friend-name">{{ friend.username }}</div>
        <div class="friend-message">
          <span v-if="friendMessages[friend.id] && friendMessages[friend.id].length">
            {{ friendMessages[friend.id].slice(-1)[0].text.slice(0, 15) }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  props: {
    friends: {
      type: Array,
      required: true
    },
    friendMessages: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      searchQuery: ''
    };
  },
  computed: {
    filteredFriends() {
      this.fetchAllAvatars();
      if (!this.searchQuery.trim()) {
        return this.friends;
      }
      return this.friends.filter(friend => friend['username'].toLowerCase().includes(this.searchQuery.toLowerCase()));
    }
  },
  mounted() {
    this.fetchAllAvatars();
  },
  methods: {

    async fetchAvatar(friend) {
      try {
        const jwt = localStorage.getItem('jwt');
        const response = await axios.get(`/api/get-friavatar/${friend.id}`, {
          headers: {
            Authorization: `Bearer ${jwt}`,
            Accept: 'application/json'
          }
        });

        if (response.data && response.data.avatar) {
          console.log(response.data.avatar)
          this.$set(friend, 'avatar', `data:image/png;base64,${response.data.avatar}`);

        } else {
          this.$set(friend, 'avatar', '../assets/default_avatar.png');
        }

      } catch (error) {
        console.error('Failed to load avatar:', error);
        this.$set(friend, 'avatar', '../assets/default_avatar.png');
      }
    },

    async fetchAllAvatars() {
      // 并行加载所有朋友的头像
      const avatarPromises = this.friends.map(friend => this.fetchAvatar(friend));
      await Promise.all(avatarPromises);
    },
    selectFriend(friend) {
      this.$emit('select-friend', friend);
    }
  },
};
</script>

<style scoped>
.search-box {
  padding: 12px;
  background-color: #dae5ff;
  //border-bottom: 2px solid #789afd;
  //box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  //border-radius: 8px;
  transition: all 0.3s ease;
}

.search-box input {
  width: 100%;
  height: 20px;
  border: none;
  //border-radius: 8px;
  //box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.05);
  outline: none;
  transition: box-shadow 0.3s ease;
  background-color: #f9f9ff;
  color: #4a4a4a;
  font-family: 'Nunito', sans-serif;
}

.search-box input:focus {
  box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.2);
}

.friend-list {
  width: 220px;
  border-top: 1px ;
  border-left: 1px ;
  //border-right: 1px solid #789afd;
  overflow-y: auto;
  height: 100%;
  background-color: #f9f9ff;
  font-family: 'Nunito', sans-serif;
  //border-radius: 8px;
}

.friend-item {
  display: flex;
  align-items: center;
  padding: 12px;
  cursor: pointer;
  border-bottom: 1px solid #dae5ff;
  transition: background-color 0.3s;
  border-radius: 8px;
  margin: 6px;
}

.friend-item:hover {
  background-color: #dae5ff;
}

.user-icon {
  width: 36px;
  height: 36px;
  margin-right: 10px;
  border-radius: 50%;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.friend-info {
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.friend-name {
  font-weight: bold;
  color: #4a4a4a;
  margin-bottom: 4px;
}

.friend-message {
  font-size: 0.85em;
  color: #789afd;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
