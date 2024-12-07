const reservationServiceUrl = "http://localhost:8082"; // Reservation service endpoint
const token = localStorage.getItem("authToken"); // Retrieve token from storage

document.addEventListener("DOMContentLoaded", () => {
    const reservationId = localStorage.getItem("reservationId"); // Retrieve reservation ID
    const reservationInput = document.getElementById("reservation-id");
    const message = document.getElementById("modify-message");
    const form = document.getElementById("modify-reservation-form");

    if (!reservationId) {
        message.textContent = "Reservation ID not found. Please try again.";
        message.style.color = "red";
        return;
    }

    reservationInput.value = reservationId;

    form.addEventListener("submit", async (e) => {
        e.preventDefault();

        const newStartTime = document.getElementById("new-start-time").value;
        const newEndTime = document.getElementById("new-end-time").value;

        // Ensure datetime strings are in the correct format
        const formattedStartTime = `${newStartTime}:00`;
        const formattedEndTime = `${newEndTime}:00`;

        console.log("Formatted times for modification:", formattedStartTime, formattedEndTime);

        try {
            const response = await fetch(`${reservationServiceUrl}/reservations/modify`, {
                method: "PUT",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${token}`,
                },
                body: JSON.stringify({
                    reservation_id: reservationId,
                    new_start_time: formattedStartTime,
                    new_end_time: formattedEndTime,
                }),
            });

            if (response.ok) {
                message.textContent = "Reservation modified successfully!";
                message.style.color = "green";
                localStorage.removeItem("reservationId"); // Clear stored ID
            } else {
                const error = await response.text();
                message.textContent = `Error: ${error}`;
                message.style.color = "red";
            }
        } catch (error) {
            console.error("Error modifying reservation:", error);
            message.textContent = "Error: Unable to modify reservation.";
            message.style.color = "red";
        }
    });
});
