import { userConstants } from "../_constants";

export function updatePayment(state = {}, action) {
  switch (action.type) {
    case userConstants.UPDATE_PAYMENT_INFO_REQUEST:
      return { updating: true };
    case userConstants.UPDATE_PAYMENT_INFO_SUCCESS:
      return {};
    case userConstants.UPDATE_PAYMENT_INFO_FAILURE:
      return {};
    default:
      return state;
  }
}
