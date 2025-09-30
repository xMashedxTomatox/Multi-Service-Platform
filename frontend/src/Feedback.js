import React, { useState, useEffect } from "react";

function Feedback({ token }) {
  const [feedbacks, setFeedbacks] = useState([]);
  const [message, setMessage] = useState("");

  // Fetch feedbacks
  useEffect(() => {
    if (!token) return; // don’t fetch if not logged in

    const fetchFeedbacks = async () => {
      try {
        const res = await fetch(`${process.env.REACT_APP_API_URL}/feedback`, {
          headers: {
            "Authorization": `Bearer ${token}`,
            "Content-Type": "application/json",
          },
        });

        if (!res.ok) throw new Error("Failed to fetch feedbacks");

        const data = await res.json();
        setFeedbacks(Array.isArray(data) ? data : []); // ✅ ensure always an array
      } catch (err) {
        console.error("Error fetching feedbacks:", err);
        setFeedbacks([]); // fallback on error
      }
    };

    fetchFeedbacks();
  }, [token]);

  // Submit feedback
  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!token) {
      console.error("No token, cannot submit feedback");
      return;
    }

    try {
      const res = await fetch(`${process.env.REACT_APP_API_URL}/feedback`, {
        method: "POST",
        headers: {
          "Authorization": `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ message }),
      });

      if (!res.ok) throw new Error("Failed to submit feedback");

      const newFb = await res.json();
      setFeedbacks((prev) => [...prev, newFb]); // ✅ safe append
      setMessage("");
    } catch (err) {
      console.error("Error submitting feedback:", err);
    }
  };

  return (
    <div>
      <h2>Feedbacks</h2>
      <form onSubmit={handleSubmit}>
        <input
          placeholder="Write feedback..."
          value={message}
          onChange={(e) => setMessage(e.target.value)}
        />
        <button type="submit">Submit</button>
      </form>

      {feedbacks.length === 0 ? ( // ✅ safe conditional render
        <p>No feedback yet</p>
      ) : (
        <ul>
          {feedbacks.map((fb) => (
            <li key={fb.id}>
              {fb.message} (by user {fb.user_id})
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}

export default Feedback;
