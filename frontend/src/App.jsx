import React from "react";
import { Routes, Route } from "react-router-dom";
import SignUpPage from "./pages/SignUp.jsx";
import SignInPage from "./pages/SignIn.jsx";
import EventsPage from "./pages/Events.jsx";
import EventDetailPage from "./pages/EventDetail.jsx";
import BookingPage from "./pages/Booking.jsx";
import PaymentPage from "./pages/Payment.jsx";

function App() {
  return (
    <Routes>
      <Route path="/" element={<SignInPage />} />
      <Route path="/signin" element={<SignInPage />} />
      <Route path="/signup" element={<SignUpPage />} />
      <Route path="/events" element={<EventsPage />} />
      <Route path="/events/:id" element={<EventDetailPage />} />
      <Route path="/booking/:id" element={<BookingPage />} />
      {/* ✅ Payment route fixed */}
      <Route path="/payment/:id" element={<PaymentPage />} />
    </Routes>
  );
}

export default App;
