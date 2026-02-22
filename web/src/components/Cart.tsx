import { useState } from 'react';
import { useCart } from '../context/CartContext';
import { api } from '../api/api';
import './Cart.css';

interface CartProps {
    onOrderConfirm: (orderId: string) => void;
}

export const Cart = ({ onOrderConfirm }: CartProps) => {
    const { cart, removeFromCart, clearCart } = useCart();
    const [coupon, setCoupon] = useState('');
    const [appliedCoupon, setAppliedCoupon] = useState<string | null>(null);
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const subtotal = cart.reduce((acc, item) => acc + item.price * item.quantity, 0);

    const totalItemsCount = cart.reduce((acc, item) => acc + item.quantity, 0);

    let discount = 0;
    if (appliedCoupon === 'HAPPYHOURS') {
        discount = subtotal * 0.18;
    } else if (appliedCoupon === 'BUYGETONE' && totalItemsCount >= 2) {
        // Flatten the cart into an array of individual item prices
        const allPrices: number[] = [];
        cart.forEach(item => {
            for (let i = 0; i < item.quantity; i++) {
                allPrices.push(item.price);
            }
        });

        // Sort the prices from lowest to highest
        allPrices.sort((a, b) => a - b);

        // Calculate how many items are free (1 for every 2 items bought)
        const freeItemsCount = Math.floor(allPrices.length / 2);

        // The discount is the sum of the cheapest `freeItemsCount` items
        discount = allPrices.slice(0, freeItemsCount).reduce((sum, price) => sum + price, 0);
    }

    const total = subtotal - discount;

    const handleApplyCoupon = () => {
        const trimmed = coupon.trim();
        if (!trimmed) {
            setAppliedCoupon(null);
            setError(null);
        } else {
            setAppliedCoupon(trimmed);
            setError(null);
        }
    };

    const handleConfirmOrder = async () => {
        setIsSubmitting(true);
        setError(null);
        try {
            const res = await api.placeOrder({
                items: cart.map((item) => ({ productId: item.id, quantity: item.quantity })),
                couponCode: appliedCoupon || undefined,
            });
            onOrderConfirm(res.orderId);
            clearCart();
        } catch (err: any) {
            setError(err.message || 'Failed to place order');
        } finally {
            setIsSubmitting(false);
        }
    };

    if (cart.length === 0) {
        return (
            <div className="cart empty">
                <h2>Your Cart (0)</h2>
                <div className="empty-state">
                    <img src="/assets/illustration-empty-cart.svg" alt="" />
                    <p>Your added items will appear here</p>
                </div>
            </div>
        );
    }

    return (
        <div className="cart">
            <h2>Your Cart ({cart.reduce((a, b) => a + b.quantity, 0)})</h2>
            <ul className="cart-items">
                {cart.map((item) => (
                    <li key={item.id} className="cart-item">
                        <div className="item-details">
                            <span className="item-name">{item.name}</span>
                            <div className="item-meta">
                                <span className="item-qty">{item.quantity}x</span>
                                <span className="item-unit-price">@ ${item.price.toFixed(2)}</span>
                                <span className="item-total-price-small">${(item.price * item.quantity).toFixed(2)}</span>
                            </div>
                        </div>
                        <button className="remove-item" onClick={() => removeFromCart(item.id)} aria-label="Remove item">
                            <div className="remove-icon">×</div>
                        </button>
                    </li>
                ))}
            </ul>

            <div className="cart-summary">
                <div className="summary-row">
                    <span>Order Total</span>
                    <span className="total-price">${total.toFixed(2)}</span>
                </div>
            </div>

            <div className="coupon-section">
                <div className="coupon-input-group">
                    <input
                        type="text"
                        placeholder="Coupon Code"
                        value={coupon}
                        onChange={(e) => setCoupon(e.target.value.toUpperCase())}
                    />
                    <button onClick={handleApplyCoupon}>Apply</button>
                </div>
                {appliedCoupon && <p className="coupon-success">Applied: {appliedCoupon}</p>}
                {error && <p className="coupon-error">{error}</p>}
            </div>

            <div className="carbon-neutral">
                <img src="/assets/icon-carbon-neutral.svg" alt="" />
                <span>This is a <strong>carbon-neutral</strong> delivery</span>
            </div>

            <button className="confirm-btn" disabled={isSubmitting} onClick={handleConfirmOrder}>
                {isSubmitting ? 'Placing Order...' : 'Confirm Order'}
            </button>
        </div>
    );
};
