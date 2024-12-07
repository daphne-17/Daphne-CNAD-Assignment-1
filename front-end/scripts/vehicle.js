// Base URL for the vehicle service
const vehicleServiceUrl = "http://localhost:8082"; // Replace with actual vehicle service URL
const token = localStorage.getItem("authToken"); // Fetch the token from localStorage

// DOM Element to display vehicles
const vehicleList = document.getElementById("vehicle-list");

// Fetch available vehicles from the backend
async function fetchAvailableVehicles() {
    try {
        const response = await fetch(`${vehicleServiceUrl}/vehicles/available`, {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${token}`, // Include the JWT token in the header
            },
        });

        if (response.ok) {
            const cars = await response.json();
            if (cars.length === 0) {
                vehicleList.innerHTML = "<p>No vehicles available at the moment.</p>";
            } else {
                vehicleList.innerHTML = `
                    ${cars.map(car => `
                        <div class="vehicle-card">
                            <p><strong>ID for booking:</strong> ${car.vehicle_id}</p>
                            <p><strong>Model:</strong> ${car.model}</p>
                            <p><strong>Rate per Hour:</strong> $${car.rate_per_hour}</p>
                            <p><strong>Status:</strong> ${car.status}</p>
                        </div>
                    `).join("<hr>")}`;
            }
        } else {
            const error = await response.text();
            vehicleList.innerHTML = `<p>Error: Unable to fetch vehicles. ${error}</p>`;
        }
    } catch (error) {
        console.error("Error fetching available vehicles:", error);
        vehicleList.innerHTML = "<p>Error: Unable to fetch vehicles. Please try again later.</p>";
    }
}

// Trigger the fetch function when the page loads
document.addEventListener("DOMContentLoaded", fetchAvailableVehicles);
