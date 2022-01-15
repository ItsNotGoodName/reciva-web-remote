<template>
  <button class="button is-success is-small" @click="togglePower">
    <span class="icon">
      <i class="fas fa-edit"></i>
    </span>
  </button>
</template>

<script>
export default {
  props: {
    preset: {
      type: Object,
      required: true,
    },
  },
};
</script>

<style lang=""></style>
