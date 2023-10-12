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