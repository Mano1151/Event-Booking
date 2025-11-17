import React, { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";

// ---------- Error Boundary ----------
class ErrorBoundary extends React.Component {
  constructor(props) {
    super(props);
    this.state = { hasError: false, error: null };
  }
  static getDerivedStateFromError(error) {
    return { hasError: true, error };
  }
  render() {
    if (this.state.hasError) {
      return (
        <div className="flex justify-center items-center h-screen text-red-500 bg-gray-900">
          <div className="text-center">
            <h1 className="text-2xl font-bold mb-4">Something went wrong!</h1>
            <p>{this.state.error?.message || "Unknown error occurred."}</p>
          </div>
        </div>
      );
    }
    return this.props.children;
  }
}

export default function EventDetailPageWrapper() {
  return (
    <ErrorBoundary>
      <EventDetailPage />
    </ErrorBoundary>
  );
}

function EventDetailPage() {
  const { id } = useParams();
  const navigate = useNavigate();

  const user = JSON.parse(localStorage.getItem("user") || "{}");
  useEffect(() => {
    if (!user?.id) navigate("/signin");
  }, [navigate, user?.id]);

  const [event, setEvent] = useState(null);
  const [selectedSeats, setSelectedSeats] = useState([]);
  const [lockedSeats, setLockedSeats] = useState([]);
  const [bookedSeats, setBookedSeats] = useState([]);
  const [loading, setLoading] = useState(true);

  const eventId = Number(id);

  const cleanSeatIds = (arr) =>
    Array.from(new Set((arr || []).map(String).filter((s) => s !== "")));

  // -------- Fetch Event Data --------
  useEffect(() => {
    if (isNaN(eventId)) return;
    const fetchEvent = async () => {
      try {
        const res = await fetch(`http://localhost:8001/show-events/${eventId}`);
        if (!res.ok) throw new Error(`Event fetch failed (${res.status})`);
        const data = await res.json();

        const eventData =
          data.showEvent ||
          data.show_event ||
          data.event ||
          data.data ||
          (data.id ? data : null);

        if (!eventData) throw new Error("Event not found");

        setEvent({
          ...eventData,
          total_seats: eventData.total_seats ?? eventData.totalSeats,
          price_per_seat: eventData.price_per_seat ?? eventData.pricePerSeat,
        });
      } catch (err) {
        console.error("Event Fetch Error:", err);
      } finally {
        setLoading(false);
      }
    };
    fetchEvent();
  }, [eventId]);

  // -------- Fetch Seats --------
  const fetchSeats = async () => {
    try {
      const [bookedRes, lockedRes] = await Promise.all([
        fetch(`http://localhost:8002/events/${eventId}/booked-seats`),
        fetch(`http://localhost:8002/events/${eventId}/locked-seats`),
      ]);

      if (!bookedRes.ok || !lockedRes.ok) throw new Error("Failed to fetch seat data");

      const bookedData = await bookedRes.json();
      const lockedData = await lockedRes.json();

      setBookedSeats(cleanSeatIds(bookedData.seatIds || bookedData.seat_ids));
      setLockedSeats(cleanSeatIds(lockedData.seatIds || lockedData.seat_ids));
    } catch (err) {
      console.error("Seat Fetch Error:", err);
      setBookedSeats([]);
      setLockedSeats([]);
    }
  };

  useEffect(() => {
    fetchSeats();
    const interval = setInterval(fetchSeats, 5000);
    return () => clearInterval(interval);
  }, [eventId]);

  // -------- Seat Selection Logic (Fixed) --------
  const toggleSeat = (seatId, tier) => {
    if (bookedSeats.includes(seatId) || lockedSeats.includes(seatId)) return;

    setSelectedSeats((prev) => {
      const exists = prev.find((s) => s.seatId === seatId);
      if (exists) {
        // Remove seat
        return prev.filter((s) => s.seatId !== seatId);
      }
      // Add seat
      const updated = [...prev, { seatId, tierName: tier.name, price: tier.price }];
      // Ensure no duplicates
      return Array.from(new Map(updated.map((s) => [s.seatId, s])).values());
    });
  };

  // -------- Proceed to Booking --------
  const proceedBooking = () => {
    if (!selectedSeats.length) {
      alert("Select at least one seat to proceed.");
      return;
    }

    navigate(`/booking/${eventId}`, {
      state: {
        event: {
          id: event.id,
          title: event.title,
          date: event.date,
          price_per_seat: event.price_per_seat,
          total_seats: event.total_seats,
        },
        selectedSeats,
        user,
      },
    });
  };

  // -------- Loading / Error States --------
  if (loading)
    return (
      <div className="flex justify-center items-center h-screen bg-gray-900 text-white">
        Loading event details...
      </div>
    );

  if (!event)
    return (
      <div className="flex justify-center items-center h-screen bg-gray-900 text-red-500">
        Event not found.
      </div>
    );

  // -------- Seat Grid & Tier Setup --------
  const seatsPerRow = 10;
  const totalRows = Math.ceil(event.total_seats / seatsPerRow);
  const getSeatLabel = (i) =>
    String.fromCharCode(65 + Math.floor(i / seatsPerRow)) + ((i % seatsPerRow) + 1);

  const silverPrice = event.price_per_seat || 250;
  const seatTiers = [
    { name: "Platinum", rows: 2, price: silverPrice * 3 },
    { name: "Gold", rows: 3, price: silverPrice * 2 },
    { name: "Silver", rows: totalRows - 5, price: silverPrice },
  ];

  const getSeatTier = (rowIndex) => {
    let total = 0;
    for (let tier of seatTiers) {
      if (rowIndex < total + tier.rows) return tier;
      total += tier.rows;
    }
    return seatTiers[seatTiers.length - 1];
  };

  // -------- Render --------
  return (
    <div className="max-w-6xl mx-auto px-6 py-10 font-sans bg-gray-900 text-white min-h-screen">
      {/* Event Header */}
      <div className="bg-gray-800 text-white rounded-lg p-6 mb-8 shadow-lg">
        <h1 className="text-4xl font-bold mb-2">{event.title}</h1>
        <p className="text-sm opacity-80 mb-2">{event.description}</p>
        <div className="flex flex-wrap gap-6 text-sm opacity-80">
          <p>📅 <strong>Date:</strong> {new Date(event.date).toLocaleString()}</p>
          <p>💰 <strong>Price:</strong> Starting from ₹{silverPrice}</p>
          <p>🪑 <strong>Total Seats:</strong> {event.total_seats}</p>
        </div>
      </div>

      {/* Seat Grid */}
      <div className="flex flex-col items-center space-y-3 mb-10">
        {Array.from({ length: totalRows }).map((_, rowIndex) => {
          const tier = getSeatTier(rowIndex);

          const tierLabel = (() => {
            let sum = 0;
            for (let t of seatTiers) {
              if (rowIndex === sum) return t.name;
              sum += t.rows;
            }
            return null;
          })();

          return (
            <React.Fragment key={rowIndex}>
              {tierLabel && (
                <div className="text-sm font-semibold text-white my-2">
                  --- {tierLabel} ---
                </div>
              )}
              <div className="flex space-x-2">
                {Array.from({ length: seatsPerRow }).map((_, colIndex) => {
                  const seatIndex = rowIndex * seatsPerRow + colIndex;
                  if (seatIndex >= event.total_seats) return null;

                  const label = getSeatLabel(seatIndex);
                  const seatId = label;

                  const isBooked = bookedSeats.includes(seatId);
                  const isLocked = lockedSeats.includes(seatId);
                  const isSelected = selectedSeats.some((s) => s.seatId === seatId);

                  const seatStyles = isBooked
                    ? "bg-red-800 text-white border border-red-900"
                    : isLocked
                    ? "bg-orange-600 text-white border border-orange-700"
                    : isSelected
                    ? "bg-green-500 text-white border border-green-600"
                    : "bg-gray-700 text-white hover:scale-110 hover:brightness-110 border border-gray-800";

                  return (
                    <div
                      key={seatId}
                      onClick={() => toggleSeat(seatId, tier)}
                      className={`w-12 h-12 flex items-center justify-center rounded-xl text-sm font-semibold cursor-pointer transition-all duration-300 ${seatStyles}`}
                      title={`${label} - ${tier.name} ₹${tier.price}`}
                    >
                      {label}
                    </div>
                  );
                })}
              </div>
            </React.Fragment>
          );
        })}
      </div>

      {/* Legend */}
      <div className="flex justify-center space-x-6 mb-10 text-sm">
        <Legend color="bg-green-500" label="Selected" />
        <Legend color="bg-orange-600" label="Locked" />
        <Legend color="bg-red-800" label="Booked" />
      </div>

      {/* Proceed Button */}
      <div className="text-center">
        <button
          onClick={proceedBooking}
          className="bg-gradient-to-r from-purple-600 to-blue-600 text-white px-8 py-3 rounded-lg shadow-lg hover:shadow-xl hover:scale-105 transform transition"
        >
          🎟️ Proceed to Booking ({selectedSeats.length} seat
          {selectedSeats.length > 1 ? "s" : ""})
        </button>
      </div>
    </div>
  );
}

// ---------- Legend ----------
function Legend({ color, label }) {
  return (
    <div className="flex items-center space-x-2">
      <div className={`w-5 h-5 rounded-xl ${color}`}></div>
      <span>{label}</span>
    </div>
  );
}
