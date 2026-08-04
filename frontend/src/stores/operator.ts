import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { Events } from '@wailsio/runtime'
import {
  LoadOperators,
  GetDataVersion,
  CheckUpdate,
  DoUpdate,
} from '../../bindings/arkinightinfo/operatorservice'
import type { Operator } from '../../bindings/arkinightinfo/data/models'

const professionOptions = [
  '全部职业', '先锋', '近卫', '重装', '狙击', '术师', '医疗', '辅助', '特种'
]

const rarityOptions = [
  '全部星级', '1星', '2星', '3星', '4星', '5星', '6星'
]

export const useOperatorStore = defineStore('operator', () => {
  const all = ref<Operator[]>([])
  const searchText = ref('')
  const selectedProfession = ref('全部职业')
  const selectedRarity = ref('全部星级')
  const selectedId = ref<string | null>(null)
  const statusText = ref('正在加载…')
  const dataVersion = ref('未知')
  const loading = ref(false)
  const progress = ref('')
  const progressPercent = ref(0)

  // Listen to download progress events from Go
  Events.On('dl-progress', (msg: string) => {
    progress.value = msg
    if (msg.includes('%')) {
      const m = msg.match(/(\d+)%/)
      if (m) progressPercent.value = parseInt(m[1])
    }
    // Detect completion
    if (msg.includes('下载完成')) {
      progress.value = ''
      progressPercent.value = 0
      load()
    }
  })

  const filtered = computed(() => {
    const keyword = searchText.value.trim().toLowerCase()
    const prof = selectedProfession.value
    const rarity = selectedRarity.value

    return all.value.filter(o => {
      if (keyword && !o.name.toLowerCase().includes(keyword) && !o.appellation.toLowerCase().includes(keyword)) {
        return false
      }
      if (prof !== '全部职业' && o.classLabel !== prof) return false
      if (rarity !== '全部星级' && o.rarityText !== rarity) return false
      return true
    })
  })

  const selected = computed(() => {
    if (!selectedId.value) return null
    return all.value.find(o => o.id === selectedId.value) ?? null
  })

  const needsDownload = computed(() => all.value.length === 0 && !loading.value)

  async function load() {
    loading.value = true
    statusText.value = '正在加载本地数据…'
    try {
      const ops = await LoadOperators()
      all.value = ops ?? []
      const ver = await GetDataVersion()
      dataVersion.value = ver || '未知'
      statusText.value = `共 ${all.value.length} 名干员 · 数据版本 ${dataVersion.value}`
    } catch (e: any) {
      statusText.value = '数据未下载，请点击「下载资源」'
    } finally {
      loading.value = false
    }
  }

  async function checkUpdate() {
    try {
      const result = await CheckUpdate()
      if (result && result.HasUpdate) {
        statusText.value = result.Message
      }
    } catch {
      // silent
    }
  }

  function startDownload() {
    if (loading.value) return
    loading.value = true
    progressPercent.value = 0
    statusText.value = '开始下载…'
    // Go runs in goroutine, progress comes via events
    DoUpdate()
  }

  function select(id: string) {
    selectedId.value = id
  }

  function resetFilters() {
    searchText.value = ''
    selectedProfession.value = '全部职业'
    selectedRarity.value = '全部星级'
  }

  return {
    all, searchText, selectedProfession, selectedRarity,
    selectedId, selected, filtered,
    statusText, dataVersion, loading,
    progress, progressPercent, needsDownload,
    professionOptions, rarityOptions,
    load, checkUpdate, startDownload, select, resetFilters
  }
})
