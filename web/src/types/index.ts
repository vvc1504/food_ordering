export interface Product {
  id: string;
  name: string;
  category: string;
  price: number;
  image?: {
    thumbnail: string;
    mobile: string;
    tablet: string;
    desktop: string;
  };
}

export interface CartItem extends Product {
  quantity: number;
}

export interface OrderReq {
  items: {
    productId: string;
    quantity: number;
  }[];
  couponCode?: string;
}

export interface OrderResponse {
  orderId: string;
}

export interface OrderItem {
  productId: string;
  quantity: number;
}

export interface Order {
  id: string;
  items: OrderItem[];
  products: Product[];
}
