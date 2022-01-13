<template>
  <button
    v-if="radio.preset == preset.number"
    :title="preset.name"
    class="button is-flex is-info"
    :class="{ 'is-loading': loading }"
    @click="setPreset"
  >
    <span class="tag is-white">{{ preset.number }}</span>
    <span class="is-flex-grow-1">
      {{ preset.name }}
    </span>
  </button>
  <button
    v-else
    :title="preset.name"
    class="button is-flex"
    :class="{ 'is-loading': loading }"
    @click="setPreset"
  >
    <span class="tag is-info">{{ preset.number }}</span>
    <span class="is-flex-grow-1">
      {{ preset.name }}
    </span>
  </button>
</template>

<script>
import { mapState } from "vuex";

export default {
  data() {
    return { loading: false };
  },
  props: {
    preset: {
      type: Object,
      required: true,
    },
  },
  computed: {
    ...mapState({
      radio: (state) => state.r.radio,
    }),
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

<style></style>
