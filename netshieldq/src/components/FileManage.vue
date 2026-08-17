<template>
  <el-container style="height: 100vh;">
    <NavigationBar/>
    <el-container>
      <el-header style="height: auto; padding: 20px; background-color: #f4f7ff;">
        <el-input
            placeholder="搜索文件..."
            v-model="searchQuery"
            clearable
            style="width: 100%; border-radius: 20px; box-shadow: 0 4px 6px rgba(0,0,0,0.1);"
        ></el-input>
      </el-header>
      <el-main style="padding: 20px;background: #f4f7ff">
        <el-row :gutter="20">
          <el-col v-for="file in filteredFiles" :key="file.FileName" :span="8">
            <el-card :body-style="{ padding: '15px' }" class="file-card">
              <div class="file-info">
                <img src="../assets/fileicon.webp" alt="File Icon" class="file-icon">
                <div class="file-details">
                  <span class="file-name">{{ file.FileName }}</span>
                  <el-tag size="mini" class="file-timestamp">{{ file.Timestamp }}</el-tag>
                </div>
              </div>
              <div class="file-actions">
                <el-button size="mini" @click="download(file.FileName)">下载</el-button>
                <el-button size="mini" type="danger" @click="deleteFile(file.FileName)">删除</el-button>
              </div>
            </el-card>
          </el-col>
        </el-row>
      </el-main>
    </el-container>
  </el-container>
</template>

<script>
import axios from 'axios';
import NavigationBar from "@/views/NavigationBar.vue";

export default {
  components: {NavigationBar},
  data() {
    return {
      files: [],
      searchQuery: ''
    };
  },
  mounted() {
    this.fetchFiles();
  },
  computed: {
    filteredFiles() {
      return this.files.filter(file =>
          file.FileName.toLowerCase().includes(this.searchQuery.toLowerCase())
      );
    }
  },
  methods: {
    fetchFiles() {
      axios.get('/api/file/userfiles', {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('jwt')}`
        }
      })
          .then(response => {
            this.files = response.data;
          })
          .catch(error => {
            console.error('There was an error fetching the files:', error);
          });
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
            this.$message.error('Failed to download file');
          });
    },
    deleteFile(fileName) {
      const fileToDelete = this.files.find(file => file.FileName === fileName);
      if (fileToDelete) {
        axios.delete(`/api/file/delete?fileName=${fileName}`)
            .then(() => {
              this.files = this.files.filter(file => file.FileName !== fileName);
              this.$message.success('File deleted successfully');
            })
            .catch(error => {
              console.error('Error deleting file:', error);
              this.$message.error('Failed to delete file');
            });
      }
    }
  }
};
</script>

<style scoped>
body {
  background-color: #dae5ff;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.el-input__inner {
  border-radius: 20px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  transition: box-shadow 0.3s ease;
}

.el-input__inner:focus {
  box-shadow: 0 6px 10px rgba(0, 0, 0, 0.15);
}

.el-card.file-card {
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  transition: box-shadow 0.3s ease;
  border-radius: 12px;
  background-color: #dae5ff;
}

.el-card.file-card:hover {
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.2);
}

.file-info {
  display: flex;
  align-items: center;
  gap: 15px;
}

.file-icon {
  width: 50px;
  height: 50px;
  border-radius: 12px;
}

.file-details {
  display: flex;
  width: 120px;
  flex-direction: column;
  justify-content: center;
  color: #4a4a4a;
}

.file-name {
  font-size: 16px;
  font-weight: bold;
}

.file-timestamp {
  font-size: 12px;
  color: #789afd;
}

.file-actions {
  margin-top: 15px;
  display: flex;
  justify-content: space-between;
}

.el-button {
  margin-right: 5px;
  background-color: #789afd;
  color: white;
  border-radius: 10px;
}

.el-button:hover {
  background-color: #5876d1;
}

.el-tag {
  background-color: #ecf2ff;
  color: #4a4a4a;
}

/* 列布局和间距 */
.el-row {
  margin: 0 -10px;
}

.el-col {
  padding: 10px;
}
</style>
