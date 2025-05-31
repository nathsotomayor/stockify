import { defineStore } from 'pinia';
import stockService, {
    type Stock,
    type FetchParams,
    type RecommendedStockView
} from '../services/stockService';

export interface StockState {
  stocks: Stock[];
  isLoading: boolean;
  error: string | null;
  searchQuery: string;
  sortBy: string;
  sortOrder: 'asc' | 'desc';

  currentPage: number;
  pageSize: number;
  totalItems: number;
  totalPages: number;

  recommendations: RecommendedStockView[];
  isLoadingRecommendations: boolean;
  recommendationsError: string | null;
}

export const useStockStore = defineStore('stock', {
  state: (): StockState => ({
    stocks: [],
    isLoading: false,
    error: null,
    searchQuery: '',
    sortBy: 'time',
    sortOrder: 'desc',

    currentPage: 1,
    pageSize: 10,
    totalItems: 0,
    totalPages: 0,

    recommendations: [],
    isLoadingRecommendations: false,
    recommendationsError: null,
  }),
  actions: {
    async fetchStocks() {
      this.isLoading = true;
      this.error = null;
      try {
        const params: FetchParams = {
          search: this.searchQuery || undefined,
          sortBy: this.sortBy || undefined,
          sortOrder: this.sortOrder || undefined,
          page: this.currentPage,
          pageSize: this.pageSize,
        };
        const responseData = await stockService.getStocks(params);
        this.stocks = responseData.stocks;
        this.totalItems = responseData.totalItems;
        this.totalPages = responseData.totalPages;
        this.currentPage = responseData.page;
        if (this.currentPage > this.totalPages && this.totalPages > 0) {
            this.currentPage = this.totalPages;
        } else if (this.totalPages === 0) {
            this.currentPage = 1;
        }

      } catch (err) {
        this.error = err instanceof Error ? err.message : 'Ocurrió un error desconocido al obtener las acciones.';
        this.stocks = [];
        this.totalItems = 0;
        this.totalPages = 0;
      } finally {
        this.isLoading = false;
      }
    },
    setSearchQuery(query: string) {
      this.searchQuery = query;
      this.currentPage = 1;
      this.fetchStocks();
    },
    setSort(sortBy: string, sortOrder: 'asc' | 'desc') {
      this.sortBy = sortBy;
      this.sortOrder = sortOrder;
      this.currentPage = 1;
      this.fetchStocks();
    },
    changePage(page: number) {
        if (page > 0 && page <= this.totalPages && page !== this.currentPage) {
            this.currentPage = page;
            this.fetchStocks();
        } else if (page === 1 && this.currentPage !== 1) {
            this.currentPage = page;
            this.fetchStocks();
        }
    },
    changePageSize(size: number) {
        if (size > 0 && size !== this.pageSize) {
            this.pageSize = size;
            this.currentPage = 1;
            this.fetchStocks();
        }
    },
    async fetchRecommendations() {
      this.isLoadingRecommendations = true;
      this.recommendationsError = null;
      try {
        const data = await stockService.getRecommendations();
        this.recommendations = data;
      } catch (err) {
        this.recommendationsError = err instanceof Error ? err.message : 'Ocurrió un error desconocido al obtener las recomendaciones.';
        this.recommendations = [];
      } finally {
        this.isLoadingRecommendations = false;
      }
    },
  },
});
