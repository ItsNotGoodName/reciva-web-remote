<template>
  <b-button
    class="is-small is-success"
    :loading="loading"
    @click="showPreset(preset.url)"
    :disabled="loading"
  >
    <b-icon icon="fa-edit"></b-icon>
  </b-button>
</template>

<script>
import BButton from "./Bulma/BButton.vue";
import BIcon from "./Bulma/BIcon.vue";

export default {
  components: {
    BButton,
    BIcon,
  },
  data() {
    return {
      loading: false,
    };
  },
  props: {
    preset: {
      type: Object,
    },
  },
  methods: {
    showPreset(url) {
      if (this.loading) return;
      this.loading = true;
      this.$store.dispatch("showPreset", url).finally(() => {
        this.loading = false;
      });
    },
  },
};
</script>

<style lang="scss" scoped></style>
