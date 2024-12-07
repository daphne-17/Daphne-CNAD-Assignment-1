const reservationServiceUrl = "http://localhost:8082"; // Reservation service endpoint
const token = localStorage.getItem("authToken"); // Retrieve token from storage

// Fetch and display user bookings
async function fetchBookings() {
    const bookingList = document.getElementById("booking-list");
    const message = document.getElementById("action-message");
    try {
        const response = await fetch(`${reservationServiceUrl}/reservations/user`, {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${token}`,
            },
        });

        if (response.ok) {
            const data = await response.json();

            // Check if `data` is an array
            if (Array.isArray(data) && data.length > 0) {
                bookingList.innerHTML = data
                    .map(
                        (booking) => `
                        <div class="booking">
                            <p><strong>Reservation ID:</strong> ${booking.reservation_id}</p>
                            <p><strong>Vehicle ID:</strong> ${booking.vehicle_id}</p>
                            <p><strong>Start Time:</strong> ${booking.start_time}</p>
                            <p><strong>End Time:</strong> ${booking.end_time}</p>
                            <p><strong>Status:</strong> ${booking.status}</p>
                            <button class="btn modify-btn" onclick="modifyBooking('${booking.reservation_id}')">Modify</button>
                            <button class="btn cancel-btn" onclick="cancelBooking('${booking.reservation_id}')">Cancel</button>
                        </div>
                    `
                    )
                    .join("<hr>");
                message.textContent = ""; // Clear any previous action messages
            } else if (data.message) {
                // Handle case where no bookings exist
                bookingList.innerHTML = `<p>${data.message}</p>`;
                message.textContent = ""; // Clear any previous action messages
            } else {
                // Unexpected response structure
                bookingList.innerHTML = "<p>Error loading bookings.</p>";
                message.textContent = "Unexpected response format.";
                message.style.color = "red";
            }
        } else {
            bookingList.innerHTML = "<p>Error loading bookings.</p>";
            const error = await response.text();
            message.textContent = `Error: ${error}`;
            message.style.color = "red";
        }
    } catch (err) {
        bookingList.innerHTML = "<p>Error fetching bookings.</p>";
        message.textContent = "Error: Unable to fetch bookings.";
        message.style.color = "red";
        console.error(err);
    }
}


// Modify an existing booking
async function modifyBooking(reservationId) {
    localStorage.setItem("reservationId", reservationId); // Store the reservation ID
    window.location.href = "modifybooking.html"; // Redirect to the modify page
}


// Cancel an existing booking
async function cancelBooking(reservationId) {
    if (!confirm("Are you sure you want to cancel this reservation?")) {
        return;
    }

    const message = document.getElementById("action-message");

    try {
        const response = await fetch(`${reservationServiceUrl}/reservations/cancel`, {
            method: "DELETE",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${token}`,
            },
            body: JSON.stringify({ reservation_id: reservationId }),
        });

        if (response.ok) {
            message.textContent = "Reservation cancelled successfully!";
            message.style.color = "green";
            fetchBookings(); // Refresh the bookings
        } else {
            const error = await response.text();
            message.textContent = `Error cancelling reservation: ${error}`;
            message.style.color = "red";
        }
    } catch (err) {
        message.textContent = "Error: Unable to cancel booking.";
        message.style.color = "red";
        console.error(err);
    }
}

// Load bookings on page load
fetchBookings();
