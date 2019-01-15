import config from 'config';
import { authHeader } from '../_helpers';

export const userService = {
    login,
    logout,
    register,
		registerLandLordProperty,
    getAll,
		getAllLandLordProperties,
		getAllNotifications,
		getCurrentUser,
		getStateList,
		getServiceRequestList,
		sendNotification,
		sendServiceReq,
		updateServiceReq,
    //delete: _delete
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
	const requestOptions = {
		method: 'GET',
		mode: 'cors',
		headers: authHeader()
	};

	//return fetch("http://10.0.0.152:8000/users/notification/all", requestOptions).then(handleResponse);
	return fetch("http://rentalmgmt.co:8000/users/notification/all", requestOptions).then(handleResponse);
}

function getServiceRequestList() {
	const requestOptions = {
		method: 'GET',
		mode: 'cors',
		headers: authHeader()
	};

	//return fetch("http://10.0.0.152:8000/users/notification/all", requestOptions).then(handleResponse);
	return fetch("http://rentalmgmt.co:8000/users/service/all", requestOptions).then(handleResponse);
}

function getAllLandLordProperties() {
	const requestOptions = {
		method: 'GET',
		mode: 'cors',
		headers: authHeader()
	};

	//return fetch("http://10.0.0.152:8000/users/notification/all", requestOptions).then(handleResponse);
	return fetch("http://rentalmgmt.co:8000/landlord/property/all", requestOptions).then(handleResponse);
}

function getCurrentUser() {
	const requestOptions = {
		method: 'GET',
		mode: 'cors',
		headers: authHeader()
	};

	//return fetch("http://10.0.0.152:8000/users/notification/all", requestOptions).then(handleResponse);
	return fetch("http://rentalmgmt.co:8000/users/currentUser", requestOptions).then(handleResponse);
}

function getStateList() {
	const requestOptions = {
		method: 'GET',
		mode: 'cors',
		headers: {'Content-Type': 'application/json'},
	};

	//return fetch("http://10.0.0.152:8000/stateList", requestOptions).then(handleResponse);
	return fetch("http://rentalmgmt.co:8000/stateList", requestOptions).then(handleResponse);
}

function register(user) {
    const requestOptions = {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(user)
    };

    //return fetch(`${config.apiUrl}/users/register`, requestOptions).then(handleResponse);
    //return fetch("http://10.0.0.152:8000/users/register", requestOptions).then(handleResponse);
    return fetch("http://rentalmgmt.co:8000/users/register", requestOptions).then(handleResponse);
}

function registerLandLordProperty(address) {
    const requestOptions = {
        method: 'POST',
        //headers: { 'Content-Type': 'application/json' },
        headers: authHeader(),
        body: JSON.stringify(address)
    };

    //return fetch(`${config.apiUrl}/users/register`, requestOptions).then(handleResponse);
    //return fetch("http://10.0.0.152:8000/users/register", requestOptions).then(handleResponse);
    return fetch("http://rentalmgmt.co:8000/users/landlord/property/register", requestOptions).then(handleResponse);
}

function sendNotification(notification) {
    const requestOptions = {
        method: 'POST',
        //headers: { 'Content-Type': 'application/json' },
        headers: authHeader(),
        body: JSON.stringify(notification)
    };

    //return fetch(`${config.apiUrl}/users/register`, requestOptions).then(handleResponse);
    //return fetch("http://10.0.0.152:8000/users/register", requestOptions).then(handleResponse);
    return fetch("http://rentalmgmt.co:8000/landlord/notification", requestOptions).then(handleResponse);
}

function sendServiceReq(serviceReq) {
    const requestOptions = {
        method: 'POST',
        //headers: { 'Content-Type': 'application/json' },
        headers: authHeader(),
        body: JSON.stringify(serviceReq)
    };

    //return fetch(`${config.apiUrl}/users/register`, requestOptions).then(handleResponse);
    //return fetch("http://10.0.0.152:8000/users/register", requestOptions).then(handleResponse);
    return fetch("http://rentalmgmt.co:8000/tenant/service/request", requestOptions).then(handleResponse);
}

function updateServiceReq(serviceReq) {
    const requestOptions = {
        method: 'POST',
        //headers: { 'Content-Type': 'application/json' },
        headers: authHeader(),
        body: JSON.stringify(serviceReq)
    };

    //return fetch(`${config.apiUrl}/users/register`, requestOptions).then(handleResponse);
    //return fetch("http://10.0.0.152:8000/users/register", requestOptions).then(handleResponse);
    return fetch("http://rentalmgmt.co:8000/landlord/service/request/update", requestOptions).then(handleResponse);
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
