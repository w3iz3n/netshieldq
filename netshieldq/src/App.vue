<template>
  <div id="app">
    <TitleBar/>
    <router-view :key="$route.fullPath"></router-view>
  </div>
</template>

<script>
import TitleBar from "@/components/TitleBar.vue";

export default {
  name: 'App',
  components: {TitleBar},
  created() {
    this.$router.push('/home');
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  height: 100vh;
  width: 100vw;
  overflow: hidden; /* 防止出现滚动条 */
  display: flex;
  flex-direction: column;
}

body {
  font-family: 'Arial', sans-serif;
  background-color: #f0f0f0;
  margin: 0;
  padding: 0;
}

TitleBar {
  height: 50px; /* 根据需求调整 */
  width: 100%;
  background-color: #dae5ff; /* 或其他符合主题的颜色 */
  flex-shrink: 0; /* 确保 TitleBar 不会缩小 */
}

router-view {
  flex-grow: 1; /* 占据剩余空间 */
  overflow: auto; /* 如果内容多于屏幕，允许滚动 */
}
</style>
