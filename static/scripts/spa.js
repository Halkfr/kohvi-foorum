document.addEventListener('click', e => {
    if (["home-nav", "create-post-nav"].includes(e.target.id)) {
        route(e)
    }
    e.preventDefault();
});

const route = (e) => {
    window.history.pushState({}, '', e.target.href);
    handleLocation();
}

const routers = {
    '/home': 'home.html',
    '/create-post': 'create-post.html',
    '/profile': 'profile.html',
    '/sing-in': '',
    '/sign-up': ''
}

const handleLocation = async () => {
    const path = window.location.pathname;
    const html = await fetch('./static/templates/' + routers[path]).then((data) => data.text()).catch(error => {
        console.error('Error:', error);
    });
    // if (path === "/sign-in" || path === 'sign-up') {

    // } else {
        document.querySelector('#content-container').innerHTML = html;
    // }
}

window.route = route;