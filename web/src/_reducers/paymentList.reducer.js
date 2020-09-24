import { userConstants } from "../_constants";

export function paymentList(state = {}, action) {
  switch (action.type) {
    case userConstants.GET_PAYMENT_LIST_REQUEST:
      return {
        loading: true,
      };
    case userConstants.GET_PAYMENT_LIST_SUCCESS:
      return {
        items: action.paymentList,
      };
    case userConstants.GET_PAYMENT_LIST_FAILURE:
      return {
        error: action.error,
      };
    default:
      return state;
  }
}
