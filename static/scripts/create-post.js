function changeDropdownText(item) {
    let element = document.getElementById("titleCategory").classList
    element.remove("btn-light", "btn-primary", "btn-warning", "btn-danger", "btn-success")
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
            if (document.querySelector('form.create-post')['titleCategory'].innerHTML == "Select Category") {
                alert("Please select category")
            } else if (document.querySelector('form.create-post').checkValidity()) {
                e.preventDefault();
                let formData = new FormData();
                formData.append('thread', document.querySelector('form.create-post')['titleCategory'].innerHTML);
                formData.append('title', document.querySelector('form.create-post #create-post-title').value);
                formData.append('image', document.querySelector('form.create-post #create-post-img').files[0]);
                formData.append('image-name', document.querySelector('form.create-post #create-post-img').value.split(/(\\|\/)/g).pop());
                formData.append('content', document.querySelector('form.create-post #create-post-content').value);

                fetch('http://127.0.0.1:8080/api/add-post', {
                    method: 'POST',
                    headers: {
                        'credentials': 'include',
                    },
                    body: formData
                }).then(response => {
                    if (response.ok) { window.history.pushState({}, '', 'home'); handleLocation(); } else { /*display error*/ }
                })
            } else {
                if (!document.querySelector('form.create-post #create-post-content').checkValidity()) {
                    document.querySelector('form.create-post #create-post-content').reportValidity()
                }
                if (!document.querySelector('form.create-post #create-post-title').checkValidity()) {
                    document.querySelector('form.create-post #create-post-title').reportValidity()
                }
            }
        }
    }
})
