async function login(event) {
    event.preventDefault();
    const email = document.getElementById('email').value.trim();
    const password = document.getElementById('password').value.trim();

    const res = await fetch('/api/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
    });

    if (res.ok) {
        // redirect to home or success page
        window.location.href = "succes";
    } else {
        const data = await res.json().catch(() => ({}));
        alert(data.error || "❌ อีเมลหรือรหัสผ่านไม่ถูกต้อง");
    }
}

// attach listener (ensure form exists)
const form = document.querySelector('form');
if (form) form.addEventListener('submit', login);