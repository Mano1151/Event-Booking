import React, { useState, useEffect } from "react";
import { useLocation, useNavigate, useParams } from "react-router-dom";

export default function PaymentPage() {
  const location = useLocation();
  const navigate = useNavigate();
  const { id } = useParams();

  const user = JSON.parse(localStorage.getItem("user") || "{}");

  // Booking state: use location.state if available
  const [booking, setBooking] = useState(location.state?.booking || null);
  const [paymentMethod, setPaymentMethod] = useState("card");
  const [loading, setLoading] = useState(false);
  const [paymentSuccess, setPaymentSuccess] = useState(false);

  // Redirect if user not logged in
  useEffect(() => {
    if (!user?.id) navigate("/signin");
  }, [user, navigate]);

  // Fetch booking if not in state (page refresh)
  useEffect(() => {
    if (!booking && id) {
      fetch(`http://localhost:8002/v1/bookings/${id}`)
        .then((res) => res.json())
        .then((data) => {
          if (!data.booking) throw new Error("Booking not found");
          setBooking({
            ...data.booking,
            // Ensure totalCost exists
            totalCost:
              data.booking.totalCost ||
              data.booking.seats.reduce((sum, s) => sum + s.price, 0),
          });
        })
        .catch((err) => {
          console.error(err);
          alert("Booking not found. Redirecting to events.");
          navigate("/events");
        });
    }
  }, [booking, id, navigate]);

  const handlePayment = async () => {
    setLoading(true);
    try {
      const res = await fetch("http://localhost:8003/v1/payments", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          bookingId: booking.id,
          paymentMethod,
          userEmail: user.email,
        }),
      });

      const data = await res.json();

      if (res.ok && data.status === "PAID") {
        await fetch(`http://localhost:8002/v1/bookings/${booking.id}/confirm`, {
          method: "PUT",
          headers: { "Content-Type": "application/json" },
        });

        setPaymentSuccess(true);
        setTimeout(() => navigate("/events"), 2000);
      } else {
        alert("❌ Payment failed. Try again.");
      }
    } catch (err) {
      console.error(err);
      alert("Payment failed: " + err.message);
    } finally {
      setLoading(false);
    }
  };

  if (!booking) {
    return (
      <div className="flex justify-center items-center h-screen bg-gray-900 text-gray-400">
        Loading booking...
      </div>
    );
  }

  if (paymentSuccess) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-900 p-6">
        <div className="text-center">
          <div className="bg-green-600 w-24 h-24 rounded-full flex items-center justify-center animate-pop-in mx-auto shadow-lg">
            <svg
              className="w-12 h-12 text-white animate-check"
              fill="none"
              stroke="currentColor"
              strokeWidth={3}
              viewBox="0 0 24 24"
            >
              <path strokeLinecap="round" strokeLinejoin="round" d="M5 13l4 4L19 7" />
            </svg>
          </div>
          <h2 className="text-2xl text-green-400 font-semibold mt-6 animate-fade-in">
            Payment Successful!
          </h2>
          <p className="text-gray-400 mt-2">Redirecting to events...</p>
        </div>

        <style>{`
          .animate-pop-in { animation: popIn 0.4s ease-out forwards; }
          .animate-check { animation: checkMark 0.6s ease-out forwards; }
          .animate-fade-in { animation: fadeIn 1s ease-out forwards; }
          @keyframes popIn {0% { transform: scale(0); opacity: 0; } 100% { transform: scale(1); opacity: 1; }}
          @keyframes checkMark {0% { stroke-dasharray: 0 24; } 100% { stroke-dasharray: 24 0; }}
          @keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
        `}</style>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-900 p-6">
      <div className="bg-gray-800 text-white p-8 rounded-2xl shadow-2xl w-full max-w-md">
        <h2 className="text-3xl font-bold text-center mb-6">💳 Secure Payment</h2>

        <div className="space-y-3 mb-6 text-center">
          <p className="text-lg">
            Booking ID: <span className="font-semibold">{booking.id}</span>
          </p>
          <p className="text-lg">
            Total Amount: <span className="font-semibold text-green-400">₹{booking.totalCost}</span>
          </p>
          <p className="text-sm mt-2 text-gray-400">
            Paying as: <strong>{user.email}</strong>
          </p>
        </div>

        <div className="space-y-2">
          <label className="block font-medium text-sm mb-2 text-gray-300">Payment Method</label>
          <select
            value={paymentMethod}
            onChange={(e) => setPaymentMethod(e.target.value)}
            className="w-full p-3 rounded-lg bg-gray-700 text-white border border-gray-600 focus:ring-2 focus:ring-cyan-400"
          >
            <option value="card">💳 Card</option>
            <option value="upi">📱 UPI</option>
          </select>
        </div>

        <button
          onClick={handlePayment}
          disabled={loading}
          className={`mt-6 w-full py-3 rounded-xl text-lg font-semibold transition-all ${
            loading
              ? "bg-gray-600 cursor-not-allowed"
              : "bg-gradient-to-r from-green-500 to-emerald-600 hover:from-emerald-600 hover:to-green-500"
          }`}
        >
          {loading ? "Processing..." : "Pay Now"}
        </button>
      </div>
    </div>
  );
}
