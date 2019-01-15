import { userConstants } from '../_constants';

export function landLordPropertyList(state = {}, action) {
        switch (action.type) {
                case userConstants.GET_ALL_LANDLORD_PROPERTIES_REQUEST:
                        return {
                                loading: true
                        };
                case userConstants.GET_ALL_LANDLORD_PROPERTIES_SUCCESS:
                        return {
                                items: action.landLordPropertyList
                        };
                case userConstants.GET_ALL_LANDLORD_PROPERTIES_FAILURE:
                        return {
                                error: action.error
                        };
                default:
                        return state
        }
}

