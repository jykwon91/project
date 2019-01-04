import { combineReducers } from 'redux';

import { authentication } from './authentication.reducer';
import { registration } from './registration.reducer';
import { users } from './users.reducer';
import { notifications } from './notifications.reducer';
import { stateList } from './stateList.reducer';
import { landLordPropertyRegistration } from './landLordPropertyRegistration.reducer';
import { alert } from './alert.reducer';

const rootReducer = combineReducers({
  authentication,
  registration,
  users,
	notifications,
	stateList,
	landLordPropertyRegistration,
  alert
});

export default rootReducer;
