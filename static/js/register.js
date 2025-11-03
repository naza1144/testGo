console.log("register.js loaded");

async function register(event) {
  event.preventDefault();
  console.log("submit handler fired");

  const nameEl = document.getElementById('name');
  const emailEl = document.getElementById('register-email');
  const passwordEl = document.getElementById('register-password');

  if (!nameEl || !emailEl || !passwordEl) {
    console.error("Missing fields", { nameEl, emailEl, passwordEl });
    return;
  }

  const payload = { name: nameEl.value.trim(), email: emailEl.value.trim(), password: passwordEl.value.trim() };
  console.log("payload ->", payload);

  try {
    const res = await fetch('/api/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    });
    console.log("fetch status:", res.status);
    const bodyText = await res.text();
    console.log("response body:", bodyText);

    if (res.ok) window.location.href = "login";
    else {
      let data = {};
      try { data = JSON.parse(bodyText); } catch {}
      alert(data.error || "สมัครไม่สำเร็จ");
    }
  } catch (err) {
    console.error("fetch error:", err);
    alert("Network error");
  }
}

const form = document.querySelector('#register-form');
if (form) {
  form.addEventListener('submit', register);
  console.log("submit listener attached");
} else {
  console.warn("No #register-form found");
}