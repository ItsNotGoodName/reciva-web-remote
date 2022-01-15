<template>
  <button
    @click="setPreset"
    class="is-fullwidth wrap-button-text"
    :title="preset.name"
    :class="[
      { 'is-loading': this.loading },
      selected ? 'button is-info' : 'button',
    ]"
    :disabled="loading"
  >
    <span class="mr-2" :class="selected ? 'tag is-white' : 'tag is-info'">
      {{ preset.number }}
    </span>
    <span class="is-flex-grow-1" style="word-break: break-all">
      {{ preset.name }}
    </span>
  </button>
</template>

<script>
export default {
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

<style scoped>
.wrap-button-text {
  height: 100%;
  white-space: normal;
}
</style>
