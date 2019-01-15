import { userConstants } from '../_constants';
import { userService } from '../_services';
import { alertActions } from './';
import { history } from '../_helpers';

export const userActions = {
    login,
    logout,
    register,
    getAll,
    getAllNotifications,
		getAllLandLordProperties,
		getServiceRequestList,
		getStateList,
		getCurrentUser,
		registerLandLordProperty,
		sendNotification,
		sendServiceReq,
		updateServiceReq,
    delete: _delete
};

function login(email, password) {
    return dispatch => {
        dispatch(request({ email }));

        userService.login(email, password)
            .then(
                user => { 
                    dispatch(success(user));
										if (user.userType == "tenant") {
                    	history.push('/');
										} else if(user.userType == "landlord") {
											history.push('/landlord');
										}
                },
                error => {
                    dispatch(failure(error.toString()));
                    dispatch(alertActions.error(error.toString()));
                }
            );
    };

    function request(user) { return { type: userConstants.LOGIN_REQUEST, user } }
    function success(user) { return { type: userConstants.LOGIN_SUCCESS, user } }
    function failure(error) { return { type: userConstants.LOGIN_FAILURE, error } }
}

function logout() {
    userService.logout();
    return { type: userConstants.LOGOUT };
}

function register(user) {
    return dispatch => {
        dispatch(request(user));

        userService.register(user)
            .then(
                user => { 
                    dispatch(success());
                    history.push('/login');
                    dispatch(alertActions.success('Registration successful'));
                },
                error => {
                    dispatch(failure(error.toString()));
                    dispatch(alertActions.error(error.toString()));
                }
            );
    };

    function request(user) { return { type: userConstants.REGISTER_REQUEST, user } }
    function success(user) { return { type: userConstants.REGISTER_SUCCESS, user } }
    function failure(error) { return { type: userConstants.REGISTER_FAILURE, error } }
}

function registerLandLordProperty(address) {
    return dispatch => {
        dispatch(request(address));

        userService.registerLandLordProperty(address)
            .then(
                address => { 
                    dispatch(success());
                    history.push('/landlord');
                    dispatch(alertActions.success('Registration successful'));
                },
                error => {
                    dispatch(failure(error.toString()));
                    dispatch(alertActions.error(error.toString()));
                }
            );
    };

    function request(address) { return { type: userConstants.REGISTER_LANDLORD_PROPERTY_REQUEST, address } }
    function success(address) { return { type: userConstants.REGISTER_LANDLORD_PROPERTY_SUCCESS, address } }
    function failure(error) { return { type: userConstants.REGISTER_LANDLORD_PROPERTY_FAILURE, error } }
}

function getAll() {
    return dispatch => {
        dispatch(request());

        userService.getAll()
            .then(
                users => dispatch(success(users)),
                error => dispatch(failure(error.toString()))
            );
    };

    function request() { return { type: userConstants.GETALL_REQUEST } }
    function success(users) { return { type: userConstants.GETALL_SUCCESS, users } }
    function failure(error) { return { type: userConstants.GETALL_FAILURE, error } }
}

function getAllNotifications() {
    return dispatch => {
        dispatch(request());

        userService.getAllNotifications()
            .then(
                notifications => dispatch(success(notifications)),
                error => dispatch(failure(error.toString()))
            );
    };

    function request() { return { type: userConstants.GETALL_NOTIFICATIONS_REQUEST } }
    function success(notifications) { return { type: userConstants.GETALL_NOTIFICATIONS_SUCCESS, notifications } }
    function failure(error) { return { type: userConstants.GETALL_NOTIFICATIONS_FAILURE, error } }
}

function getAllLandLordProperties() {
    return dispatch => {
        dispatch(request());

        userService.getAllLandLordProperties()
            .then(
                landLordPropertyList => dispatch(success(landLordPropertyList)),
                error => dispatch(failure(error.toString()))
            );
    };

    function request() { return { type: userConstants.GET_ALL_LANDLORD_PROPERTIES_REQUEST } }
    function success(landLordPropertyList) { return { type: userConstants.GET_ALL_LANDLORD_PROPERTIES_SUCCESS, landLordPropertyList } }
    function failure(error) { return { type: userConstants.GET_ALL_LANDLORD_PROPERTIES_FAILURE, error } }
}

