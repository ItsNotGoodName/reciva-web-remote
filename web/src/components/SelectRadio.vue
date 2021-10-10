<script setup>

defineProps({
  uuid: String,
  radios: Object,
  selectRadio: Function,
  discoverRadios: Function
})

</script>

<template>
  <div class="flex">
    <button
      :key="k"
      v-for="k in Object.keys(radios)"
      v-bind:class="{ 'bg-blue-300': k == uuid }"
      @click="selectRadio(k)"
      class="rounded-xl p-3"
    >{{ radios[k] }}</button>
    <button @click="discoverRadios()" class="ml-auto bg-yellow-300 rounded-xl p-3">Discover</button>
  </div>
</template>