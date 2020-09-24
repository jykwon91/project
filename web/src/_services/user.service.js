import config from "config";
import { authHeader } from "../_helpers";

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
  getLandLordList,
  sendNotification,
  sendServiceReq,
  updateServiceReq,
  updateUser,
  getTenantList,
  getPaymentList,
  getPaymentOverview,
  updatePayment,
  createTestPayment,
  //delete: _delete
};

//production
//const apiUrl = "http://rentalmgmt.co:8080";

//dev
const apiUrl = "http://10.0.0.152:8080";

function login(email, password) {
  let header = new Headers({
    "Access-Control-Allow-Origin": "*",
    "Content-Type": "application/json",
  });

  const requestOptions = {
    method: "POST",
    mode: "cors",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  };

  //return fetch("http://rentalmgmt.co:8000/users/authenticate", requestOptions)
  return fetch(apiUrl + "/users/authenticate", requestOptions)
    .then(handleResponse)
    .then((user) => {
      // login successful if there's a jwt token in the response
      if (user.token) {
        // store user details and jwt token in local storage to keep user logged in between page refreshes
        localStorage.setItem("user", JSON.stringify(user));
      }

      return user;
    });
}

function logout() {
  // remove user from local storage to log user out
  localStorage.removeItem("user");
}

function getAll() {
  const requestOptions = {
    method: "GET",
    mode: "cors",
    headers: authHeader(),
  };

  //return fetch("http://rentalmgmt.co:8000/users/all", requestOptions).then(handleResponse);
  return fetch(apiUrl + "/users/all", requestOptions).then(handleResponse);
}

function getAllNotifications() {
  const requestOptions = {
    method: "GET",
    mode: "cors",
    headers: authHeader(),
  };

  //return fetch("http://rentalmgmt.co:8000/users/notification/all", requestOptions).then(handleResponse);
  return fetch(apiUrl + "/users/notification/all", requestOptions).then(
    handleResponse
  );
}

function getServiceRequestList() {
  const requestOptions = {
    method: "GET",
    mode: "cors",
    headers: authHeader(),
  };

  return fetch(apiUrl + "/users/service/all", requestOptions).then(
    handleResponse
  );
}

function getAllLandLordProperties() {
  const requestOptions = {
    method: "GET",
    mode: "cors",
    headers: authHeader(),
  };

  return fetch(apiUrl + "/landlord/property/all", requestOptions).then(
    handleResponse
  );
}

function getCurrentUser() {
  const requestOptions = {
    method: "GET",
    mode: "cors",
    headers: authHeader(),
  };

  return fetch(apiUrl + "/users/currentUser", requestOptions).then(
    handleResponse
  );
}

function getStateList() {
  const requestOptions = {
    method: "GET",
    mode: "cors",
    headers: { "Content-Type": "application/json" },
  };

  //return fetch("http://rentalmgmt.co:8000/stateList", requestOptions).then(handleResponse);
  return fetch(apiUrl + "/stateList", requestOptions).then(handleResponse);
}

function getLandLordList() {
  const requestOptions = {
    method: "GET",
    mode: "cors",
    headers: { "Content-Type": "application/json" },
  };

  //return fetch("http://rentalmgmt.co:8000/users/landlord/all", requestOptions).then(handleResponse);
  return fetch(apiUrl + "/users/landlord/all", requestOptions).then(
    handleResponse
  );
}

function getTenantList() {
  const requestOptions = {
    method: "GET",
    mode: "cors",
    headers: authHeader(),
  };

  //return fetch("http://rentalmgmt.co:8000/landlord/tenant/all", requestOptions).then(handleResponse);
  return fetch(apiUrl + "/landlord/tenant/all", requestOptions).then(
    handleResponse
  );
}

function register(user) {
  const requestOptions = {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(user),
  };

  //return fetch("http://rentalmgmt.co:8000/users/register", requestOptions).then(handleResponse);
  return fetch(apiUrl + "/users/register", requestOptions).then(handleResponse);
}

function registerLandLordProperty(address) {
  const requestOptions = {
    method: "POST",
    headers: authHeader(),
    body: JSON.stringify(address),
  };

  //return fetch("http://rentalmgmt.co:8000/users/landlord/property/register", requestOptions).then(handleResponse);
  return fetch(
    apiUrl + "/users/landlord/property/register",
    requestOptions
  ).then(handleResponse);
}

function sendNotification(notification) {
  const requestOptions = {
    method: "POST",
    headers: authHeader(),
    body: JSON.stringify(notification),
  };

  //return fetch("http://rentalmgmt.co:8000/landlord/notification", requestOptions).then(handleResponse);
  return fetch(apiUrl + "/landlord/notification", requestOptions).then(
    handleResponse
  );
}

function sendServiceReq(serviceReq) {
  const requestOptions = {
    method: "POST",
    headers: authHeader(),
    body: JSON.stringify(serviceReq),
  };

  //return fetch("http://rentalmgmt.co:8000/tenant/service/request", requestOptions).then(handleResponse);
  return fetch(apiUrl + "/tenant/service/request", requestOptions).then(
    handleResponse
  );
}

function updateServiceReq(serviceReq) {
  const requestOptions = {
    method: "POST",
    headers: authHeader(),
    body: JSON.stringify(serviceReq),
  };

  //return fetch("http://rentalmgmt.co:8000/landlord/service/request/update", requestOptions).then(handleResponse);
  return fetch(
    apiUrl + "/landlord/service/request/update",
    requestOptions
  ).then(handleResponse);
}

function updateUser(user) {
  const requestOptions = {
    method: "POST",
    headers: authHeader(),
    body: JSON.stringify(user),
  };

  //return fetch("http://rentalmgmt.co:8000/users/update", requestOptions).then(handleResponse);
  return fetch(apiUrl + "/users/update", requestOptions).then(handleResponse);
}

function getPaymentList() {
  const requestOptions = {
    method: "GET",
    mode: "cors",
    headers: authHeader(),
  };

  //return fetch("http://rentalmgmt.co:8000/users/payment/all", requestOptions).then(handleResponse);
  return fetch(apiUrl + "/users/payment/all", requestOptions).then(
    handleResponse
  );
}

function getPaymentOverview(landLordID) {
  const requestOptions = {
    method: "GET",
    mode: "cors",
    headers: authHeader(),
  };

  //return fetch("http://rentalmgmt.co:8000/tenant/payment/overview/" + landLordID, requestOptions).then(handleResponse);
  return fetch(
    apiUrl + "/tenant/payment/overview/" + landLordID,
    requestOptions
  ).then(handleResponse);
}

function createTestPayment() {
  const requestOptions = {
    method: "GET",
    mode: "cors",
    headers: authHeader(),
  };

  //return fetch("http://rentalmgmt.co:8000/tenant/payment/overview/" + landLordID, requestOptions).then(handleResponse);
  return fetch(apiUrl + "/test/create/payment", requestOptions).then(
    handleResponse
  );
}

function updatePayment(paymentInfo) {
  const requestOptions = {
    method: "POST",
    headers: authHeader(),
    body: JSON.stringify(paymentInfo),
  };

  //return fetch("http://rentalmgmt.co:8000/users/update", requestOptions).then(handleResponse);
  return fetch(apiUrl + "/tenant/payment/update", requestOptions).then(
    handleResponse
  );
}

function handleResponse(response) {
  return response.text().then((text) => {
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
