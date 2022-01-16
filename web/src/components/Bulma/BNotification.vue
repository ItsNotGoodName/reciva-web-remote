<template>
  <div class="notification">
    <button class="delete" @click="$emit('close')"></button>
    <slot></slot>
  </div>
</template>

<script>
export default {};
</script>

<style lang="scss" scoped></style>
