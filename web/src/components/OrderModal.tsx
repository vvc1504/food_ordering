import './OrderModal.css';

interface OrderModalProps {
    orderId: string;
    onClose: () => void;
}

export const OrderModal = ({ orderId, onClose }: OrderModalProps) => {
    return (
        <div className="modal-overlay">
            <div className="modal-content">
                <div className="order-confirmed-icon">
                    <svg width="48" height="48" viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <path d="M24 0C10.7452 0 0 10.7452 0 24C0 37.2548 10.7452 48 24 48C37.2548 48 48 37.2548 48 24C48 10.7452 37.2548 0 24 0ZM36.7071 18.7071L22.7071 32.7071C22.3166 33.0976 21.6834 33.0976 21.2929 32.7071L11.2929 22.7071C10.9024 22.3166 10.9024 21.6834 11.2929 21.2929L12.7071 19.8787C13.0976 19.4882 13.7308 19.4882 14.1213 19.8787L22 27.7574L33.8787 15.8787C34.2692 15.4882 34.9024 15.4882 35.2929 15.8787L36.7071 17.2929C37.0976 17.6834 37.0976 18.3166 36.7071 18.7071Z" fill="#1EA94C" />
                    </svg>
                </div>
                <h1>Order Confirmed</h1>
                <p className="order-details-text">We hope you enjoy your food!</p>
                <div className="order-summary-box">
                    <div className="order-id-label">Order ID</div>
                    <div className="order-id-value">{orderId}</div>
                </div>
                <button className="start-new-order-btn" onClick={onClose}>
                    Start New Order
                </button>
            </div>
        </div>
    );
};
