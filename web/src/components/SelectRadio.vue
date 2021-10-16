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
    <div class="flex space-x-1">
      <button
        :key="k"
        v-for="k in Object.keys(radios)"
        v-bind:class="{ 'bg-blue-500 text-white': k == uuid, 'hover:bg-blue-300 hover:text-white': k != uuid }"
        @click="selectRadio(k)"
        class="rounded-xl p-3"
      >{{ radios[k] }}</button>
    </div>
    <div class="ml-auto">
      <button
        @click="discoverRadios()"
        class="ml-auto bg-yellow-300 hover:bg-yellow-500 rounded-xl p-3"
      >Discover</button>
    </div>
  </div>
</template>