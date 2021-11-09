<template>
  <Dropdown
    v-model="selectedRadio"
    :options="radios"
    :loading="radioConnecting || loading"
    optionLabel="name"
    placeholder="Select a Radio"
  />
</template>

<script>
import Dropdown from "primevue/dropdown";
import { mapState } from "vuex";

export default {
  components: {
    Dropdown,
  },
  computed: {
    selectedRadio: {
      get() {
        return this.$store.state.selectedRadio;
      },
      set(selectedRadio) {
        this.$store.dispatch("setSelectedRadio", selectedRadio);
      },
    },
    ...mapState(["radios", "radioConnecting", "loading"]),
  },
};
</script>

<style lang="scss" scoped>
</style>