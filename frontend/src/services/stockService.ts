import axios, { type AxiosError } from 'axios';

const API_BASE_URL = import.meta.env.VITE_APP_API_BASE_URL || 'http://localhost:8080/api';

export interface Stock {
  ID?: number;
  CreatedAt?: string;
  UpdatedAt?: string;
  DeletedAt?: string | null;
  ticker: string;
  company: string;
  brokerage: string;
  action: string;
  rating_to: string;
  rating_from?: string | null;
  target_to?: number | null;
  target_from?: number | null;
  time: string;
}

export interface FetchParams {
  search?: string;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
  page?: number;
  pageSize?: number;
}

export interface PaginatedStocksResponse {
  stocks: Stock[];
  totalItems: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

export interface RecommendationReason {
  type: string;
  details: string;
}

export interface RecommendedStockView extends Stock {
  reasons: RecommendationReason[];
  score: number;
}

const stockService = {
  async getStocks(params: FetchParams = {}): Promise<PaginatedStocksResponse> {
    try {
      const response = await axios.get<PaginatedStocksResponse>(`${API_BASE_URL}/stocks`, { params });
      return response.data;
    } catch (error) {
      throw error;
    }
  },

  async getStockByTicker(ticker: string): Promise<Stock | null> {
    if (!ticker) {
      return Promise.resolve(null);
    }
    const fullUrl = `${API_BASE_URL}/stocks/${ticker}`;

    try {
      const response = await axios.get<Stock>(fullUrl, {
        headers: {
          'Cache-Control': 'no-cache',
          'Pragma': 'no-cache',
          'Expires': '0',
          'Accept': 'application/json',
        }
      });

      const contentType = response.headers['content-type'];
      if (contentType && contentType.includes('application/json')) {
        return response.data;
      } else {
        throw new Error(`Respuesta inesperada del servidor: se esperaba JSON pero se obtuvo ${contentType}`);
      }
    } catch (error) {
      const axiosError = error as AxiosError;
      if (axiosError.isAxiosError) {
        if (axiosError.response?.status === 404) {
          return null;
        }
      }
      throw error;
    }
  },

  async getRecommendations(): Promise<RecommendedStockView[]> {
    try {
      const response = await axios.get<{ recommendations: RecommendedStockView[] }>(`${API_BASE_URL}/stocks/recommendations`);
      return response.data.recommendations || [];
    } catch (error) {
      throw error;
    }
  },
};

export default stockService;
