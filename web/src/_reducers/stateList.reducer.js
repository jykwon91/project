import { userConstants } from '../_constants';

export function stateList(state = {}, action) {
	switch (action.type) {
		case userConstants.GET_STATE_LIST_REQUEST:
			return {
				loading: true
			};
		case userConstants.GET_STATE_LIST_SUCCESS:
			return {
				items: action.stateList
			};
		case userConstants.GET_STATE_LIST_FAILURE:
			return {
				error: action.error
			};
		default:
			return state
	}
}
