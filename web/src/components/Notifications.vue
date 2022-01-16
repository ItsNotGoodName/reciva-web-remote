<template>
  <div class="container">
    <b-notification
      @close="deletePopup(id)"
      :class="notificationClass(popup.type)"
      class="m-4"
      v-for="(popup, id) in popups"
      :key="id"
    >
      {{ popup.text }}
    </b-notification>
  </div>
</template>

<script>
import { mapActions, mapState } from "vuex";
import {
  MESSAGE_SUCCESS,
  MESSAGE_ERROR,
  MESSAGE_WARNING,
  MESSAGE_INFO,
} from "../constants";

import BNotification from "./Bulma/BNotification.vue";

export default {
  components: { BNotification },
  computed: {
    ...mapState(["popups"]),
  },
  methods: {
    ...mapActions(["deletePopup"]),
    notificationClass(type) {
      if (type == MESSAGE_SUCCESS) {
        return "is-success";
      }
      if (type == MESSAGE_ERROR) {
        return "is-danger";
      }
      if (type == MESSAGE_WARNING) {
        return "is-warning";
      }
      if (type == MESSAGE_INFO) {
        return "is-info";
      }
      return "is-primary";
    },
  },
};
</script>

<style lang="scss" scoped></style>
