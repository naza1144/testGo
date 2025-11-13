async function login(event) {
    event.preventDefault();
    const email = document.getElementById('email').value.trim();
    const password = document.getElementById('password').value.trim();

    if (email === "" || password === "") {
        alert('กรุณากรอกอีเมลและรหัสผ่าน');
        return;
    }

    const res = await fetch('/api/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
    });

    // Expect server to return JSON { ok: true, role: "admin" } on success
    if (res.ok) {
        const data = await res.json().catch(() => ({}));
        if (data && data.role === 'admin') {
            // go to admin page
            window.location.href = '/admin';
        } else {
            // normal user success page
            window.location.href = '/succes';
        }
    } else {
        const data = await res.json().catch(() => ({}));
        alert(data.error || "❌ อีเมลหรือรหัสผ่านไม่ถูกต้อง");
    }
}

// attach listener (ensure form exists)
const form = document.querySelector('form');
if (form) form.addEventListener('submit', login);