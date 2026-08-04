<script setup lang="ts">
import { useOperatorStore } from '../stores/operator'

const store = useOperatorStore()
</script>

<template>
  <div>
    <div class="flex items-center gap-2.5 px-3.5 pt-3">
      <input
        v-model="store.searchText"
        type="text"
        placeholder="搜索干员名或代号…"
        class="w-[220px] px-2 py-1.5 bg-[#1A1B26] border border-[#2A2B3D] rounded text-sm text-[#CDD6F4] placeholder-[#5B5D6E] outline-none focus:border-[#89B4FA]"
      />
      <select
        v-model="store.selectedProfession"
        class="w-[130px] px-2 py-1.5 bg-[#1A1B26] border border-[#2A2B3D] rounded text-sm text-[#CDD6F4] outline-none focus:border-[#89B4FA]"
      >
        <option v-for="p in store.professionOptions" :key="p" :value="p">{{ p }}</option>
      </select>
      <select
        v-model="store.selectedRarity"
        class="w-[100px] px-2 py-1.5 bg-[#1A1B26] border border-[#2A2B3D] rounded text-sm text-[#CDD6F4] outline-none focus:border-[#89B4FA]"
      >
        <option v-for="r in store.rarityOptions" :key="r" :value="r">{{ r }}</option>
      </select>
      <button
        class="px-3 py-1.5 border border-[#2A2B3D] text-[#8A8DA0] text-sm rounded hover:bg-[#222336] hover:text-[#CDD6F4] transition-colors"
        @click="store.resetFilters()"
      >
        重置
      </button>
      <div class="flex-1" />

      <!-- First time: big download button -->
      <button
        v-if="store.needsDownload"
        :disabled="store.loading"
        class="px-5 py-1.5 bg-[#F9A825] text-[#12131C] text-sm font-medium rounded hover:bg-[#E09820] disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
        @click="store.startDownload()"
      >
        下载资源
      </button>

      <!-- Normal: update button -->
      <button
        v-else
        :disabled="store.loading"
        class="px-4 py-1.5 bg-[#89B4FA] text-[#12131C] text-sm font-medium rounded hover:bg-[#74A8F7] disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
        @click="store.startDownload()"
      >
        检查并更新
      </button>
    </div>

    <!-- Progress bar -->
    <div v-if="store.loading && store.progress" class="px-3.5 pb-2">
      <div class="flex items-center gap-2 text-xs text-[#8A8DA0] mb-1">
        <span>{{ store.progress }}</span>
        <span v-if="store.progressPercent > 0">{{ store.progressPercent }}%</span>
      </div>
      <div class="w-full h-1.5 bg-[#252640] rounded-full overflow-hidden">
        <div
          class="h-full bg-[#89B4FA] rounded-full transition-all duration-300"
          :style="{ width: store.progressPercent > 0 ? store.progressPercent + '%' : '100%' }"
          :class="store.progressPercent === 0 ? 'animate-pulse' : ''"
        />
      </div>
    </div>
  </div>
</template>
