import { userConstants } from '../_constants';

export function landLordPropertyRegistration(state = {}, action) {
  switch (action.type) {
    case userConstants.REGISTER_LANDLORD_PROPERTY_REQUEST:
      return { registering: true };
    case userConstants.REGISTER_LANDLORD_PROPERTY_SUCCESS:
      return {};
    case userConstants.REGISTER_LANDLORD_PROPERTY_FAILURE:
      return {};
    default:
      return state
  }
}
