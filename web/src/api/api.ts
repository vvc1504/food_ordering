import type { Product, OrderReq, OrderResponse, Order } from '../types';

const API_BASE_URL = 'http://localhost:8080';

export const api = {
    getProducts: async (): Promise<Product[]> => {
        const response = await fetch(`${API_BASE_URL}/product`);
        if (response.status === 503) {
            throw new Error('INDEXING');
        }
        if (!response.ok) throw new Error('Failed to fetch products');
        return response.json();
    },

    getProduct: async (id: string): Promise<Product> => {
        const response = await fetch(`${API_BASE_URL}/product/${id}`);
        if (!response.ok) throw new Error('Failed to fetch product');
        return response.json();
    },

    placeOrder: async (order: OrderReq): Promise<OrderResponse> => {
        const response = await fetch(`${API_BASE_URL}/order`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'api_key': 'apitest',
            },
            body: JSON.stringify(order),
        });
        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText || 'Failed to place order');
        }
        return response.json();
    },

    getOrders: async (): Promise<Order[]> => {
        const response = await fetch(`${API_BASE_URL}/orders`);
        if (!response.ok) throw new Error('Failed to fetch orders');
        return response.json();
    }
};
