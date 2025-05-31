<template>
  <div v-if="totalPages > 0" class="flex items-center justify-between mt-6 px-4 py-3 sm:px-6 bg-white shadow rounded-md">
    <div class="flex-1 flex justify-between sm:hidden">
      <button
        @click="changePage(currentPage - 1)"
        :disabled="currentPage === 1"
        :class="['relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50', currentPage === 1 ? 'opacity-50 cursor-not-allowed' : '']"
      >
        Anterior
      </button>
      <span class="text-sm text-gray-700 px-2 py-2">
        Pág {{ currentPage }} de {{ totalPages }}
      </span>
      <button
        @click="changePage(currentPage + 1)"
        :disabled="currentPage === totalPages"
        :class="['ml-3 relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50', currentPage === totalPages ? 'opacity-50 cursor-not-allowed' : '']"
      >
        Siguiente
      </button>
    </div>

    <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
      <div>
        <p class="text-sm text-gray-700">
          Mostrando
          <span class="font-medium">{{ (currentPage - 1) * pageSize + 1 }}</span>
          a
          <span class="font-medium">{{ Math.min(currentPage * pageSize, totalItems) }}</span>
          de
          <span class="font-medium">{{ totalItems }}</span>
          resultados
        </p>
      </div>
      <div>
        <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px" aria-label="Pagination">
          <button
            @click="changePage(currentPage - 1)"
            :disabled="currentPage === 1"
            :class="['relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50', currentPage === 1 ? 'opacity-50 cursor-not-allowed' : '']"
            aria-label="Anterior"
          >
            <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
              <path fill-rule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clip-rule="evenodd" />
            </svg>
          </button>

          <template v-for="(pageNumber, index) in pageRange" :key="`page-${index}-${pageNumber}`">
            <button
              v-if="typeof pageNumber === 'number'"
              @click="changePage(pageNumber)"
              :aria-current="currentPage === pageNumber ? 'page' : undefined"
              :class="[
                'relative inline-flex items-center px-4 py-2 border text-sm font-medium',
                currentPage === pageNumber ? 'z-10 bg-indigo-50 border-indigo-500 text-indigo-600' : 'bg-white border-gray-300 text-gray-500 hover:bg-gray-50'
              ]"
            >
              {{ pageNumber }}
            </button>
            <span v-else class="relative inline-flex items-center px-4 py-2 border border-gray-300 bg-white text-sm font-medium text-gray-700">
              ...
            </span>
          </template>

          <button
            @click="changePage(currentPage + 1)"
            :disabled="currentPage === totalPages"
            :class="['relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50', currentPage === totalPages ? 'opacity-50 cursor-not-allowed' : '']"
            aria-label="Siguiente"
          >
            <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
              <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
            </svg>
          </button>
        </nav>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, defineProps, defineEmits } from 'vue';

const props = defineProps<{
  currentPage: number;
  totalPages: number;
  pageSize: number;
  totalItems: number;
}>();

const emit = defineEmits(['page-changed']);

const changePage = (page: number) => {
  if (page > 0 && page <= props.totalPages && page !== props.currentPage) {
    emit('page-changed', page);
  }
};

const pageRange = computed(() => {
  const delta = 1;
  const left = props.currentPage - delta;
  const right = props.currentPage + delta + 1;
  const range: (number)[] = [];
  const rangeWithDots: (number | string)[] = [];
  let l: number | undefined;

  for (let i = 1; i <= props.totalPages; i++) {
    if (i === 1 || i === props.totalPages || (i >= left && i < right)) {
      range.push(i);
    }
  }

  for (const i of range) {
    if (l != undefined) {
      if (i - l === 2) {
        rangeWithDots.push(l + 1);
      } else if (i - l !== 1) {
        rangeWithDots.push('...');
      }
    }

    rangeWithDots.push(i);
    l = i;
  }

  return rangeWithDots;
});
</script>
