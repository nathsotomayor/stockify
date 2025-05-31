import { setActivePinia, createPinia } from 'pinia';
import { useStockStore } from '../stockStore';
import stockService, {
    type PaginatedStocksResponse,
    type RecommendedStockView,
    type Stock,
    type RecommendationReason
} from '../../services/stockService';
import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest';

vi.mock('../../services/stockService', () => ({
  default: {
    getStocks: vi.fn(),
    getStockByTicker: vi.fn(),
    getRecommendations: vi.fn(),
  },
}));
describe('Stock Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    vi.mocked(stockService.getStocks).mockClear();
    vi.mocked(stockService.getRecommendations).mockClear();
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  it('initial state is correct', () => {
    const store = useStockStore();
    expect(store.stocks).toEqual([]);
    expect(store.isLoading).toBe(false);
    expect(store.error).toBeNull();
    expect(store.currentPage).toBe(1);
    expect(store.pageSize).toBe(10);
    expect(store.totalItems).toBe(0);
    expect(store.totalPages).toBe(0);
    expect(store.recommendations).toEqual([]);
    expect(store.isLoadingRecommendations).toBe(false);
    expect(store.recommendationsError).toBeNull();
  });

  describe('fetchStocks action', () => {
    it('fetches stocks successfully and updates state', async () => {
      const store = useStockStore();
      const mockResponse: PaginatedStocksResponse = {
        stocks: [{ ticker: 'AAPL', company: 'Apple', time: new Date().toISOString() } as Stock],
        totalItems: 1,
        page: 1,
        pageSize: 10,
        totalPages: 1,
      };
      vi.mocked(stockService.getStocks).mockResolvedValue(mockResponse);

      await store.fetchStocks();

      expect(store.isLoading).toBe(false);
      expect(store.stocks).toEqual(mockResponse.stocks);
      expect(store.totalItems).toBe(mockResponse.totalItems);
      expect(store.totalPages).toBe(mockResponse.totalPages);
      expect(store.currentPage).toBe(mockResponse.page);
      expect(store.error).toBeNull();
      expect(stockService.getStocks).toHaveBeenCalledTimes(1);
      expect(stockService.getStocks).toHaveBeenCalledWith({
        search: store.searchQuery || undefined,
        sortBy: store.sortBy || undefined,
        sortOrder: store.sortOrder || undefined,
        page: 1,
        pageSize: 10,
      });
    });

    it('handles error when fetching stocks', async () => {
      const store = useStockStore();
      vi.mocked(stockService.getStocks).mockRejectedValue(new Error('Network Error'));

      await store.fetchStocks();

      expect(store.isLoading).toBe(false);
      expect(store.stocks).toEqual([]);
      expect(store.error).toBe('Network Error');
      expect(store.totalItems).toBe(0);
      expect(store.totalPages).toBe(0);
    });
  });

  describe('fetchRecommendations action', () => {
    it('fetches recommendations successfully', async () => {
        const store = useStockStore();
        const mockReasons: RecommendationReason[] = [{ type: "TEST", details: "Test reason"}];
        const mockRecs: RecommendedStockView[] = [
            { ticker: 'REC1', company: 'Rec Co 1', reasons: mockReasons, score: 90, time: new Date().toISOString() } as RecommendedStockView,
        ];
        vi.mocked(stockService.getRecommendations).mockResolvedValue(mockRecs);

        await store.fetchRecommendations();

        expect(store.isLoadingRecommendations).toBe(false);
        expect(store.recommendations).toEqual(mockRecs);
        expect(store.recommendationsError).toBeNull();
        expect(stockService.getRecommendations).toHaveBeenCalledTimes(1);
    });

    it('handles error when fetching recommendations', async () => {
        const store = useStockStore();
        vi.mocked(stockService.getRecommendations).mockRejectedValue(new Error('API Error'));

        await store.fetchRecommendations();

        expect(store.isLoadingRecommendations).toBe(false);
        expect(store.recommendations).toEqual([]);
        expect(store.recommendationsError).toBe('API Error');
    });
  });

  it('setSearchQuery updates query, resets page, and fetches stocks', () => {
    const store = useStockStore();
    const fetchStocksSpy = vi.spyOn(store, 'fetchStocks');

    store.setSearchQuery('TEST');

    expect(store.searchQuery).toBe('TEST');
    expect(store.currentPage).toBe(1);
    expect(fetchStocksSpy).toHaveBeenCalledTimes(1);
  });

  it('setSort updates sort params, resets page, and fetches stocks', () => {
    const store = useStockStore();
    const fetchStocksSpy = vi.spyOn(store, 'fetchStocks');

    store.setSort('company', 'desc');

    expect(store.sortBy).toBe('company');
    expect(store.sortOrder).toBe('desc');
    expect(store.currentPage).toBe(1);
    expect(fetchStocksSpy).toHaveBeenCalledTimes(1);
  });

  it('changePage updates page and fetches stocks if page is valid', () => {
    const store = useStockStore();
    store.totalPages = 5;
    const fetchStocksSpy = vi.spyOn(store, 'fetchStocks');

    store.changePage(3);
    expect(store.currentPage).toBe(3);
    expect(fetchStocksSpy).toHaveBeenCalledTimes(1);

    store.changePage(0);
    expect(store.currentPage).toBe(3);
    expect(fetchStocksSpy).toHaveBeenCalledTimes(1);

    store.changePage(6);
    expect(store.currentPage).toBe(3);
    expect(fetchStocksSpy).toHaveBeenCalledTimes(1);

    store.changePage(3);
    expect(store.currentPage).toBe(3);
    expect(fetchStocksSpy).toHaveBeenCalledTimes(1);
  });
});
