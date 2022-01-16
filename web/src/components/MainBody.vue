<template>
  <div class="container">
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

<style></style>
