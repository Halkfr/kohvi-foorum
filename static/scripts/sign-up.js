function hide1() {
    window.history.replaceState({}, '', 'http://127.0.0.1:8080/sign-in')
    window.location.href = '/sign-in'
}

function load1() {
    signUpForm = document.querySelector('form.sign-up')
    if (signUpForm['sign-up-fname'].checkValidity() && signUpForm['sign-up-lname'].checkValidity()) {
        document.getElementById("name-form").classList.add("d-none")
        document.getElementById("age-gender-form").classList.remove("d-none")
    } else {
        if (!signUpForm['sign-up-lname'].checkValidity()) {
            signUpForm['sign-up-lname'].reportValidity()
        }
        if (!signUpForm['sign-up-fname'].checkValidity()) {
            signUpForm['sign-up-fname'].reportValidity()
        }
    }
}

function hide2() {
    document.getElementById("age-gender-form").classList.add("d-none")
    document.getElementById("name-form").classList.remove("d-none")
}

function load2() {
    if (signUpForm['sign-up-dbirth'].checkValidity()) {
        document.getElementById("age-gender-form").classList.add("d-none")
        document.getElementById("sign-in-data-form").classList.remove("d-none")
    } else {
        signUpForm['sign-up-dbirth'].reportValidity()
    }
}

function hide3() {
    document.getElementById("sign-in-data-form").classList.add("d-none")
    document.getElementById("age-gender-form").classList.remove("d-none")
}

document.addEventListener('click', function (e) {
    if (location.pathname === "/sign-up") {
        signUpForm = document.querySelector('form.sign-up')
        if (signUpForm && signUpForm.checkValidity()) {
            if (e.target.matches('#sign-up-submit-button')) {
                let x = document.querySelector('form.sign-up').elements;
                e.preventDefault()

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
                            "password": x['sign-up-pass'].value,
                            "sessionStatus": "Offline"
                        },
                    )
                }).then(response => {
                    if (response.ok) { window.history.replaceState({}, '', 'http://127.0.0.1:8080/sign-in'), window.location.href = '/sign-in' } else {
                        if (response.status === 403) {
                            alert('Account with such Username or Email already exists');
                        } else {
                            console.error(data.error);
                        }
                    }
                })
            }
        }
    }
})