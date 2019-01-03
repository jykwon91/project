import config from 'config';
import { authHeader } from '../_helpers';

export const userService = {
    login,
    logout,
    register,
    getAll,
		getAllNotifications,
    getById,
    update,
    delete: _delete
};

function login(email, password) {

	let header = new Headers({
		'Access-Control-Allow-Origin':'*',
		'Content-Type':'application/json'
	});

    const requestOptions = {
        method: 'POST',
				mode: 'cors',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password })
    };

    //return fetch(`${config.apiurl}/users/authenticate`, requestOptions)
    //return fetch("http://192.168.1.125:8000/users/authenticate", requestOptions)
    //return fetch("http://10.0.0.152:8000/users/authenticate", requestOptions)
    return fetch("http://rentalmgmt.co:8000/users/authenticate", requestOptions)
        .then(handleResponse)
        .then(user => {
            // login successful if there's a jwt token in the response
            if (user.token) {
                // store user details and jwt token in local storage to keep user logged in between page refreshes
                localStorage.setItem('user', JSON.stringify(user));
            }

            return user;
        });
}

function logout() {
    // remove user from local storage to log user out
    localStorage.removeItem('user');
}

function getAll() {
    const requestOptions = {
        method: 'GET',
				mode: 'cors',
        headers: authHeader()
    };

    //return fetch(`${config.apiUrl}/users`, requestOptions).then(handleResponse);
    //return fetch("http://10.0.0.152:8000/users/all", requestOptions).then(handleResponse);
    return fetch("http://rentalmgmt.co:8000/users/all", requestOptions).then(handleResponse);
}

function getAllNotifications() {
	console.log("_action getallnotifications");
	const requestOptions = {
		method: 'GET',
		mode: 'cors',
		headers: authHeader()
	};

	return fetch("http://rentalmgmt.co:8000/users/notification/all", requestOptions).then(handleResponse);
}

function getById(id) {
    const requestOptions = {
        method: 'GET',
        headers: authHeader()
    };

    return fetch(`${config.apiUrl}/users/${id}`, requestOptions).then(handleResponse);
}

function register(user) {
    const requestOptions = {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(user)
    };

    //return fetch(`${config.apiUrl}/users/register`, requestOptions).then(handleResponse);
    //return fetch("http://192.168.1.125:8000/users/register", requestOptions).then(handleResponse);
    //return fetch("http://10.0.0.152:8000/users/register", requestOptions).then(handleResponse);
    return fetch("http://rentalmgmt.co:8000/users/register", requestOptions).then(handleResponse);
}

function update(user) {
    const requestOptions = {
        method: 'PUT',
        headers: { ...authHeader(), 'Content-Type': 'application/json' },
        body: JSON.stringify(user)
    };

    return fetch(`${config.apiUrl}/users/${user.id}`, requestOptions).then(handleResponse);
}

// prefixed function name with underscore because delete is a reserved word in javascript
function _delete(id) {
    const requestOptions = {
        method: 'DELETE',
        headers: authHeader()
    };

    return fetch(`${config.apiUrl}/users/${id}`, requestOptions).then(handleResponse);
}

function handleResponse(response) {
		console.log(response);
    return response.text().then(text => {
        const data = text && JSON.parse(text);
        if (!response.ok) {
            if (response.status === 401) {
                // auto logout if 401 response returned from api
                logout();
                location.reload(true);
            }

            const error = (data && data.message) || response.statusText;
            return Promise.reject(error);
        }

        return data;
    });
}
