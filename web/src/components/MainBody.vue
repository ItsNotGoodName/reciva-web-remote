<template>
  <div class="container px-2" style="margin-top: 1rem; margin-bottom: 6rem">
    <PresetPlay v-if="page == 'play'" />
    <PresetEdit v-else />
  </div>
</template>

<script>
import { mapState } from "vuex";

import PresetPlay from "./PresetPlay/index.vue";
import PresetEdit from "./PresetEdit/index.vue";

export default {
  computed: {
    ...mapState(["page"]),
  },
  components: {
    PresetPlay,
    PresetEdit,
  },
};
</script>

<style lang="scss" scoped></style>
