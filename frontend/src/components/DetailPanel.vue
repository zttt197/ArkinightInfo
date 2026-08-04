<script setup lang="ts">
import { computed } from 'vue'
import { useOperatorStore } from '../stores/operator'
import RangeGrid from './RangeGrid.vue'

const store = useOperatorStore()

const op = computed(() => store.selected)
</script>

<template>
  <div class="h-full overflow-hidden">
    <!-- Empty state -->
    <div
      v-if="!op"
      class="h-full flex items-center justify-center text-sm text-[#5B5D6E]"
    >
      从左侧选择一名干员查看详情
    </div>

    <!-- Detail -->
    <div v-else class="h-full overflow-y-auto pl-1 pr-2">
      <!-- Basic info -->
      <div class="flex gap-4">
        <div class="w-24 h-24 rounded-lg bg-[#252640] flex items-center justify-center flex-shrink-0 overflow-hidden">
          <img v-if="op.avatarPath" :src="op.avatarPath" class="w-full h-full object-cover" />
          <span v-else class="text-[38px] text-[#5B5D6E]">{{ op.initial }}</span>
        </div>
        <div class="min-w-0">
          <div class="flex items-center gap-2.5 flex-wrap">
            <span class="text-[26px] font-bold text-[#CDD6F4]">{{ op.name }}</span>
            <span class="text-sm text-[#6C6E80] mt-1">{{ op.appellation }}</span>
            <span class="bg-[#2A2218] text-[#F9E2AF] text-xs px-1.5 py-0.5 rounded mt-1">{{ op.rarityText }}</span>
          </div>
          <div class="text-base text-[#F9E2AF] mt-1.5">{{ op.stars }}</div>
          <div class="flex flex-wrap items-center gap-2 mt-2.5">
            <span class="bg-[#252640] text-[#A5ADC6] text-xs px-2 py-0.5 rounded">{{ op.classLabel }}</span>
            <span class="bg-[#252640] text-[#A5ADC6] text-xs px-2 py-0.5 rounded">{{ op.positionLabel }}</span>
            <span class="bg-[#252640] text-[#A5ADC6] text-xs px-2 py-0.5 rounded">{{ op.nationLabel }}</span>
            <span class="text-xs text-[#8A8DA0]">{{ op.tagsText }}</span>
          </div>
          <div class="text-xs text-[#8A8DA0] mt-1.5">部署费用 {{ op.deployCost }}</div>
          <div class="text-xs text-[#8A8DA0] mt-0.5">再部署 {{ op.redeployText }}</div>
        </div>
      </div>

      <!-- Phase table -->
      <div class="mt-4.5">
        <h3 class="text-[15px] font-bold text-[#CDD6F4] mt-4.5 mb-1.5">面板属性（满级）</h3>
        <div class="grid grid-cols-[150px_80px_80px_70px_60px_70px_90px_60px] text-xs font-bold text-[#8A8DA0] mb-1">
          <span>阶段</span><span>生命</span><span>攻击</span><span>防御</span><span>法抗</span><span>部署</span><span>再部署</span><span>阻挡</span>
        </div>
        <div
          v-for="p in op.phases"
          :key="p.label"
          class="grid grid-cols-[150px_80px_80px_70px_60px_70px_90px_60px] text-[13px] text-[#CDD6F4] py-0.5 border-b border-[#1E1F2E] last:border-b-0"
        >
          <span>{{ p.label }}</span><span>{{ p.hp }}</span><span>{{ p.atk }}</span><span>{{ p.def }}</span>
          <span>{{ p.res }}</span><span>{{ p.cost }}</span><span>{{ p.redeploy }}</span><span>{{ p.block }}</span>
        </div>
      </div>

      <!-- Attack Range -->
      <template v-if="op.phases.some(p => p.rangeGrid && p.rangeGrid.rows > 0)">
        <h3 class="text-[15px] font-bold text-[#CDD6F4] mt-4.5 mb-1">攻击范围</h3>
        <div v-for="p in op.phases" :key="'r'+p.label" class="mt-1.5">
          <span class="text-xs text-[#8A8DA0] mr-2">{{ p.label }}</span>
          <RangeGrid :grid="p.rangeGrid" />
        </div>
      </template>

      <!-- Trust -->
      <h3 class="text-[15px] font-bold text-[#CDD6F4] mt-4.5 mb-1">信赖加成</h3>
      <p class="text-[13px] text-[#A5ADC6]">{{ op.trustText }}</p>

      <!-- Talents -->
      <h3 class="text-[15px] font-bold text-[#CDD6F4] mt-4.5 mb-1">天赋</h3>
      <div v-for="t in op.talents" :key="t.title" class="mt-1.5">
        <div class="flex items-center gap-2">
          <span class="text-sm font-semibold">{{ t.title }}</span>
          <span class="bg-[#252640] text-[#8A8DA0] text-[11px] px-1.5 py-px rounded">{{ t.meta }}</span>
        </div>
        <p class="text-[13px] text-[#A5ADC6] mt-1 leading-5">{{ t.desc }}</p>
      </div>

      <!-- Skills -->
      <h3 class="text-[15px] font-bold text-[#CDD6F4] mt-4.5 mb-1">技能</h3>
      <div v-for="s in op.skills" :key="s.title" class="mt-1.5">
        <div class="flex items-center gap-2 flex-wrap">
          <span class="text-sm font-semibold">{{ s.title }}</span>
          <span class="text-[11px] text-[#6C6E80]">{{ s.meta }}</span>
        </div>
        <p class="text-[13px] text-[#A5ADC6] mt-1 leading-5">{{ s.desc }}</p>
        <RangeGrid v-if="s.rangeGrid && s.rangeGrid.rows > 0" :grid="s.rangeGrid" class="mt-2" />
      </div>

      <!-- Base skills -->
      <h3 class="text-[15px] font-bold text-[#CDD6F4] mt-4.5 mb-1">基建技能</h3>
      <div v-for="b in op.baseSkills" :key="b.title" class="mt-1.5">
        <div class="flex items-center gap-2">
          <span class="text-sm font-semibold">{{ b.title }}</span>
          <span class="text-[11px] text-[#6C6E80]">{{ b.meta }}</span>
        </div>
        <p class="text-[13px] text-[#A5ADC6] mt-1 leading-5">{{ b.desc }}</p>
      </div>

      <div class="h-4" />
    </div>
  </div>
</template>
