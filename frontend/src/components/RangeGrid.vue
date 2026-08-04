<script setup lang="ts">
import type { RangeGrid } from '../../bindings/arkinightinfo/data/models'

defineProps<{ grid: RangeGrid }>()
</script>

<template>
  <div
    v-if="grid.rows > 0 && grid.cols > 0"
    class="inline-grid gap-0.5 rounded overflow-hidden"
    :style="{ gridTemplateColumns: `repeat(${grid.cols}, 24px)`, gridTemplateRows: `repeat(${grid.rows}, 24px)` }"
  >
    <div
      v-for="(cell, i) in grid.grid"
      :key="i"
      class="w-6 h-6 rounded-sm"
      :class="{
        'bg-[#F9A825]': cell === 'self',
        'bg-[#3D2E1A]': cell === 'range',
        'bg-transparent': cell === 'empty'
      }"
    />
  </div>
</template>
