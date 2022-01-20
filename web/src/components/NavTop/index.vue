<template>
  <nav class="navbar is-fixed-top is-flex p-2" role="navigation" aria-label="top navigation">
    <radio-state class="mr-2" />
    <radio-title class="is-flex-grow-1" />
  </nav>
</template>

<script>
import RadioState from "./RadioState.vue";
import RadioTitle from "./RadioTitle.vue";

export default {
  components: {
    RadioState,
    RadioTitle,
  },
};
</script>

<style lang="scss" scoped>
nav {
  border-bottom: 1px solid #ddd;
}
</style>
