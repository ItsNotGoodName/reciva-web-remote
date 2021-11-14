<template>
  <Button
    v-if="page == 'edit'"
    icon="pi pi-play"
    @click="togglePage"
    label="Play"
    title="Play presets."
  />
  <Button
    v-else
    icon="pi pi-pencil"
    @click="togglePage"
    label="Edit"
    title="Edit presets."
  />
</template>

<script>
import { mapState, mapActions } from "vuex";

import Button from "primevue/button";

export default {
  components: { Button },
  methods: {
    ...mapActions(["showEditPage", "showPlayerPage"]),
    togglePage(event) {
      if (this.page == "edit") {
        this.showPlayerPage();
      } else {
        this.showEditPage();
      }
      this.toggle(event);
    },
  },
  props: {
    toggle: {
      type: Function,
      default: () => {},
    },
  },
  computed: {
    ...mapState(["page"]),
  },
};
</script>

<style lang="scss" scoped>
</style>