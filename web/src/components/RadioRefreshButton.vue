<template>
  <b-button
    title="Refresh"
    @click="refresh"
    :loading="loading"
    :class="{ 'is-danger': !radioWSConnected }"
  >
    <b-icon icon="fa-redo" />
  </b-button>
</template>

<script>
import { mapState } from "vuex";
import BButton from "./Bulma/BButton.vue";
import BIcon from "./Bulma/BIcon.vue";

export default {
  data() {
    return {
      loading: false,
    };
  },
  components: {
    BButton,
    BIcon,
  },
  computed: {
    ...mapState({
      radioWSConnected: (state) => state.r.radioWSConnected,
    }),
  },
  methods: {
    refresh() {
      if (this.loading) return;

      this.loading = true;
      this.$store.dispatch("refreshRadioAll").finally(() => {
        this.loading = false;
      });
    },
  },
};
</script>

<style lang="scss" scoped></style>
