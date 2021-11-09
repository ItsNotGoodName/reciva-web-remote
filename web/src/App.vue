<script>
import { mapActions, mapState } from "vuex";
import Toast from "primevue/toast";

export default {
  components: {
    Toast,
  },
  mounted() {
    this.$toast.add({
      severity: "success",
      summary: "Success Message",
      detail: "Order submitted",
      life: 3000,
    });
  },
  methods: {
    ...mapActions(["clearMessage"]),
  },
  computed: {
    ...mapState({}),
  },
};
</script>

<template>
  <div>
    <Toast position="bottom-right" />
  </div>
</template>

<style>
</style>
