import { userConstants } from '../_constants';

export function sendNotification(state = {}, action) {
  switch (action.type) {
    case userConstants.SEND_NOTIFICATION_REQUEST:
      return { sending: true };
    case userConstants.SEND_NOTIFICATION_SUCCESS:
      return {};
    case userConstants.SEND_NOTIFICATION_FAILURE:
      return {};
    default:
      return state
  }
}
