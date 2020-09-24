import { userConstants } from "../_constants";

export function paymentOverview(state = {}, action) {
  switch (action.type) {
    case userConstants.GET_PAYMENT_OVERVIEW_REQUEST:
      return {
        loading: true,
      };
    case userConstants.GET_PAYMENT_OVERVIEW_SUCCESS:
      return {
        items: action.paymentOverview,
      };
    case userConstants.GET_PAYMENT_OVERVIEW_FAILURE:
      return {
        error: action.error,
      };
    default:
      return state;
  }
}
