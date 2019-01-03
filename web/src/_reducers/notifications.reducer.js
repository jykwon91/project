import { userConstants } from '../_constants';

export function notifications(state = {}, action) {
	switch (action.type) {
		case userConstants.GETALL_NOTIFICATIONS_REQUEST:
			return {
				loading: true
			};
		case userConstants.GETALL_NOTIFICATIONS_SUCCESS:
			return {
				items: action.notifications
			};
		case userConstants.GETALL_NOTIFICATIONS_FAILURE:
			return {
				error: action.error
			};
		default:
			return state
	}
}
