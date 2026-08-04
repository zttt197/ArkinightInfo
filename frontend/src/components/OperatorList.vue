<script setup lang="ts">
import { useOperatorStore } from '../stores/operator'

const store = useOperatorStore()
</script>

<template>
  <div class="bg-[#1A1B26] border border-[#2A2B3D] rounded-lg flex flex-col overflow-hidden">
    <div class="flex-1 overflow-y-auto">
      <div
        v-for="op in store.filtered"
        :key="op.id"
        class="flex items-center gap-3 px-3 py-2 cursor-pointer border-b border-[#252640] last:border-b-0 transition-colors"
        :class="store.selectedId === op.id ? 'bg-[#1E2548]' : 'hover:bg-[#222336]'"
        @click="store.select(op.id)"
      >
        <div class="w-12 h-12 rounded-md bg-[#252640] flex items-center justify-center flex-shrink-0 overflow-hidden">
          <img v-if="op.avatarPath" :src="op.avatarPath" class="w-full h-full object-cover" />
          <span v-else class="text-xl text-[#5B5D6E]">{{ op.initial }}</span>
        </div>
        <div class="min-w-0">
          <div class="flex items-baseline gap-2">
            <span class="text-[15px] font-semibold text-[#CDD6F4] truncate">{{ op.name }}</span>
            <span class="text-[11px] text-[#6C6E80] flex-shrink-0">{{ op.appellation }}</span>
          </div>
          <div class="text-xs text-[#F9E2AF] mt-0.5">{{ op.stars }}</div>
          <div class="text-[11px] text-[#8A8DA0] mt-0.5">{{ op.classLabel }}</div>
        </div>
      </div>
      <div
        v-if="store.filtered.length === 0 && !store.loading"
        class="flex items-center justify-center h-full text-sm text-[#5B5D6E]"
      >
        没有匹配的干员
      </div>
    </div>
  </div>
</template>
