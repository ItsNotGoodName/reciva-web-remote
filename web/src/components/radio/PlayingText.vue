<template>
  <div v-if="radioReady" class="flex flex-wrap">
    <Button class="my-auto flex-none sm:w-auto w-full">
      {{ radio.title }}
    </Button>
    <span class="m-auto p-2 flex-grow-1">
      {{ radio.metadata }}
    </span>
  </div>
</template>

<script>
import Button from "primevue/button";

import { mapGetters, mapState } from "vuex";
export default {
  components: {
    Button,
  },
  computed: {
    ...mapState(["radio"]),
    ...mapGetters(["radioReady"]),
  },
};
</script>

<style scoped>
</style>