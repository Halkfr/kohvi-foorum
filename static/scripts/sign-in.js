document.addEventListener('click', function (e) {
    if (location.pathname === "/sign-in") {
        e.preventDefault()
        if (e.target.matches('#sign-in-submit-button')) {
            let x = document.querySelector('form.sign-in').elements;
            if (x['sign-in-username-or-email'].value !== "" || x['sign-in-pass'].value !== "") {
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
                    else {
                        alert("invalid username or password")
                    }
                })
            } else {
                alert("field should be filled")
            }
        }
    }
})