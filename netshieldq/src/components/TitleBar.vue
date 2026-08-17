<template>
  <div id="title-bar">
    <div id="title-bar-content">
      <button @click="minimizeWindow">-</button>
      <button @click="toggleMaximizeWindow">⬜</button>
      <button @click="closeWindow">✕</button>
    </div>
  </div>
</template>

<script>
export default {
  name: "TitleBar",
  methods: {
    minimizeWindow() {
      window.electron.send('minimize-window');
    },
    toggleMaximizeWindow() {
      window.electron.send('maximize-window');
    },
    closeWindow() {
      window.electron.send('close-window');
    }
  }
}
</script>

<style scoped>
#title-bar {
  background-color: #dae5ff;
  height: 30px;
  -webkit-app-region: drag;
  display: flex;
  justify-content: flex-end;
  align-items: center;
}

#title-bar-content button {
  background-color: transparent;
  border: none;
  color: #4a4a4a;
  font-size: 14px;
  width: 45px;
  height: 100%;
  cursor: pointer;
  -webkit-app-region: no-drag;
}

#title-bar-content button:hover {
  background-color: rgba(74, 74, 74, 0.1);
}
</style>
