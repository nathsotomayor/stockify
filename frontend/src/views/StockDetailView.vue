<template>
  <div class="stock-detail-view"> <div class="max-w-4xl mx-auto">
      <button @click="goBack" class="mb-6 text-indigo-600 hover:text-indigo-700 font-medium flex items-center group transition-colors">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2 group-hover:-translate-x-1 transition-transform" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clip-rule="evenodd" />
        </svg>
        Volver a la lista
      </button>

      <div v-if="isLoading" class="text-center py-20">
        <svg class="mx-auto h-12 w-12 text-indigo-500 animate-spin" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <p class="text-lg text-slate-600 mt-4">Cargando detalles para <strong>{{ ticker }}</strong>...</p>
      </div>
      <div v-else-if="error && !isLoading" class="text-center py-10 bg-red-50 p-6 rounded-lg shadow-md">
        <p class="text-xl font-semibold text-red-700 mb-2">Oops! Algo salió mal.</p>
        <p class="text-md text-red-600">{{ error }}</p>
      </div>
      <div v-else-if="stock && !isLoading && !error" class="bg-white shadow-xl rounded-lg overflow-hidden">
        <div class="bg-gradient-to-r from-slate-700 to-slate-800 p-6 sm:p-8 text-white">
          <h1 class="text-3xl sm:text-4xl font-bold">{{ stock.company }}</h1>
          <p class="text-xl sm:text-2xl text-indigo-300 font-semibold mt-1">{{ stock.ticker }}</p>
        </div>

        <div class="p-6 sm:p-8 space-y-6">
          <section>
            <h2 class="text-xl font-semibold text-slate-700 mb-3 border-b border-slate-200 pb-2">Información del Evento</h2>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-x-6 gap-y-4">
              <DetailItem label="Brokerage" :value="stock.brokerage" />
              <DetailItem label="Fecha del Evento (API)" :value="formatDate(stock.time)" />
              <DetailItem label="Acción del Broker" :value="stock.action" class="sm:col-span-2"/>
            </div>
          </section>
          <section>
            <h2 class="text-xl font-semibold text-slate-700 mb-3 border-b border-slate-200 pb-2">Análisis de Rating</h2>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-x-6 gap-y-4">
              <DetailItem label="Rating Anterior" :value="stock.rating_from || 'N/A'" />
              <DetailItem label="Rating Actual" :value="stock.rating_to" />
            </div>
          </section>
          <section>
            <h2 class="text-xl font-semibold text-slate-700 mb-3 border-b border-slate-200 pb-2">Análisis de Precio Objetivo</h2>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-x-6 gap-y-4">
              <DetailItem label="Target Anterior" :value="stock.target_from !== null && stock.target_from !== undefined ? '$' + stock.target_from.toFixed(2) : 'N/A'" />
              <DetailItem label="Target Actual" :value="stock.target_to !== null && stock.target_to !== undefined ? '$' + stock.target_to.toFixed(2) : 'N/A'" />
            </div>
          </section>
          <section v-if="stock.CreatedAt || stock.UpdatedAt">
            <h2 class="text-xl font-semibold text-slate-700 mb-3 border-b border-slate-200 pb-2">Registro en Sistema</h2>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-x-6 gap-y-4">
                <DetailItem v-if="stock.CreatedAt" label="Registrado en BD" :value="formatDate(stock.CreatedAt)"/>
                <DetailItem v-if="stock.UpdatedAt" label="Última Actualización en BD" :value="formatDate(stock.UpdatedAt)"/>
            </div>
          </section>
        </div>
      </div>
      <div v-else-if="!stock && !isLoading && !error" class="text-center py-20">
        <svg class="mx-auto h-16 w-16 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
          <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607zM13.5 10.5h-6" />
        </svg>
        <p class="text-xl text-slate-600 mt-4">No se encontró información para el ticker <strong class="text-indigo-600">{{ ticker }}</strong>.</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, defineProps, watch } from 'vue';
import { useRouter } from 'vue-router';
import stockService, { type Stock } from '../services/stockService';
import DetailItem from '../components/common/DetailItem.vue';

const props = defineProps<{
  ticker: string
}>();

const stock = ref<Stock | null>(null);
const isLoading = ref(true);
const error = ref<string | null>(null);
const router = useRouter();

const fetchStockDetails = async (tickerValue: string) => {
  if (!tickerValue) {
    error.value = 'Ticker no proporcionado en la ruta.';
    isLoading.value = false;
    stock.value = null;
    return;
  }

  isLoading.value = true;
  error.value = null;
  stock.value = null;

  try {
    const data = await stockService.getStockByTicker(tickerValue);
    if (data) {
      stock.value = data;
    }
  } catch (err) {
    if (err instanceof Error) {
        error.value = err.message;
    } else {
        error.value = 'Un error desconocido ocurrió al cargar los detalles del stock.';
    }
    stock.value = null;
  } finally {
    isLoading.value = false;
  }
};

onMounted(() => {
  fetchStockDetails(props.ticker);
});

watch(() => props.ticker, (newTicker, oldTicker) => {
  if (newTicker && newTicker !== oldTicker) {
    fetchStockDetails(newTicker);
  } else if (newTicker && !stock.value && !error.value && !isLoading.value) {
    fetchStockDetails(newTicker);
  }
});

const formatDate = (dateString?: string): string => {
  if (!dateString) return 'N/A';
  try {
    const date = new Date(dateString);
    if (isNaN(date.getTime())) {
        return 'Fecha inválida';
    }
    return date.toLocaleString('es-CO', { year: 'numeric', month: 'long', day: 'numeric', hour: '2-digit', minute: '2-digit', hour12: true });
  } catch(e) {
    return 'Error de fecha';
  }
};

const goBack = (): void => {
  router.push({ name: 'StockList' });
};
</script>
