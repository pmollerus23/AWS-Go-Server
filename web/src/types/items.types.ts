export interface Item {
  id: number;
  name: string;
  description: string;
}

export interface CreateItemRequest {
  name: string;
  description: string;
}

export interface CreateItemResponse {
  id: number;
  name: string;
  description: string;
}
