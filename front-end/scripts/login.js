document.addEventListener("DOMContentLoaded", () => {
    document.getElementById("login-form").addEventListener("submit", async (e) => {
        e.preventDefault();

        // Get the identifier (email or phone) and password values
        const identifier = document.getElementById("identifier").value;
        const password = document.getElementById("password").value;

        // Ensure login-message exists
        const loginMessage = document.getElementById("login-message");
        if (!loginMessage) {
            console.error("Element with ID 'login-message' not found");
            return;
        }

        try {
            // Send login request to the backend
            const response = await fetch("http://localhost:8080/login", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ identifier, password }),
            });

            if (response.ok) {
                const data = await response.json();

                // Store token and user ID in localStorage
                localStorage.setItem("authToken", data.token);
                localStorage.setItem("userID", data.user_id);

                // Display success message and redirect
                loginMessage.textContent = "Login successful! Redirecting...";
                loginMessage.style.color = "green";

                setTimeout(() => {
                    window.location.href = "dashboard.html";
                }, 1000);
            } else {
                const error = await response.text();
                loginMessage.textContent = `Error: ${error}`;
                loginMessage.style.color = "red";
            }
        } catch (error) {
            loginMessage.textContent = "Error: Unable to log in. Please try again.";
            loginMessage.style.color = "red";
            console.error("Login error:", error);
        }
    });
});
