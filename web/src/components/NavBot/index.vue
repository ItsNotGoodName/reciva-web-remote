<template>
  <nav
    class="navbar is-fixed-bottom"
    role="navigation"
    aria-label="bottom navigation"
  >
    <RadioSelect />
    <RadioControl />
  </nav>
</template>

<script>
import RadioSelect from "./RadioSelect.vue";
import RadioControl from "./RadioControl.vue";

export default {
  components: {
    RadioSelect,
    RadioControl,
  },
};
</script>

<style></style>
