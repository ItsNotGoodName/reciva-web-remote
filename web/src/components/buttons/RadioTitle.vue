<template>
  <div class="dropdown is-hoverable is-flex">
    <div class="dropdown-trigger is-flex-grow-1">
      <button
        :title="radio.title"
        class="button is-info has-text-left is-fullwidth"
        :class="{ 'is-loading': !radioLoaded }"
        aria-haspopup="true"
        aria-controls="dropdown-menu4"
      >
        <span class="tag mr-2 is-white">{{ radio.preset }}</span>
        <span class="is-flex-grow-1">
          {{ radio.title }}
        </span>
      </button>
    </div>
    <div class="dropdown-menu" id="dropdown-menu4" role="menu">
      <div class="dropdown-content">
        <div class="dropdown-item has-text-centered">Stream Information</div>
        <div class="dropdown-item">
          <span class="tag mr-2 is-rounded is-warning">Metadata</span>
          {{ radio.metadata }}
        </div>
        <div class="dropdown-item">
          <span class="tag mr-2 is-rounded is-info">URL</span>
          {{ radio.url }}
        </div>
        <div class="dropdown-item">
          <span class="tag mr-2 is-rounded is-success">New URL</span>
          {{ radio.newUrl }}
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapState, mapGetters } from "vuex";

export default {
  computed: {
    ...mapGetters(["radioLoaded"]),
    ...mapState({
      radio: (state) => state.r.radio,
    }),
  },
};
</script>

<style scoped>
.button.has-text-left {
  justify-content: flex-start;
}
</style>
