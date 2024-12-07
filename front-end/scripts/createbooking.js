const reservationServiceUrl = "http://localhost:8082"; // Reservation service endpoint
const billingServiceUrl = "http://localhost:8083"; // Billing service endpoint
const token = localStorage.getItem("authToken");

document.addEventListener("DOMContentLoaded", () => {
    const bookingForm = document.getElementById("create-booking-form");
    const calculateCostBtn = document.getElementById("calculate-cost-btn");
    const billingSection = document.getElementById("billing-section");
    const estimatedCostElement = document.getElementById("estimated-cost");
    const confirmPaymentBtn = document.getElementById("confirm-payment-btn");
    const message = document.getElementById("reservation-message");

    let reservationDetails = {};

    // Step 1: Calculate Cost
    calculateCostBtn.addEventListener("click", async () => {
        try {
            // Validate inputs
            const vehicleId = document.getElementById("vehicle-id").value;
            const startTime = document.getElementById("start-time").value;
            const endTime = document.getElementById("end-time").value;
            const promoCode = document.getElementById("promo-code").value;

            if (!vehicleId || !startTime || !endTime) {
                throw new Error("All fields are required.");
            }

            // Format times
            const formattedStartTime = `${startTime}:00`;
            const formattedEndTime = `${endTime}:00`;

            reservationDetails = {
                vehicle_id: vehicleId,
                start_time: formattedStartTime,
                end_time: formattedEndTime,
                promo_code: promoCode,
            };

            // Fetch estimated cost
            const response = await fetch(`${billingServiceUrl}/billing/calculate`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${token}`,
                },
                body: JSON.stringify(reservationDetails),
            });

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(errorText);
            }

            const data = await response.json();
            estimatedCostElement.textContent = `Estimated Cost: $${data.final_cost.toFixed(2)}`;
            billingSection.style.display = "block";
            message.textContent = ""; // Clear previous messages
        } catch (error) {
            console.error("Error calculating cost:", error);
            message.textContent = `Error: ${error.message || "Unable to calculate cost."}`;
            message.style.color = "red";
            billingSection.style.display = "none";
        }
    });

    // Step 2: Confirm Reservation
    confirmPaymentBtn.addEventListener("click", async () => {
        try {
            // Validate required reservation details
            if (!reservationDetails.vehicle_id || !reservationDetails.start_time || !reservationDetails.end_time) {
                throw new Error("Reservation details are incomplete. Please calculate the cost again.");
            }

            // Send the reservation request
            const response = await fetch(`${reservationServiceUrl}/reservations`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${token}`,
                },
                body: JSON.stringify({
                    vehicle_id: reservationDetails.vehicle_id,
                    start_time: reservationDetails.start_time,
                    end_time: reservationDetails.end_time,
                }),
            });

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(errorText);
            }

            const responseData = await response.text();
            message.textContent = responseData; // Display server response
            message.style.color = "green";
            bookingForm.reset();
            billingSection.style.display = "none";
        } catch (error) {
            console.error("Error making reservation:", error);
            message.textContent = `Error: ${error.message || "Unable to create reservation."}`;
            message.style.color = "red";
        }
    });
});
