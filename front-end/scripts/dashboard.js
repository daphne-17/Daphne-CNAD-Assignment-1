// Base URLs for each microservice
const authServiceUrl = "http://localhost:8080";
const userServiceUrl = "http://localhost:8081";
const vehicleServiceUrl = "http://localhost:8082";

// Token from auth-service (assume it is already stored after login)
const token = localStorage.getItem("authToken");

// Redirect to login if the user is not authenticated
if (!token) {
    alert("Please log in to access the dashboard.");
    window.location.href = "login.html";
}

// DOM Elements
const dashboardContent = document.getElementById("dashboard-content");

// View Profile
document.getElementById("view-profile-btn").addEventListener("click", () => {
    window.location.href = "userprofile.html"; // Redirect to the user profile page
});


// Redirect to the "Available Vehicles" page
document.getElementById("view-cars-btn").addEventListener("click", () => {
    window.location.href = "vehicles.html";
});

// Redirect to the "Make a Reservation" page
document.getElementById("make-reservation-btn").addEventListener("click", () => {
    window.location.href = "createbooking.html";
});

// Redirect to the "View Reservations" page
document.getElementById("view-reservations-btn").addEventListener("click", () => {
    window.location.href = "viewbooking.html";
});

