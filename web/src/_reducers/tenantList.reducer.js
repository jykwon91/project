import { userConstants } from '../_constants';

export function tenantList(state = {}, action) {
  switch (action.type) {
    case userConstants.GET_TENANT_LIST_REQUEST:
      return {
        loading: true
      };
    case userConstants.GET_TENANT_LIST_SUCCESS:
      return {
        items: action.tenantList
      };
    case userConstants.GET_TENANT_LIST_FAILURE:
      return { 
        error: action.error
      };
    case userConstants.DELETE_TENANT_REQUEST:
      // add 'deleting:true' property to user being deleted
      return {
        ...state,
        items: state.items.map(tenant =>
          tenant.id === action.id
            ? { ...tenant, deleting: true }
            : tenant
        )
      };
    case userConstants.DELETE_TENANT_SUCCESS:
      // remove deleted user from state
      return {
        items: state.items.filter(tenant => tenant.id !== action.id)
      };
    case userConstants.DELETE_TENANT_FAILURE:
      // remove 'deleting:true' property and add 'deleteError:[error]' property to user 
      return {
        ...state,
        items: state.items.map(tenant => {
          if (tenant.id === action.id) {
            // make copy of user without 'deleting:true' property
            const { deleting, ...tenantCopy } = tenant;
            // return copy of user with 'deleteError:[error]' property
            return { ...tenantCopy, deleteError: action.error };
          }

          return tenant;
        })
      };
    default:
      return state
  }
}
