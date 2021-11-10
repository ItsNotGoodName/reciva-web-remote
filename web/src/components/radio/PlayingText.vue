<template>
  <div v-if="radioReady" class="flex">
    <Button
      class="p-button-rounded pl-1 mr-2"
      v-if="radio.state == 'Playing'"
      icon="pi pi-play"
    />
    <Button
      class="p-button-warning p-button-rounded mr-2"
      v-else-if="radio.state == 'Connecting'"
      icon="pi pi-spin pi-spinner"
    />
    <Button
      class="p-button-danger p-button-rounded mr-2"
      v-else-if="radio.state == 'Stopped'"
      icon="pi pi-stop"
    />
    <Button class="h-full flex-grow-1">
      {{ radio.title }}
    </Button>
  </div>
</template>

<script>
import Button from "primevue/button";

import { mapGetters, mapState } from "vuex";
export default {
  components: {
    Button,
  },
  computed: {
    ...mapState(["radio"]),
    ...mapGetters(["radioReady"]),
  },
};
</script>

<style scoped>
</style>