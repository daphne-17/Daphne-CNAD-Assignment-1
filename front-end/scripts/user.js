const userServiceUrl = "http://localhost:8081"; // Base URL for User-Service
const token = localStorage.getItem("authToken"); // Retrieve token from localStorage
const userID = localStorage.getItem("userID"); // Retrieve user ID from localStorage

// Ensure the user is logged in
if (!token || !userID) {
    alert("Please log in to access your profile.");
    window.location.href = "login.html";
}

// DOM Elements
const profileDetails = document.getElementById("profile-details");
const membershipDetails = document.getElementById("membership-details");
const rentalHistory = document.getElementById("rental-history");
const updateProfileForm = document.getElementById("update-profile-form");
const updateMessage = document.getElementById("update-message");

// Fetch and display profile details
async function fetchUserProfile() {
    try {
        const response = await fetch(`${userServiceUrl}/users/profile`, {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${token}`,
            },
        });

        if (response.ok) {
            const data = await response.json();
            profileDetails.innerHTML = `
                <p><strong>Name:</strong> ${data.name}</p>
                <p><strong>Email:</strong> ${data.email}</p>
                <p><strong>Phone:</strong> ${data.phone}</p>
            `;

            // Prefill update form fields
            document.getElementById("name").value = data.name;
            document.getElementById("email").value = data.email;
            document.getElementById("phone").value = data.phone;
        } else {
            profileDetails.innerHTML = `<p>Error: Unable to fetch profile details.</p>`;
        }
    } catch (error) {
        profileDetails.innerHTML = `<p>Error: Failed to load profile details. Please try again later.</p>`;
        console.error(error);
    }
}

// Fetch and display membership status
async function fetchMembershipStatus() {
    try {
        const response = await fetch(`${userServiceUrl}/users/membership`, {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${token}`,
            },
        });

        if (response.ok) {
            const data = await response.json();
            membershipDetails.innerHTML = `
                <p><strong>Membership Tier:</strong> ${data.tier_name}</p>
                <p><strong>Hourly Rate:</strong> $${data.hourly_rate}</p>
                <p><strong>Booking Limit:</strong> ${data.booking_limit}</p>
            `;
        } else {
            membershipDetails.innerHTML = `<p>Error: Unable to fetch membership status.</p>`;
        }
    } catch (error) {
        membershipDetails.innerHTML = `<p>Error: Failed to load membership status. Please try again later.</p>`;
        console.error(error);
    }
}

// Fetch and display rental history
async function fetchRentalHistory() {
    try {
        const response = await fetch(`${userServiceUrl}/users/rentalhistory`, {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${token}`,
            },
        });

        if (response.ok) {
            const data = await response.json();
            const rentals = data.rental_history;

            if (!rentals || rentals.length === 0) {
                rentalHistory.innerHTML = "<p>No rental history available.</p>";
            } else {
                rentalHistory.innerHTML = rentals
                    .map(
                        (rental) => `
                        <p>
                            <strong>Reservation ID:</strong> ${rental.reservation_id}<br>
                            <strong>Vehicle ID:</strong> ${rental.vehicle_id}<br>
                            <strong>Start Time:</strong> ${rental.start_time}<br>
                            <strong>End Time:</strong> ${rental.end_time}<br>
                            <strong>Status:</strong> ${rental.status}
                        </p>
                    `
                    )
                    .join("<hr>");
            }
        } else {
            rentalHistory.innerHTML = `<p>Error: Unable to fetch rental history.</p>`;
        }
    } catch (error) {
        rentalHistory.innerHTML = `<p>Error: Failed to load rental history. Please try again later.</p>`;
        console.error(error);
    }
}

// Handle profile update
updateProfileForm.addEventListener("submit", async (e) => {
    e.preventDefault();

    const name = document.getElementById("name").value;
    const email = document.getElementById("email").value;
    const phone = document.getElementById("phone").value;

    try {
        const response = await fetch(`${userServiceUrl}/users/update`, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${token}`,
            },
            body: JSON.stringify({ name, email, phone }),
        });

        if (response.ok) {
            updateMessage.textContent = "Profile updated successfully!";
            updateMessage.style.color = "green";

            // Refresh profile details
            fetchUserProfile();
        } else {
            const error = await response.text();
            updateMessage.textContent = `Error: ${error}`;
            updateMessage.style.color = "red";
        }
    } catch (error) {
        updateMessage.textContent = "Error: Unable to update profile. Please try again.";
        updateMessage.style.color = "red";
        console.error(error);
    }
});

// Fetch all data on page load
fetchUserProfile();
fetchMembershipStatus();
fetchRentalHistory();
