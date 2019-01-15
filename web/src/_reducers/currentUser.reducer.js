import { userConstants } from '../_constants';

export function currentUser(state = {}, action) {
	switch (action.type) {
		case userConstants.GET_CURRENT_USER_REQUEST:
			return {
				loading: true
			};
		case userConstants.GET_CURRENT_USER_SUCCESS:
			return {
				items: action.currentUser
			};
		case userConstants.GET_CURRENT_USER_FAILURE:
			return {
				error: action.error
			};
		default:
			return state
	}
}
