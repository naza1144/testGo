async function  login(event) {
    event.preventDefault();
    const username = document.querySelector('username').value.trim();
    const password = document.querySelector('password').value.trim();

    const res = await fetch('/api/users/login', {
        metod: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password }),
    });

    const data = await Response.json();

    if (data.status == "ok") {
        window.location.href = "/templates/page/succes";
    } else {
        alert("❌ อีเมลหรือรหัสผ่านไม่ถูกต้อง");
    }
}