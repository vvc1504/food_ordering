import { useEffect, useState } from 'react';
import { api } from './api/api';
import type { Product } from './types';
import { ProductCard } from './components/ProductCard';
import { Cart } from './components/Cart';
import { OrderModal } from './components/OrderModal';
import { OrderList } from './components/OrderList';
import './App.css';

function App() {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [isIndexing, setIsIndexing] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [confirmedOrderId, setConfirmedOrderId] = useState<string | null>(null);

  const fetchProducts = () => {
    api.getProducts()
      .then((data) => {
        setProducts(data);
        setIsIndexing(false);
        setLoading(false);
        setError(null);
      })
      .catch((err) => {
        if (err.message === 'INDEXING') {
          setIsIndexing(true);
          setLoading(false);
          // Retry every 5 seconds
          setTimeout(fetchProducts, 5000);
        } else {
          setError(err.message);
          setLoading(false);
        }
      });
  };

  useEffect(() => {
    fetchProducts();
  }, []);

  return (
    <main className="app-container">
      <div className="container main-layout">
        <section className="products-section">
          <h1>Menu</h1>
          {loading && <p>Loading menu...</p>}
          {isIndexing && (
            <div className="indexing-box">
              <h3>Syncing Data</h3>
              <p>We are currently setting up the menu and special offers.</p>
              <p>This should only take a moment. The page will refresh automatically!</p>
              <div className="spinner"></div>
            </div>
          )}
          {error && <p className="error">{error}</p>}
          {!isIndexing && !loading && (
            <div className="products-grid">
              {products.map((product) => (
                <ProductCard key={product.id} product={product} />
              ))}
            </div>
          )}

          {!isIndexing && !loading && <OrderList />}
        </section>

        <aside className="cart-sidebar">
          <Cart onOrderConfirm={(id) => setConfirmedOrderId(id)} />
        </aside>
      </div>

      {confirmedOrderId && (
        <OrderModal
          orderId={confirmedOrderId}
          onClose={() => setConfirmedOrderId(null)}
        />
      )}
    </main>
  );
}

export default App;
