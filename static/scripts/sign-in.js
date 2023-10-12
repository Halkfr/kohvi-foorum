document.addEventListener('click', function (e) {
    if (location.pathname === "/sign-in") {
        document.getElementById("remember-me").onclick = changeValue
        function changeValue() {
            let value = document.getElementById("remember-me").value
            if (value === "off") {
                document.getElementById("remember-me").value = "on"
            } else {
                document.getElementById("remember-me").value = "off"
            }
        }

        if (e.target.matches('#sign-in-submit-button')) {
            e.preventDefault();
            let x = document.querySelector('form.sign-in').elements;

            fetch('http://127.0.0.1:8080/api/sign-in', {
                method: 'POST',
                headers: {
                    'Accept': 'application/json',
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(
                    {
                        "username": x['sign-in-username-or-email'].value,
                        "password": x['sign-in-pass'].value
                    },
                )
            }).then(response => {
                if (response.ok) { window.history.replaceState({}, '', 'http://127.0.0.1:8080/home'), window.location.href = '/home' }
                else { /*display error*/ }
            })
        }
    }
})