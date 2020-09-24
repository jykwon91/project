import { userConstants } from "../_constants";

export function updateUser(state = {}, action) {
  switch (action.type) {
    case userConstants.UPDATE_USER_REQUEST:
      return { updating: true };
    case userConstants.UPDATE_USER_SUCCESS:
      return {};
    case userConstants.UPDATE_USER_FAILURE:
      return {};
    default:
      return state;
  }
}
