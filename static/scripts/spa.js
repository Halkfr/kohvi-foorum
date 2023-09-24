document.addEventListener('click', e => {
    if (["home-nav", "create-post-nav", "profile-nav", "sign-out-nav", "register-link"].includes(e.target.id)) {
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
    '/view-profile': 'profile.html',
    '/sign-out': 'sign-in.html',
    '/sing-in': '',
    '/sign-up': 'sign-up.html'
}

const handleLocation = async () => {
    const path = window.location.pathname;
    const html = await fetch('./static/templates/' + routers[path]).then((data) => data.text()).catch(error => {
        console.error('Error:', error);
    });
    console.log(html)
    if (path === "/sign-out") {
        document.body.innerHTML = html;
    } else {
        if (path === "/sign-up") {
            const header = await fetch('./static/templates/header.html').then((data) => data.text());
            document.body.innerHTML = header + html;
        } else {
            const header = await fetch('./static/templates/header.html').then((data) => data.text());
            const sidepanel = await fetch('./static/templates/sidepanel.html').then((data) => data.text());
            document.body.innerHTML = header + sidepanel;
            document.querySelector('#content-container').innerHTML = html;
        }
    }
}

window.onpopstate = handleLocation;
window.route = route;
handleLocation();