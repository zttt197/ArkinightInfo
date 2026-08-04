<script setup lang="ts">
import { useOperatorStore } from '../stores/operator'

const store = useOperatorStore()
</script>

<template>
  <div class="px-3.5 py-2.5 text-xs text-[#5B5D6E]">
    {{ store.statusText }}
  </div>
</template>
