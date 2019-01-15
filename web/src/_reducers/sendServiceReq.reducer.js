import { userConstants } from '../_constants';

export function sendServiceReq(state = {}, action) {
  switch (action.type) {
    case userConstants.SEND_SERVICEREQ_REQUEST:
      return { sending: true };
    case userConstants.SEND_SERVICEREQ_SUCCESS:
      return {};
    case userConstants.SEND_SERVICEREQ_FAILURE:
      return {};
    default:
      return state
  }
}
