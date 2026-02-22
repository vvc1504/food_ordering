import { useCart } from '../context/CartContext';
import type { Product } from '../types';
import './ProductCard.css';

// Import assets
import chickenWaffle from '../assets/products/chicken_waffle.png';
import beefBurger from '../assets/products/beef_burger.png';
import veggiePizza from '../assets/products/veggie_pizza.png';
import classicFries from '../assets/products/classic_fries.png';
import icedTea from '../assets/products/iced_tea.png';

const imageMap: Record<string, string> = {
    '1': chickenWaffle,
    '2': beefBurger,
    '3': veggiePizza,
    '4': classicFries,
    '5': icedTea,
};

interface ProductCardProps {
    product: Product;
}

export const ProductCard = ({ product }: ProductCardProps) => {
    const { cart, addToCart, updateQuantity } = useCart();
    const cartItem = cart.find((item) => item.id === product.id);
    const quantity = cartItem?.quantity || 0;

    const image = imageMap[product.id] || chickenWaffle;

    return (
        <div className={`product-card ${quantity > 0 ? 'active' : ''}`}>
            <div className="product-image-container">
                <img src={image} alt={product.name} className="product-image" />
                {quantity === 0 ? (
                    <button className="add-to-cart-btn" onClick={() => addToCart(product)}>
                        <img src="/assets/icon-add-to-cart.svg" alt="" />
                        Add to Cart
                    </button>
                ) : (
                    <div className="quantity-control">
                        <button onClick={() => updateQuantity(product.id, -1)} aria-label="Decrease quantity">
                            <span className="icon-minus">-</span>
                        </button>
                        <span>{quantity}</span>
                        <button onClick={() => updateQuantity(product.id, 1)} aria-label="Increase quantity">
                            <span className="icon-plus">+</span>
                        </button>
                    </div>
                )}
            </div>
            <div className="product-info">
                <span className="product-category">{product.category}</span>
                <h3 className="product-name">{product.name}</h3>
                <span className="product-price">${product.price.toFixed(2)}</span>
            </div>
        </div>
    );
};
