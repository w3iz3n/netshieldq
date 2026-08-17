<template>
  <div class="friendship-window">
    <FriendRequest class="friend-request"/>
  </div>
</template>
<script>
import FriendRequest from '@/components/friendRequest.vue';
export default {
  components: {
    FriendRequest,
  },
};
</script>
<style scoped>
.friendship-window {
  display: flex;
  flex-direction: column;
  position: relative;
}
.friend-request {
  position: relative;
  top: 0;
  right: 0;
}
</style>