<template>
  <div class="message">
    <span>{{ message.sender }}:</span>
    <p>{{ message.text }}</p>
  </div>
</template>

<script>
export default {
  props: {
    message: {
      type: Object,
      required: true
    },
  },
  data() {
    return {
      name:'Message',
    }
  }
};
</script>

<style scoped>
.message {
  margin-bottom: 10px;
}

.message span {
  font-weight: bold;
}

.message p {
  margin: 0;
}
</style>
