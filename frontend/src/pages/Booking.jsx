import React from "react";
import { useLocation, useNavigate } from "react-router-dom";

export default function BookingPage() {
  const location = useLocation();
  const navigate = useNavigate();

  const { event, selectedSeats } = location.state || {};
  const user = JSON.parse(localStorage.getItem("user") || "{}");

  if (!event || !selectedSeats || !user?.id) {
    return (
      <div className="flex justify-center items-center h-screen bg-gray-900 text-red-500">
        No booking data found or user not logged in.
      </div>
    );
  }

  // Total Price Calculation
  const totalPrice = selectedSeats.reduce((sum, seat) => sum + seat.price, 0);
  const seatList = selectedSeats.map((s) => s.seatId).join(", ");

  const confirmBooking = async () => {
    try {
      const res = await fetch("http://localhost:8002/v1/bookings", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          event_id: event.id,
          user_id: user.id,
          seat_ids: selectedSeats.map((s) => s.seatId),
          totalCost: totalPrice, // pass totalPrice
        }),
      });

      if (!res.ok) throw new Error("Booking failed");

      const data = await res.json();
      const booking = data.booking;

      if (!booking) throw new Error("Invalid booking response");

      // Pass totalPrice to PaymentPage
      navigate(`/payment/${booking.id}`, {
        state: {
          booking: {
            ...booking,
            totalCost: totalPrice,
          },
        },
      });
    } catch (err) {
      console.error(err);
      alert(err.message || "An error occurred during booking.");
    }
  };

  return (
    <div className="min-h-screen bg-gray-900 flex items-center justify-center px-6 font-sans">
      <div className="max-w-3xl w-full bg-gray-800 p-10 rounded-lg shadow-lg text-white text-center font-[Inter]">
        <h1 className="text-3xl font-semibold mb-6">🎟️ Booking Details</h1>

        <div className="mb-4">
          <p className="text-lg font-medium">
            <strong>Event:</strong> {event.title}
          </p>
          <p className="text-gray-300">
            <strong>Date:</strong> {new Date(event.date).toLocaleString()}
          </p>
        </div>

        <div className="mb-4">
          <p className="font-medium">
            <strong>User:</strong> {user.name} ({user.email})
          </p>
        </div>

        <div className="mb-4">
          <p className="font-medium">
            <strong>Selected Seats:</strong> {seatList}
          </p>
          <p className="font-medium">
            <strong>Total Seats:</strong> {selectedSeats.length}
          </p>
        </div>

        <div className="mb-6 text-xl font-semibold">
          <p>
            💰 <strong>Total Price:</strong> ₹{totalPrice.toLocaleString()}
          </p>
        </div>

        <div className="text-center">
          <button
            onClick={confirmBooking}
            className="bg-green-600 text-white px-6 py-3 rounded-lg shadow-lg hover:shadow-xl hover:scale-105 transform transition font-semibold"
          >
            ✅ Confirm Booking
          </button>
        </div>
      </div>
    </div>
  );
}
