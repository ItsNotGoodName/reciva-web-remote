<template>
  <span class="icon">
    <i class="fas" :class="icon" />
  </span>
</template>

<script>
export default {
  props: {
    icon: {
      type: String,
      required: true,
    },
  },
};
</script>

<style lang="scss"></style>
