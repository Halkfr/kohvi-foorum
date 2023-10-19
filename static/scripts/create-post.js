function changeDropdownText(item) {
    let element = document.getElementById("titleCategory").classList
    element.remove("btn-default", "btn-primary", "btn-warning", "btn-danger", "btn-success")
    switch (item.innerHTML) {
        case "Question":
            document.getElementById("titleCategory").classList.add("btn-primary")
            break;
        case "Buy/Sell":
            document.getElementById("titleCategory").classList.add("btn-warning")
            break;
        case "Help!":
            document.getElementById("titleCategory").classList.add("btn-danger")
            break;
        case "Discussion":
            document.getElementById("titleCategory").classList.add("btn-success")
            break;
    }
    document.getElementById("titleCategory").innerHTML = item.innerHTML;
}

document.addEventListener('click', function (e) {
    if (location.pathname === "/create-post") {
        if (e.target.matches('#post-submit-btn')) {
            e.preventDefault();
            let x = document.querySelector('form.create-post').elements;

            fetch('http://127.0.0.1:8080/api/add-post', {
                method: 'POST',
                headers: {
                    'Accept': 'application/json',
                    'Content-Type': 'application/json',
                    'credentials': 'include',
                },
                body: JSON.stringify(
                    {
                        thread: x['titleCategory'].innerHTML,
                        title: x['create-post-title'].value,
                        image: x['create-post-img'].value,
                        content: x['create-post-content'].value
                    },
                )
            }).then(response => {
                if (response.ok) { window.location.href = '/home' } else { /*write error*/ }
            })
        }
    }
})