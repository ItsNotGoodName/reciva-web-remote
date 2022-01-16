<template>
  <button class="button" :class="{ 'is-loading': loading }">
    <slot></slot>
  </button>
</template>

<script>
export default {
  props: {
    loading: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

<style lang="scss"></style>
