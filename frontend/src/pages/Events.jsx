import React, { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";

// Helper: Format date with AM/PM
const formatDate = (dateStr) => {
  const date = new Date(dateStr);
  const day = String(date.getDate()).padStart(2, "0");
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const year = date.getFullYear();

  let hours = date.getHours();
  const minutes = String(date.getMinutes()).padStart(2, "0");
  const ampm = hours >= 12 ? "PM" : "AM";
  hours = hours % 12 || 12;

  return `${day}-${month}-${year} ${String(hours).padStart(2, "0")}:${minutes} ${ampm}`;
};

// Helper: Countdown timer
const getCountdown = (dateStr) => {
  const now = new Date();
  const eventDate = new Date(dateStr);
  const diff = eventDate - now;

  if (diff <= 0) return "Event Started";

  const days = Math.floor(diff / (1000 * 60 * 60 * 24));
  const hours = Math.floor((diff / (1000 * 60 * 60)) % 24);
  const minutes = Math.floor((diff / (1000 * 60)) % 60);

  return `${days}d ${hours}h ${minutes}m`;
};

// Helper: Check if event is closed (1 day before)
const isEventClosed = (eventDateStr) => {
  const now = new Date();
  const eventDate = new Date(eventDateStr);
  const cutoffDate = new Date(eventDate.getTime() - 24 * 60 * 60 * 1000);
  return now > cutoffDate; // true if event is closed
};

export default function EventsPage() {
  const [events, setEvents] = useState([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState("");
  const [sortBy, setSortBy] = useState("date");
  const [favorites, setFavorites] = useState(
    JSON.parse(localStorage.getItem("favorites") || "[]")
  );
  const [darkMode, setDarkMode] = useState(
    localStorage.getItem("darkMode") === "true"
  );
  const [showFavorites, setShowFavorites] = useState(false);

  const navigate = useNavigate();
  const user = JSON.parse(localStorage.getItem("user") || "{}");

  // Redirect if not logged in
  useEffect(() => {
    if (!user?.id) navigate("/signin");
  }, [navigate, user]);

  // Fetch events
  useEffect(() => {
    if (!user?.id) return;
    const fetchEvents = async () => {
      try {
        const res = await fetch("http://localhost:8001/show-events");
        if (!res.ok) throw new Error("Failed to fetch events");
        const data = await res.json();
        setEvents(data.showEvents || []);
      } catch (err) {
        console.error("Error fetching events:", err);
      } finally {
        setLoading(false);
      }
    };
    fetchEvents();
  }, [user?.id]);

  // Real-time countdown update every minute
  useEffect(() => {
    const interval = setInterval(() => setEvents([...events]), 60000);
    return () => clearInterval(interval);
  }, [events]);

  // Logout
  const handleLogout = () => {
    localStorage.removeItem("user");
    localStorage.removeItem("token");
    navigate("/signin");
  };

  // Toggle favorites
  const toggleFavorite = (id) => {
    const updated = favorites.includes(id)
      ? favorites.filter((f) => f !== id)
      : [...favorites, id];
    setFavorites(updated);
    localStorage.setItem("favorites", JSON.stringify(updated));
  };

  // Toggle dark mode
  const toggleDarkMode = () => {
    const newMode = !darkMode;
    setDarkMode(newMode);
    localStorage.setItem("darkMode", newMode);
  };

  // Filter + sort logic
  const filteredEvents = events
    .filter(
      (event) =>
        (!showFavorites || favorites.includes(event.id)) &&
        event.title.toLowerCase().includes(searchTerm.toLowerCase())
    )
    .sort((a, b) => {
      if (sortBy === "date") return new Date(a.date) - new Date(b.date);
      if (sortBy === "seats") return b.availableSeats - a.availableSeats;
      return 0;
    });

  if (!user?.id) return null;

  return (
    <div
      className={`min-h-screen font-sans transition-colors duration-300 ${
        darkMode ? "bg-gray-900 text-gray-100" : "bg-gray-50 text-gray-800"
      }`}
    >
      {/* Header */}
      <header
        className={`shadow-md py-4 mb-8 ${
          darkMode ? "bg-gray-800" : "bg-white"
        }`}
      >
        <div className="container mx-auto px-4 flex justify-between items-center">
          <h1 className="text-2xl font-bold">Event Ticket Booking</h1>
          <div className="flex items-center space-x-4">
            <button
              onClick={toggleDarkMode}
              className="px-3 py-2 rounded text-sm bg-gray-200 hover:bg-gray-300 dark:bg-gray-700 dark:hover:bg-gray-600 transition"
            >
              {darkMode ? "☀️ Light" : "🌙 Dark"}
            </button>
            <span className="font-medium">Hello, {user.name}</span>
            <button
              onClick={handleLogout}
              className="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600 transition"
            >
              Logout
            </button>
          </div>
        </div>
      </header>

      {/* Search + Favorites */}
      <div className="container mx-auto px-4 mb-8 flex flex-col md:flex-row items-center justify-between gap-4">
        <input
          type="text"
          placeholder="🔍 Search events..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="w-full md:w-1/2 p-2 rounded border focus:ring focus:outline-none text-gray-800"
        />
        <button
          onClick={() => setShowFavorites(!showFavorites)}
          className="px-3 py-2 rounded bg-yellow-400 hover:bg-yellow-500 transition"
        >
          {showFavorites ? "Show All Events" : "Show Favorites"}
        </button>
        <select
          value={sortBy}
          onChange={(e) => setSortBy(e.target.value)}
          className="p-2 rounded border text-gray-800"
        >
          <option value="date">Sort by Date</option>
          <option value="seats">Sort by Available Seats</option>
        </select>
      </div>

      {/* Events Section */}
      <main className="container mx-auto px-4">
        <h2 className="text-3xl font-bold text-center mb-8">Available Events</h2>

        {loading && <p className="text-center text-gray-500">Loading events...</p>}
        {!loading && filteredEvents.length === 0 && (
          <p className="text-center text-gray-500">No events found.</p>
        )}

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8 justify-items-center">
          {filteredEvents.map((event) => {
            const closed = isEventClosed(event.date);
            return (
              <div
                key={event.id}
                className={`relative shadow-lg rounded-xl p-6 w-full max-w-sm transition transform hover:scale-[1.02] ${
                  darkMode ? "bg-gray-800" : "bg-white"
                }`}
              >
                {event.imageUrl && (
                  <img
                    src={event.imageUrl}
                    alt={event.title}
                    className="rounded-lg mb-3 h-48 w-full object-cover"
                  />
                )}

                {/* Closed overlay */}
                {closed && (
                  <div className="absolute inset-0 bg-black bg-opacity-50 flex items-center justify-center rounded-lg z-10">
                    <span className="text-white font-bold text-lg">Event Closed</span>
                  </div>
                )}

                <div className="flex justify-between items-start mb-2 relative z-0">
                  <h3 className="text-xl font-semibold">{event.title}</h3>
                  <button
                    onClick={() => toggleFavorite(event.id)}
                    className={`text-xl ${
                      favorites.includes(event.id)
                        ? "text-red-500"
                        : "text-gray-400 hover:text-red-400"
                    } transition`}
                  >
                    ❤️
                  </button>
                </div>
                <p className="text-gray-400 mb-2">{event.description}</p>
                <p className="text-sm mb-1">📅 {formatDate(event.date)}</p>
                <p className="text-sm mb-1 text-green-500">⏱️ {getCountdown(event.date)}</p>
                <p className="text-sm mb-4">
                  🎟️ Seats: {event.availableSeats} / {event.totalSeats}
                </p>

                <Link to={`/events/${event.id}`}>
                  <button
                    disabled={closed}
                    className={`bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 w-full transition ${
                      closed ? "cursor-not-allowed opacity-50" : ""
                    }`}
                  >
                    View Details
                  </button>
                </Link>
              </div>
            );
          })}
        </div>
      </main>
    </div>
  );
}
