import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
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
      statusText.value = '数据加载失败：' + (e?.message || e)
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

  async function doUpdate() {
    if (loading.value) return
    loading.value = true
    try {
      statusText.value = '正在下载更新…'
      await DoUpdate()
      statusText.value = '数据已更新，正在重新加载…'
      await load()
    } catch (e: any) {
      statusText.value = '更新失败：' + (e?.message || e)
    } finally {
      loading.value = false
    }
  }

  function select(id: string) {
    selectedId.value = id
  }

  return {
    all, searchText, selectedProfession, selectedRarity,
    selectedId, selected, filtered,
    statusText, dataVersion, loading,
    professionOptions, rarityOptions,
    load, checkUpdate, doUpdate, select
  }
})
