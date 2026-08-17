<template>
  <div class="container">
    <!-- 发送好友申请 -->
    <div class="submit-request-container">
      <h1>添加好友</h1>
      <form @submit.prevent="submitFriendRequest" class="submit-form">
        <input v-model="username" placeholder="输入用户名" class="input-field">
        <button type="submit" class="submit-btn">发送请求</button>
      </form>
    </div>
    <div class="divider"></div>
    <!-- 处理收到的好友申请 -->
    <div class="friend-requests-container">
      <h1>好友申请</h1>
      <div class="friend-requests-list">
        <div v-for="request in paginatedFriendRequests" :key="request.username" class="request-item">
          <span class="request-text">{{ request.username }} 请求添加你为好友</span>
          <el-select v-model="request.status" @change="updateFriendRequestStatus(request)" class="status-select" :disabled="request.status !== 'pending'" size="small">
            <el-option label="等待中" value="pending"></el-option>
            <el-option label="接受" value="accepted"></el-option>
            <el-option label="拒绝" value="rejected"></el-option>
          </el-select>
        </div>
      </div>
      <div class="pagination">
        <el-button @click="changePage(currentPage - 1)" :disabled="currentPage <= 1" class="page-btn" size="small">上一页</el-button>
        <el-button @click="changePage(currentPage + 1)" :disabled="currentPage >= totalPages" class="page-btn" size="small">下一页</el-button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      friendRequests: [],
      username: '',
      currentPage: 1,
      pageSize: 3,
    };
  },
  computed: {
    totalPages() {
      return Math.ceil(this.friendRequests.length / this.pageSize);
    },
    paginatedFriendRequests() {
      const start = (this.currentPage - 1) * this.pageSize;
      const end = start + this.pageSize;
      return this.friendRequests.slice(start, end);
    }
  },
  async created() {
    await this.loadFriendRequests();
  },
  methods: {
    async submitFriendRequest() {
      const response = await fetch('/api/friend/add', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('jwt')}`
        },
        body: JSON.stringify({ username: this.username })
      });
      if (response.ok) {
        this.$message.success('好友申请已发送');
        this.username = '';
        this.loadFriendRequests();
      } else {
        this.$message.error('好友申请发送失败');
      }
    },
    async updateFriendRequestStatus(request) {
      const response = await fetch(`/api/friend/accept`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('jwt')}`
        },
        body: JSON.stringify({ username: request.username, status: request.status })
      });
      if (response.ok) {
        this.$message.success('状态已更新');
        this.loadFriendRequests();
      } else {
        const errorMessage = await response.text();
        this.$message.error(`状态更新失败: ${errorMessage}`);
        this.loadFriendRequests(); // 可选：在错误时重置状态
      }
    },
    async loadFriendRequests() {
      const response = await fetch('/api/friend/requests', {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('jwt')}`
        }
      });
      if (response.ok) {
        this.friendRequests = await response.json();
      } else {
        this.$message.error('加载好友请求失败');
      }
    },
    changePage(page) {
      this.currentPage = page;
    }
  }
};
</script>

<style scoped>
.container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  width: 90%;
  max-width: 1200px;
  margin: 8px auto;
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  background-color: #fff;
}

.submit-request-container, .friend-requests-container {
  padding: 20px;
  background-color: #f4f7ff;
  border-radius: 10px;
}

h1 {
  font-size: 15px;
  color: #4a4a4a;
  margin-bottom: 15px;
}

.submit-form {
  display: flex;
  gap: 10px;
  align-items: center;
}

.input-field {
  flex-grow: 1;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 5px;
  background-color: #f9f9f9;
}

.status-select {
  width: 150px;
}

.submit-btn {
  padding: 10px 20px;
  background-color: #789afd;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.submit-btn:hover {
  background-color: #5680e9;
}

.friend-requests-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.request-item {
  padding: 10px;
  border: 1px solid #eee;
  border-radius: 5px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #fff;
}

.request-text {
  flex-grow: 1;
  margin-right: 10px;
  color: #4a4a4a;
}

.pagination {
  display: flex;
  justify-content: space-between;
  padding: 10px 0;
}

.divider {
  border-left: 1px solid #eee;
  margin: 20px 0;
  height: auto;
}
</style>
