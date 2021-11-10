<template>
  <Button
    @click="togglePage"
    :class="{ 'p-button-outlined': page != 'edit' }"
    label="Edit"
  />
</template>

<script>
import { mapState, mapActions } from "vuex";

import Button from "primevue/button";

export default {
  components: { Button },
  methods: {
    ...mapActions(["showEditPage", "showPlayerPage"]),
    togglePage() {
      if (this.page == "edit") {
        this.showPlayerPage();
      } else {
        this.showEditPage();
      }
    },
  },
  computed: {
    ...mapState(["page"]),
  },
};
</script>

<style lang="scss" scoped>
</style>