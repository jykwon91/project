import { userConstants } from '../_constants';

export function createTestPayment(state = {}, action) {
  switch (action.type) {
    case userConstants.CREATE_TEST_PAYMENT_REQUEST:
      return { sending: true };
    case userConstants.CREATE_TEST_PAYMENT_SUCCESS:
      return {};
    case userConstants.CREATE_TEST_PAYMENT_FAILURE:
      return {};
    default:
      return state
  }
}
