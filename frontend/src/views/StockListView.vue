<template>
  <div class="stock-list-view p-4 md:p-6">
    <h1 class="text-3xl sm:text-4xl font-bold mb-6 sm:mb-8 text-slate-800 text-center">
      Acciones del Mercado
    </h1>

    <section v-if="stockStore.recommendations.length > 0 || stockStore.isLoadingRecommendations || stockStore.recommendationsError"
              class="mb-8 sm:mb-10 p-5 sm:p-6 bg-gradient-to-br from-sky-300 to-indigo-900 text-white shadow-2xl rounded-xl max-w-screen-2xl mx-auto">
      <h2 class="text-2xl font-semibold mb-4 text-center sm:text-center">Top de Acciones Para Invertir Hoy</h2>

      <div v-if="stockStore.isLoadingRecommendations" class="text-center py-3">
        <p class="animate-pulse">Buscando las mejores oportunidades...</p>
      </div>
      <div v-else-if="stockStore.recommendationsError" class="text-center py-3 text-red-200 bg-red-700 bg-opacity-60 p-3 rounded-md">
        <p class="font-semibold">⚠️ Error al obtener recomendaciones:</p>
        <p class="text-sm">{{ stockStore.recommendationsError }}</p>
      </div>
      <div v-else-if="stockStore.recommendations.length > 0" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-3 gap-4 sm:gap-5">
        <div v-for="rec in stockStore.recommendations" :key="rec.ticker + '-rec'"
              class="bg-white/10 backdrop-blur-md p-5 rounded-lg shadow-lg hover:bg-white/20 transition-all duration-300 ease-in-out transform hover:scale-105">
          <h3 class="text-xl font-bold truncate mb-1">
            <router-link :to="{ name: 'StockDetail', params: { ticker: rec.ticker } }" class="hover:underline" :title="rec.company">
              {{ rec.company }} ({{ rec.ticker }})
            </router-link>
          </h3>
          <p class="text-sm mt-1">Rating: <span class="font-semibold">{{ rec.rating_to }}</span> | Target: <span class="font-semibold">{{ rec.target_to ? '$'+rec.target_to.toFixed(2) : 'N/A' }}</span></p>
          <div class="text-xs mt-2 space-y-1 opacity-90">
            <p class="font-semibold mb-1">Análisis Clave (Score: {{ rec.score.toFixed(0) }}):</p>
            <ul class="list-disc list-inside pl-1">
                <li v-for="(reason, idx) in rec.reasons" :key="idx" class="flex items-start text-xs mb-0.5">
                  <span class="mr-2 text-green-300">{{ getReasonIcon(reason.type) }}</span>
                  <span>{{ reason.details }}</span>
                </li>
            </ul>
          </div>
        </div>
      </div>
      <div v-else class="text-center py-3 text-indigo-100">
        <p>📊 El mercado está tranquilo... No hay nuevas recomendaciones destacadas por ahora.</p>
      </div>
    </section>

    <div class="max-w-screen-2xl mx-auto">
      <div class="mb-6 p-4 bg-white shadow-md rounded-lg flex flex-col md:flex-row gap-4 items-end">
        <div class="flex-grow w-full md:w-auto">
          <label for="search" class="block text-sm font-medium text-slate-700 mb-1">Buscar Acción</label>
          <input
            type="text"
            id="search"
            v-model="searchInputValue"
            @input="debouncedSearch"
            placeholder="Ticker o Compañía..."
            class="w-full px-3 py-2 border border-slate-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 text-slate-900"
          />
        </div>
        <div class="flex-shrink-0 w-full md:w-auto">
          <label for="sort" class="block text-sm font-medium text-slate-700 mb-1">Ordenar Por</label>
          <select id="sort" v-model="currentSortKey" @change="handleSortChange" class="w-full px-3 py-2 border border-slate-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 bg-white text-slate-900">
            <option value="time">Fecha Evento</option>
            <option value="ticker">Ticker</option>
            <option value="company">Compañía</option>
            <option value="rating_to">Rating</option>
            <option value="target_to">Target</option>
            <option value="brokerage">Brokerage</option>
          </select>
        </div>
        <div class="flex-shrink-0 w-full md:w-auto">
          <label for="order" class="block text-sm font-medium text-slate-700 mb-1">Orden</label>
          <select id="order" v-model="currentSortOrder" @change="handleSortChange" class="w-full px-3 py-2 border border-slate-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 bg-white text-slate-900">
            <option value="desc">Descendente</option>
            <option value="asc">Ascendente</option>
          </select>
        </div>
      </div>

      <div v-if="stockStore.isLoading && stockStore.stocks.length === 0" class="text-center py-12">
        <svg class="mx-auto h-12 w-12 text-indigo-500 animate-spin" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <p class="text-lg text-slate-600 mt-4">Buscando datos del mercado...</p>
      </div>
      <div v-else-if="stockStore.error && !stockStore.isLoading" class="text-center py-10 bg-red-50 p-6 rounded-lg shadow-md">
        <p class="text-xl font-semibold text-red-700 mb-2">Houston, tenemos un problema...</p>
        <p class="text-md text-red-600">{{ stockStore.error }}</p>
        <p class="text-sm text-slate-500 mt-3">Intenta recargar o verifica tu conexión.</p>
      </div>
      <div v-else-if="!stockStore.stocks.length && !stockStore.isLoading && !stockStore.error" class="text-center py-12">
        <svg class="mx-auto h-12 w-12 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
          <path vector-effect="non-scaling-stroke" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 13h6m-3-3v6m-9 1V7a2 2 0 012-2h6l2 2h6a2 2 0 012 2v8a2 2 0 01-2 2H5a2 2 0 01-2-2z" />
        </svg>
        <p class="text-lg text-slate-600 mt-4">No hay acciones que coincidan. ¡Ajusta tus filtros!</p>
      </div>

      <div v-if="stockStore.stocks.length > 0 && !stockStore.error" class="overflow-x-auto bg-white shadow-lg rounded-lg relative">
          <div v-if="stockStore.isLoading"
                class="absolute inset-0 bg-white bg-opacity-60 flex items-center justify-center z-20 rounded-lg">
            <p class="text-lg text-indigo-600 font-semibold animate-pulse">Actualizando lista...</p>
          </div>
        <table class="min-w-full divide-y divide-slate-200">
          <thead class="bg-slate-100">
            <tr>
              <th @click="requestSort('ticker')" class="px-6 py-3 text-left text-xs font-bold text-slate-600 uppercase tracking-wider cursor-pointer hover:bg-slate-200 transition-colors">Ticker</th>
              <th @click="requestSort('company')" class="px-6 py-3 text-left text-xs font-bold text-slate-600 uppercase tracking-wider cursor-pointer hover:bg-slate-200 transition-colors">Compañía</th>
              <th @click="requestSort('brokerage')" class="px-6 py-3 text-left text-xs font-bold text-slate-600 uppercase tracking-wider cursor-pointer hover:bg-slate-200 transition-colors hidden md:table-cell">Brokerage</th>
              <th class="px-6 py-3 text-left text-xs font-bold text-slate-600 uppercase tracking-wider">Acción Broker</th>
              <th @click="requestSort('rating_to')" class="px-6 py-3 text-left text-xs font-bold text-slate-600 uppercase tracking-wider cursor-pointer hover:bg-slate-200 transition-colors">Rating</th>
              <th @click="requestSort('target_to')" class="px-6 py-3 text-left text-xs font-bold text-slate-600 uppercase tracking-wider cursor-pointer hover:bg-slate-200 transition-colors hidden sm:table-cell">Target</th>
              <th @click="requestSort('time')" class="px-6 py-3 text-left text-xs font-bold text-slate-600 uppercase tracking-wider cursor-pointer hover:bg-slate-200 transition-colors hidden lg:table-cell">Fecha Evento</th>
              <th class="px-6 py-3 text-left text-xs font-bold text-slate-600 uppercase tracking-wider">Detalles</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-slate-200">
            <tr v-for="stock in stockStore.stocks" :key="stock.ID || stock.ticker" class="hover:bg-indigo-50/50 transition-colors duration-150">
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-indigo-600">
                <router-link :to="{ name: 'StockDetail', params: { ticker: stock.ticker } }" class="hover:underline">
                  {{ stock.ticker }}
                </router-link>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-slate-800 truncate" :title="stock.company">{{ stock.company }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-slate-500 hidden md:table-cell truncate" :title="stock.brokerage">{{ stock.brokerage }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-slate-500">
                <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium" :class="getActionClass(stock.action)">
                  {{ getActionIcon(stock.action) }} <span class="ml-1">{{ stock.action }}</span>
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm font-semibold" :class="getRatingClass(stock.rating_to)">{{ stock.rating_to }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-slate-500 hidden sm:table-cell">
                {{ stock.target_to !== null && stock.target_to !== undefined ? '$' + stock.target_to.toFixed(2) : 'N/A' }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-slate-500 hidden lg:table-cell">{{ formatDate(stock.time) }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                <router-link :to="{ name: 'StockDetail', params: { ticker: stock.ticker } }" class="text-indigo-600 hover:text-indigo-800 hover:underline">
                  Ver
                </router-link>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <PaginationControls
        v-if="stockStore.totalPages > 1 && !stockStore.error && stockStore.stocks.length > 0"
        class="mt-6"
        :current-page="stockStore.currentPage"
        :total-pages="stockStore.totalPages"
        :page-size="stockStore.pageSize"
        :total-items="stockStore.totalItems"
        @page-changed="handlePageChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue';
import { useStockStore } from '../stores/stockStore';
import PaginationControls from '../components/common/PaginationControls.vue';

const stockStore = useStockStore();
const searchInputValue = ref(stockStore.searchQuery);
const currentSortKey = ref(stockStore.sortBy);
const currentSortOrder = ref(stockStore.sortOrder);
let debounceTimer: number | undefined;

onMounted(() => {
  stockStore.fetchRecommendations();
  stockStore.fetchStocks();
});

const formatDate = (dateString: string): string => {
  if (!dateString) return 'N/A';
  const date = new Date(dateString);
    if (isNaN(date.getTime())) return 'Fecha inválida';
  return date.toLocaleDateString('es-CO', { year: 'numeric', month: 'short', day: 'numeric' });
};
const debouncedSearch = (): void => {
  clearTimeout(debounceTimer);
  debounceTimer = window.setTimeout(() => {
    stockStore.setSearchQuery(searchInputValue.value);
  }, 500);
};
const requestSort = (sortKey: string): void => {
  let newOrder: 'asc' | 'desc' = 'asc';
  if (currentSortKey.value === sortKey) {
    newOrder = currentSortOrder.value === 'asc' ? 'desc' : 'asc';
  }
  currentSortKey.value = sortKey;
  currentSortOrder.value = newOrder;
  stockStore.setSort(currentSortKey.value, currentSortOrder.value);
};
const handleSortChange = (): void => { stockStore.setSort(currentSortKey.value, currentSortOrder.value) };
const handlePageChange = (newPage: number): void => { stockStore.changePage(newPage); };

watch(() => stockStore.searchQuery, (newQuery) => { searchInputValue.value = newQuery; });
watch(() => stockStore.sortBy, (newSortBy) => { currentSortKey.value = newSortBy; });
watch(() => stockStore.sortOrder, (newSortOrder) => { currentSortOrder.value = newSortOrder; });

const getActionIcon = (action: string): string => {
  const lowerAction = action.toLowerCase();
  if (lowerAction.includes("upgrade") || lowerAction.includes("raised")) return '⬆️';
  if (lowerAction.includes("downgrade") || lowerAction.includes("lowered")) return '⬇️';
  if (lowerAction.includes("initiate")) return '✨';
  if (lowerAction.includes("reiterate") || lowerAction.includes("maintain")) return '➡️';
  return 'ℹ️';
};
const getActionClass = (action: string): string => {
  const lowerAction = action.toLowerCase();
  if (lowerAction.includes("upgrade") || lowerAction.includes("raised")) return 'bg-green-100 text-green-700';
  if (lowerAction.includes("downgrade") || lowerAction.includes("lowered")) return 'bg-red-100 text-red-700';
  if (lowerAction.includes("initiate")) return 'bg-sky-100 text-sky-700';
  return 'bg-slate-100 text-slate-700';
};
const getRatingClass = (rating: string): string => {
    const lowerRating = rating.toLowerCase();
    if (["buy", "strong buy", "outperform", "overweight", "positive", "add", "accumulate"].includes(lowerRating)) return 'text-green-600';
    if (["sell", "strong sell", "underperform", "underweight", "negative"].includes(lowerRating)) return 'text-red-600';
    if (["hold", "neutral", "equal-weight", "market perform"].includes(lowerRating)) return 'text-amber-600';
    return 'text-slate-700';
};

const getReasonIcon = (reasonType: string): string => {
  switch (reasonType) {
      case 'POSITIVE_RATING': return '👍';
      case 'TARGET_INCREASED': return '📈';
      case 'TARGET_ATTRACTIVE': return '🎯';
      case 'BROKER_UPGRADE': return '🌟';
      case 'NEW_POSITIVE_COVERAGE': return '💡';
      case 'RECENT_EVENT': return '⏳';
      default: return 'ℹ️';
  }
};
</script>

<style scoped>
th {
  transition: background-color 0.2s ease-in-out;
  position: sticky;
  top: 0;
  background-color: #f3f4f6;
  z-index: 10;
}
th:hover {
  background-color: #e5e7eb;
}
.truncate {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 150px;
}
@media (min-width: 768px) {
  .truncate { max-width: 200px; }
}
@media (min-width: 1024px) {
  .truncate { max-width: 250px; }
}
</style>
