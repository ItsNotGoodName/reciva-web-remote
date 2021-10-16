<script setup>

defineProps({
    connecting: Boolean,
    connected: Boolean
})
</script>

<template>
    <div v-if="!connecting && !connected" class="p-2 bg-red-200">
        <div class="container mx-auto px-2">Lost Connection to Server</div>
    </div>
</template>