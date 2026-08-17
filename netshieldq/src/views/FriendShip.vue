<template>
  <div class="friendship-app">
    <NavigationBar/>
    <FriendshipWindow class="flex-item"/>
  </div>
</template>
<script setup>
import NavigationBar from "@/views/NavigationBar.vue";
import FriendshipWindow from "@/views/FriendshipWindow.vue";
</script>
<style>
.friendship-app {
  display: flex;
  height: 100vh;
}
.flex-item {
  flex: 1;
}
</style>
