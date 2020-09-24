import { userConstants } from "../_constants";

export function payment(state = {}, action) {
  switch (action.type) {
    case userConstants.POST_PAYMENT_SUCCESS:
      return {
        items: action.payment,
      };
    case userConstants.POST_PAYMENT_FAILURE:
      return {
        error: action.error,
      };
    default:
      return state;
  }
}
