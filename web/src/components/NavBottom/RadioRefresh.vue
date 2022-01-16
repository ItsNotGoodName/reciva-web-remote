<template>
  <b-button
    title="Refresh"
    :placeholder="radioRefreshing"
    :class="{ 'is-danger': !radioWSConnected }"
    @click="refreshRadio"
  >
    <b-icon icon="fa-redo" />
  </b-button>
</template>

<script>
import { mapState, mapActions, mapGetters } from "vuex";

import BButton from "../Bulma/BButton.vue";
import BIcon from "../Bulma/BIcon.vue";

export default {
  components: {
    BButton,
    BIcon,
  },
  computed: {
    ...mapState({
      radioRefreshing: (state) => state.r.radioRefreshing,
      radioWSConnected: (state) => state.r.radioWSConnected,
    }),
    ...mapGetters(["radioLoaded"]),
  },
  methods: {
    ...mapActions(["refreshRadio"]),
  },
};
</script>

<style lang="scss" scoped></style>
