import { combineReducers } from 'redux';

import { authentication } from './authentication.reducer';
import { registration } from './registration.reducer';
import { users } from './users.reducer';
import { notifications } from './notifications.reducer';
import { stateList } from './stateList.reducer';
import { landLordPropertyList } from './landLordPropertyList.reducer';
import { landLordPropertyRegistration } from './landLordPropertyRegistration.reducer';
import { currentUser } from './currentUser.reducer';
import { sendNotification } from './sendNotification.reducer';
import { sendServiceReq } from './sendServiceReq.reducer';
import { serviceRequestList } from './serviceRequestList.reducer';
import { updateServiceReq } from './updateServiceReq.reducer';
import { landLordList } from './landLordList.reducer';
import { tenantList } from './tenantList.reducer';
import { updateUser } from './updateUser.reducer';
import { paymentList } from './paymentList.reducer';
import { payment } from './payment.reducer';
import { paymentOverview } from './paymentOverview.reducer';
import { createTestPayment } from './createTestPayment.reducer';
import { alert } from './alert.reducer';

const rootReducer = combineReducers({
  authentication,
  registration,
  users,
	notifications,
	stateList,
	landLordPropertyList,
	landLordPropertyRegistration,
	currentUser,
	sendNotification,
	sendServiceReq,
	serviceRequestList,
	updateServiceReq,
	landLordList,
	tenantList,
	updateUser,
	paymentList,
	payment,
	paymentOverview,
	createTestPayment,
  alert
});

export default rootReducer;
