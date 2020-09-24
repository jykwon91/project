import { userConstants } from "../_constants";

export function serviceRequestList(state = {}, action) {
  switch (action.type) {
    case userConstants.GET_SERVICEREQ_LIST_REQUEST:
      return {
        loading: true,
      };
    case userConstants.GET_SERVICEREQ_LIST_SUCCESS:
      return {
        items: action.serviceRequestList,
      };
    case userConstants.GET_SERVICEREQ_LIST_FAILURE:
      return {
        error: action.error,
      };
    default:
      return state;
  }
}
