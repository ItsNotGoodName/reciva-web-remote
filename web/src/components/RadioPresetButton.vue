<template>
  <b-button
    class="is-justify-content-flex-start is-clipped"
    :class="{ 'is-info': selected }"
    :loading="loading"
    :title="preset.name"
    @click="setPreset"
  >
    <b-tag class="mr-2" :class="selected || 'is-info'">{{ preset.number }}</b-tag>
    <span class="is-flex-grow-1">{{ preset.name }}</span>
  </b-button>
</template>

<script>
import BButton from "./Bulma/BButton.vue";
import BTag from "./Bulma/BTag.vue";

export default {
  components: {
    BButton,
    BTag,
  },
  data() {
    return {
      loading: false,
    };
  },
  props: {
    preset: {
      type: Object,
      required: true,
    },
    radio: {
      type: Object,
      required: true,
    },
  },
  computed: {
    selected() {
      return this.radio.preset == this.preset.number;
    },
  },
  methods: {
    setPreset() {
      if (this.loading) return;
      this.loading = true;
      this.$store.dispatch("setRadioPreset", this.preset.number).finally(() => {
        this.loading = false;
      });
    },
  },
};
</script>

<style lang="scss" scoped></style>
