import { apiClient } from './client';
import type { Item, CreateItemRequest, CreateItemResponse } from '../types/items.types';

export const itemsApi = {
  getAll: async (): Promise<Item[]> => {
    return apiClient.get<Item[]>('/api/v1/items');
  },

  create: async (data: CreateItemRequest): Promise<CreateItemResponse> => {
    return apiClient.post<CreateItemResponse>('/api/v1/items', data);
  },
};
