<script setup lang="ts">
import { onMounted } from 'vue'
import { useOperatorStore } from './stores/operator'
import FilterBar from './components/FilterBar.vue'
import OperatorList from './components/OperatorList.vue'
import DetailPanel from './components/DetailPanel.vue'
import StatusBar from './components/StatusBar.vue'

const store = useOperatorStore()

onMounted(async () => {
  await store.load()
  store.checkUpdate()
})
</script>

<template>
  <div class="h-screen flex flex-col bg-[#12131C] select-none">
    <FilterBar />
    <div class="flex-1 flex overflow-hidden px-3.5">
      <OperatorList class="w-[310px] flex-shrink-0 mr-3.5" />
      <DetailPanel class="flex-1 min-w-0" />
    </div>
    <StatusBar />
  </div>
</template>