function getServiceRequestList() {
    return dispatch => {
        userService.getServiceRequestList()
            .then(
                serviceRequestList => dispatch(success(serviceRequestList)),
                error => dispatch(failure(error.toString()))
            );
    };

    function request() { return { type: userConstants.GET_SERVICEREQ_LIST_REQUEST } }
    function success(serviceRequestList) { return { type: userConstants.GET_SERVICEREQ_LIST_SUCCESS, serviceRequestList } }
    function failure(error) { return { type: userConstants.GET_SERVICEREQ_LIST_FAILURE, error } }
}

function getCurrentUser() {
    return dispatch => {
        userService.getCurrentUser()
            .then(
                currentUser => dispatch(success(currentUser)),
                error => dispatch(failure(error.toString()))
            );
    };

    function request() { return { type: userConstants.GET_CURRENT_USER_REQUEST } }
    function success(currentUser) { return { type: userConstants.GET_CURRENT_USER_SUCCESS, currentUser } }
    function failure(error) { return { type: userConstants.GET_CURRENT_USER_FAILURE, error } }
}

function getStateList() {
    return dispatch => {
        dispatch(request());

        userService.getStateList()
            .then(
                stateList => dispatch(success(stateList)),
                error => dispatch(failure(error.toString()))
            );
    };

    function request() { return { type: userConstants.GET_STATE_LIST_REQUEST } }
    function success(stateList) { return { type: userConstants.GET_STATE_LIST_SUCCESS, stateList } }
    function failure(error) { return { type: userConstants.GET_STATE_LIST_FAILURE, error } }
}

function sendNotification(notification) {
    return dispatch => {
        dispatch(request(notification));

        userService.sendNotification(notification)
            .then(
                notification => { 
                    dispatch(success());
                    history.push('/landlord');
                    dispatch(alertActions.success('Sent notification(s)'));
                },
                error => {
                    dispatch(failure(error.toString()));
                    dispatch(alertActions.error(error.toString()));
                }
            );
    };

    function request(notification) { return { type: userConstants.SEND_NOTIFICATION_REQUEST, notification } }
    function success(notification) { return { type: userConstants.SEND_NOTIFICATION_SUCCESS, notification } }
    function failure(error) { return { type: userConstants.SEND_NOTIFICATION_FAILURE, error } }
}

function sendServiceReq(serviceReq) {
    return dispatch => {
        dispatch(request(serviceReq));

        userService.sendServiceReq(serviceReq)
            .then(
                serviceReq => { 
                    dispatch(success());
                    history.push('/');
                    dispatch(alertActions.success('Sent service request'));
                },
                error => {
                    dispatch(failure(error.toString()));
                    dispatch(alertActions.error(error.toString()));
                }
            );
    };

    function request(serviceReq) { return { type: userConstants.SEND_SERVICEREQ_REQUEST, serviceReq } }
    function success(serviceReq) { return { type: userConstants.SEND_SERVICEREQ_SUCCESS, serviceReq } }
    function failure(error) { return { type: userConstants.SEND_SERVICEREQ_FAILURE, error } }
}

function updateServiceReq(serviceReq) {
    return dispatch => {
        dispatch(request(serviceReq));

        userService.updateServiceReq(serviceReq)
            .then(
                serviceReq => { 
                    dispatch(success());
                },
                error => {
                    dispatch(failure(error.toString()));
                    dispatch(alertActions.error(error.toString()));
                }
            );
    };

    function request(serviceReq) { return { type: userConstants.UPDATE_SERVICEREQ_REQUEST, serviceReq } }
    function success(serviceReq) { return { type: userConstants.UPDATE_SERVICEREQ_SUCCESS, serviceReq } }
    function failure(error) { return { type: userConstants.UPDATE_SERVICEREQ_FAILURE, error } }
}

// prefixed function name with underscore because delete is a reserved word in javascript
function _delete(id) {
    return dispatch => {
        dispatch(request(id));

        userService.delete(id)
            .then(
                user => dispatch(success(id)),
                error => dispatch(failure(id, error.toString()))
            );
    };

    function request(id) { return { type: userConstants.DELETE_REQUEST, id } }
    function success(id) { return { type: userConstants.DELETE_SUCCESS, id } }
    function failure(id, error) { return { type: userConstants.DELETE_FAILURE, id, error } }
}
