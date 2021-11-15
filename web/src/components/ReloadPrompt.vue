<template>
  <Message
    v-if="offlineReady || needRefresh"
    @close="close()"
    severity="info"
    :sticky="needRefresh"
  >
    <div class="flex">
      <Button class="mr-2" v-if="needRefresh" @click="updateServiceWorker()">
        Reload
      </Button>
      <div class="my-auto">
        <span v-if="offlineReady">App ready to work offline</span>
        <span v-else>
          New content is available, click on the reload button to update.
        </span>
      </div>
    </div>
  </Message>
</template>

<script setup>
import { useRegisterSW } from "virtual:pwa-register/vue";
import Message from "primevue/message";
import Button from "primevue/button";

const { offlineReady, needRefresh, updateServiceWorker } = useRegisterSW();

const close = async () => {
  offlineReady.value = false;
  needRefresh.value = false;
};
</script>


<style>
</style>
