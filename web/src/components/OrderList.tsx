import { useEffect, useState } from 'react';
import type { Order } from '../types';
import { api } from '../api/api';
import './OrderList.css';

export const OrderList = () => {
    const [orders, setOrders] = useState<Order[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const fetchOrders = () => {
        setLoading(true);
        api.getOrders()
            .then((data) => {
                setOrders(data);
                setError(null);
            })
            .catch((err) => {
                setError(err.message);
            })
            .finally(() => {
                setLoading(false);
            });
    };

    useEffect(() => {
        fetchOrders();
    }, []);

    if (loading) return <div className="order-list-loading">Loading past orders...</div>;
    if (error) return <div className="order-list-error">Error loading orders: {error}</div>;

    return (
        <div className="order-list-container">
            <div className="order-list-header">
                <h2>Past Orders ({orders.length})</h2>
                <button className="refresh-btn" onClick={fetchOrders}>Refresh</button>
            </div>

            {orders.length === 0 ? (
                <p className="no-orders-msg">You haven't placed any orders yet.</p>
            ) : (
                <div className="orders-grid">
                    {orders.map((order) => {
                        // Calculate total for this order
                        const orderTotal = order.items.reduce((acc, item) => {
                            const product = order.products.find(p => p.id === item.productId);
                            return acc + (product?.price || 0) * item.quantity;
                        }, 0);

                        return (
                            <div key={order.id} className="order-history-card">
                                <div className="order-card-header">
                                    <span className="order-id-badge">Order #{order.id.slice(0, 8)}</span>
                                    <span className="order-total-lbl">Total: ${orderTotal.toFixed(2)}</span>
                                </div>

                                <ul className="order-items-list">
                                    {order.items.map((item, idx) => {
                                        const product = order.products.find(p => p.id === item.productId);
                                        return (
                                            <li key={idx} className="order-item-history">
                                                <span className="qty-badge">{item.quantity}x</span>
                                                <span className="product-name">{product?.name || 'Unknown Item'}</span>
                                            </li>
                                        );
                                    })}
                                </ul>
                            </div>
                        );
                    })}
                </div>
            )}
        </div>
    );
};
