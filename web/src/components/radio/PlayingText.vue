<template>
  <div v-if="radioReady" class="flex">
    <Badge size="large" class="my-auto mr-2">
      {{ radio.title }}
    </Badge>
    <div class="my-auto flex-1 ellipsis">
      {{ radio.metadata }}
    </div>
  </div>
</template>

<script>
import Badge from "primevue/badge";

import { mapGetters, mapState } from "vuex";
export default {
  components: {
    Badge,
  },
  computed: {
    ...mapState(["radio"]),
    ...mapGetters(["radioReady"]),
  },
};
</script>

<style scoped>
.ellipsis {
  text-overflow: ellipsis;

  /* Required for text-overflow to do anything */
  white-space: nowrap;
  overflow: hidden;
}
</style>