document.addEventListener('click', e => {
    if (["home-nav", "create-post-nav"].includes(e.target.id)) {
        route(e)
    }
    e.preventDefault();
});

const route = (e) => {
    window.history.pushState({}, '', e.target.href);
}

window.route = route;