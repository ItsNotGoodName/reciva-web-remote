<template>
  <b-button
    :class="{ 'button is-info': selected }"
    :loading="loading"
    :title="preset.name"
    @click="setPreset"
  >
    <span class="mr-2" :class="selected ? 'tag is-white' : 'tag is-info'">
      {{ preset.number }}
    </span>
    <span class="is-flex-grow-1">
      {{ preset.name }}
    </span>
  </b-button>
</template>

<script>
import BButton from "../Bulma/BButton.vue";

export default {
  components: {
    BButton,
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

<style></style>
