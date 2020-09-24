import { userConstants } from "../_constants";

export function updateServiceReq(state = {}, action) {
  switch (action.type) {
    case userConstants.UPDATE_SERVICEREQ_REQUEST:
      return { updating: true };
    case userConstants.UPDATE_SERVICEREQ_SUCCESS:
      return {};
    case userConstants.UPDATE_SERVICEREQ_FAILURE:
      return {};
    default:
      return state;
  }
}
