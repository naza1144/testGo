console.log("succes.js loaded");

document.addEventListener('DOMContentLoaded', function() {
    const logoutBtn  = document.getElementById('logout-rly');

    if (logoutBtn) {
        logoutBtn.addEventListener('click', function(event) {
            event.preventDefault();

            if (confirm("คุณแน่ใจหรือว่าต้องการออกจากระบบ?")) {
                window.location.href = "login";
            }
        });
    }
});