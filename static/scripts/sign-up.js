function hide1() {
    window.history.replaceState({}, '', 'http://127.0.0.1:8080/sign-in')
    window.location.href = '/sign-in'
}

function load1() {
    // allow if first and last names are filled
    document.getElementById("name-form").classList.add("d-none")
    document.getElementById("age-gender-form").classList.remove("d-none")
}

function hide2() {
    document.getElementById("age-gender-form").classList.add("d-none")
    document.getElementById("name-form").classList.remove("d-none")
}

function load2() {
    // allow if date of birth is selected
    document.getElementById("age-gender-form").classList.add("d-none")
    document.getElementById("sign-in-data-form").classList.remove("d-none")
}

function hide3() {
    document.getElementById("sign-in-data-form").classList.add("d-none")
    document.getElementById("age-gender-form").classList.remove("d-none")
}

document.addEventListener('click', function (e) {
    if (location.pathname === "/sign-up") {
        if (e.target.matches('#sign-up-submit-button')) {
            e.preventDefault();
            let x = document.querySelector('form.sign-up').elements;

            fetch('http://127.0.0.1:8080/api/sign-up', {
                method: 'POST',
                headers: {
                    'Accept': 'application/json',
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(
                    {
                        "firstName": x['sign-up-fname'].value,
                        "lastName": x['sign-up-lname'].value,
                        "birthdate": x['sign-up-dbirth'].value,
                        "gender": x['sign-up-gender'].value,
                        "username": x['sign-up-username'].value,
                        "email": x['sign-up-email'].value,
                        "password": x['sign-up-pass'].value
                    },
                )
            }).then(response => {
                if (response.ok) { window.history.replaceState({}, '', 'http://127.0.0.1:8080/sign-in'), window.location.href = '/sign-in' } else { /*write error*/ }
            })
        }
    }
})