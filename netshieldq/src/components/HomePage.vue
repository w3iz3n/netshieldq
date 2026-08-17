<template>
  <div class="container">
    <div class="logo-title-container">
      <img src="../assets/logo.png" alt="NetShielQ Logo" class="logo">
      <h1 class="title">NetShielQ</h1>
      <p class="subtitle">基于后量子Kyber的内外网通信系统</p>
    </div>
    <div class="button-container">
      <button @click="redirectToIntraScan">内网登录</button><br>
      <button @click="redirectToPublicLogin">公网登录</button><br>
    </div>
  </div>
</template>
<script>
export default {
  name: 'HomePage',
  methods: {
    redirectToIntraScan() {
      this.$router.push('/intrascan');
    },
    redirectToPublicLogin() {
      this.$router.push('/plogin');
      console.log("公网登录");
    }
  }
}
</script>

<style scoped>
.container {
  background-image: url('../assets/bg.jpg');
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
}
.logo-title-container {
  display: flex;
  flex-direction: column;
  align-items: center;
}
.logo {
  width: 120px; /* 根据实际需要调整 */
  height: auto;
}
.title {
  color: transparent;
  background-image: linear-gradient(to right, #ffffff, #f2f2f2);
  -webkit-background-clip: text;
  background-clip: text;
  margin: 5px 0; /* 小间距让标题与logo靠近 */
}
.subtitle {
  color: #f2f2f2; /* 浅灰色 */
  font-size: 1.2em; /* 根据实际需要调整 */
}
.description {
  color: transparent;
  background-image: linear-gradient(to right, #ffffff, #f2f2f2);
  -webkit-background-clip: text;
  background-clip: text;
  margin-top: 20px; /* 调整描述与按钮间的间距 */
}
.button-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-top: 20px;
  width: 250px;
}
button {
  width: 200px;
  height: 60px;
  background: linear-gradient(to bottom, #ffffff, #f2f2f2);
  border: none;
  border-radius: 10px;
  color: #333;
  cursor: pointer;
  margin-bottom: 10px;
}
button:hover {
  background: linear-gradient(to bottom, #0080ff, #0056b3);
}
body{
  margin: 0;
  padding: 0;
}
</style>

