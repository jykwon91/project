import { userConstants } from '../_constants';

export function landLordList(state = {}, action) {
	switch (action.type) {
		case userConstants.GET_LANDLORD_LIST_REQUEST:
			return {
				loading: true
			};
		case userConstants.GET_LANDLORD_LIST_SUCCESS:
			return {
				items: action.landLordList
			};
		case userConstants.GET_LANDLORD_LIST_FAILURE:
			return {
				error: action.error
			};
		default:
			return state
	}
}
