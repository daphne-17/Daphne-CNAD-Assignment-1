document.getElementById('signup-form').addEventListener('submit', async (e) => {
    e.preventDefault();

    const username = document.getElementById('username').value;
    const email = document.getElementById('email').value;
    const phone = document.getElementById('phone').value;
    const password = document.getElementById('password').value;

    if (!/^\d{8}$/.test(phone)) {
        alert("Phone number must be exactly 8 digits.");
        return;
    }

    try {
        const response = await fetch('http://localhost:8080/signup', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, email, phone, password }),
        });

        if (response.ok) {
            alert('Signup successful!');
            window.location.href = 'login.html';
        } else {
            const error = await response.text();
            alert('Signup failed: ' + error);
        }
    } catch (err) {
        alert('Error: ' + err.message);
    }
});
