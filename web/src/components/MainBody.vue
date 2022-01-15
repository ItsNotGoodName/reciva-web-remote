<template>
  <div class="container">
    <RadioPreset v-if="page == 'play'" />
    <PresetTable v-else />
  </div>
</template>

<script>
import { mapState } from "vuex";

import RadioPreset from "./RadioPreset.vue";
import PresetTable from "./PresetTable.vue";

export default {
  data() {
    return {
      loading: false,
    };
  },
  computed: {
    ...mapState(["page"]),
  },
  components: {
    RadioPreset,
    PresetTable,
  },
};
</script>

<style></style>
