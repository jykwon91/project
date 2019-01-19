import React from 'react';
import { Link } from 'react-router-dom';
import { connect } from 'react-redux';
import { ListGroup, ListGroupItem, ListGroupItemHeading, ListGroupItemText, ButtonGroup, Button } from 'reactstrap';
import { FormGroup, FormControl, ControlLabel } from 'react-bootstrap';
import { Container, Col, Row } from 'reactstrap';
import DropIn from "braintree-web-drop-in-react";

import { userActions } from '../_actions';

class HomePage extends React.Component {
		instance;

		constructor(props) {
			super(props);

			this.state = {
				notifications: [{}],
				serviceReq: {
					Message: '',
					RentalAddress: {},
					TenantName: '',
				},
				serviceRequestList: [{}],
				currentUser: {},
				instance: '',
				clientToken: '',
				fields: {},
				errors: {}
			}

			this.handleChange = this.handleChange.bind(this);
			this.handleSubmit = this.handleSubmit.bind(this);
		}

		static getDerivedStateFromProps(props, state) {
			if (props.notifications !== state.notifications) {
				return {
					notifications: props.notifications
				};
			}

			if (props.currentUser !== state.currentUser) {
				return {
					currentUser: props.currentUser
				};
			}

			if (props.serviceRequestList !== state.serviceRequestList) {
				return {
					serviceRequestList: props.serviceRequestList
				};
			}

			return null;
		}

		handleSubmit(event) {
			event.preventDefault();
			event.target.reset();

			this.setState({ submitted: true });
			const { serviceReq, currentUser } = this.state;
			const { dispatch } = this.props;

			serviceReq.RentalAddress = currentUser.items.RentalAddress;
			serviceReq.TenantName = currentUser.items.FirstName + ' ' + currentUser.items.LastName

			if (serviceReq.Message) {
				dispatch(userActions.sendServiceReq(serviceReq));
				this.setState({
					serviceReq: {
						Message: '',
						RentalAddress: {},
						TenantName: '',
					},
				})
			}
		}

		handleChange(field, event) {
			const { name, value } = event.target;
			const { serviceReq } = this.state;
			var fields = this.state.fields;
			fields[field] = event.target.value;

			this.setState({
				serviceReq: {
					...serviceReq,
					[name]: value
				},
				fields: fields
			})
		}

    componentDidMount() {
			this.props.dispatch(userActions.getCurrentUser());
			this.props.dispatch(userActions.getAllNotifications());
			this.props.dispatch(userActions.getServiceRequestList());
			setInterval(() => {
				this.props.dispatch(userActions.getServiceRequestList());
				this.props.dispatch(userActions.getCurrentUser());
			}, 5000);
    }

		convertEpochTime(time) {
			var t = new Date(time * 1000);
			var arrMonth = ['January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December'];

			var hours = t.getHours()
			var afternoon = false;
			if (hours > 12) {
				hours = hours - 12;
				afternoon = true;
			}

			var minutes = t.getMinutes();
			if (minutes < 10) {
				minutes = '0' + minutes
			}
			if (!afternoon) {
				minutes = minutes + ' AM';
			} else {
				minutes = minutes + ' PM';
			}

			var formatted = arrMonth[t.getMonth()] + ' ' + t.getDate() + ' ' + t.getFullYear() + ' ' + hours + ':' + minutes;
			return formatted;
		}

		async buy() {
			const { nonce } = await this.instance.requestPaymentMethod();
			console.log(this.state.currentUser);
			var requestBody = {
				LandLordID: this.state.currentUser.items.LandLordID,
				TenantID: this.state.currentUser.items.UserID,
				RentalPaymentAmt: this.state.currentUser.items.RentalPaymentAmt,
			};

			const requestOptions = {
				method: 'POST',
				mode: 'cors',
				body: JSON.stringify(requestBody),
			}
			console.log(requestOptions);
			await fetch(`http://10.0.0.152:8000/tenant/pay/${nonce}`, requestOptions)
							.then(response => console.log(response))
		}

    render() {
        const { notifications, currentUser, serviceRequestList } = this.props;
				const { serviceReq } = this.state;
        return (
					<div>
					<Row>
						<Col xs="6" sm="4">
							<p>Notications:</p>
							{currentUser.error && <span className="text-danger">ERROR: {currentUser.error}</span>}
							{currentUser.items && currentUser.items.NotificaionList &&
								<ListGroup style={{overflow: "scroll", height: "550px", width: "100%"}}>
									{currentUser.items.NotificationList.map((notification, index) =>
										<ListGroupItem key={index}>
											<ListGroupItemHeading>{notification.Time}</ListGroupItemHeading>
											<ListGroupItemHeading>From: {notification.From}</ListGroupItemHeading>
											<ListGroupItemText>{notification.Message}</ListGroupItemText>
										</ListGroupItem>
									)}
								</ListGroup>
							}
						</Col>
						<Col xs="6" sm="4">
							<Row >
							<p>Service Requests:</p>
							{serviceRequestList.error && <span className="text-danger">ERROR: {serviceRequestList.error}</span>}
							{serviceRequestList.items &&
								<ListGroup style={{overflow: "scroll", height: "400px", width: "100%"}}>
									{serviceRequestList.items.map((serviceReq, index) =>
										<ListGroupItem key={index} color={(serviceReq.Status === "processing") ? 'warning' : (serviceReq.Status === "completed") ? 'success' : 'info'}>
											<ListGroupItemHeading>{this.convertEpochTime(serviceReq.RequestTime)}</ListGroupItemHeading>
											<ListGroupItemText style={{overflowWrap: "break-word"}}>{serviceReq.Message}</ListGroupItemText>
										</ListGroupItem>
									)}
								</ListGroup>
							}

							</Row>
							<Row>
								<form onSubmit={this.handleSubmit.bind(this)}>
									<FormGroup controlId="serviceReqForm">
										<ControlLabel>Service Request</ControlLabel>
										<FormControl type="text" name="Message" value={serviceReq.Message} componentClass="textarea" onChange={this.handleChange.bind(this, "Message")} placeholder="Enter a brief description of the issue..." style={{height: "100px", resize: "none"}}/>
									</FormGroup>
									<div style={{float: "right"}}>
										<Button type="submit">Submit</Button>
									</div>
								</form>
							</Row>
						</Col>
						<Col xs="6" sm="4">
							<p>Payments:</p>
<DropIn
            options={{ authorization: 'eyJ2ZXJzaW9uIjoyLCJhdXRob3JpemF0aW9uRmluZ2VycHJpbnQiOiJkOGY5MGU4ZDQ4NzExMzI3M2MzNDY5MmU5YjZiZDhhMjczZjBhZDM4NjdiMWM1YWRiOTg4Y2NiOWM0YTg3MDdlfGNyZWF0ZWRfYXQ9MjAxOS0wMS0xN1QwMjoyMjo0Ny4xMjcwMzY4NjgrMDAwMFx1MDAyNm1lcmNoYW50X2lkPWs1eW4ydzlzcTY5Nm43YnJcdTAwMjZwdWJsaWNfa2V5PXg4OHhicmt5enE0OWg0N2IiLCJjb25maWdVcmwiOiJodHRwczovL2FwaS5zYW5kYm94LmJyYWludHJlZWdhdGV3YXkuY29tOjQ0My9tZXJjaGFudHMvazV5bjJ3OXNxNjk2bjdici9jbGllbnRfYXBpL3YxL2NvbmZpZ3VyYXRpb24iLCJncmFwaFFMIjp7InVybCI6Imh0dHBzOi8vcGF5bWVudHMuc2FuZGJveC5icmFpbnRyZWUtYXBpLmNvbS9ncmFwaHFsIiwiZGF0ZSI6IjIwMTgtMDUtMDgifSwiY2hhbGxlbmdlcyI6WyJjdnYiXSwiZW52aXJvbm1lbnQiOiJzYW5kYm94IiwiY2xpZW50QXBpVXJsIjoiaHR0cHM6Ly9hcGkuc2FuZGJveC5icmFpbnRyZWVnYXRld2F5LmNvbTo0NDMvbWVyY2hhbnRzL2s1eW4ydzlzcTY5Nm43YnIvY2xpZW50X2FwaSIsImFzc2V0c1VybCI6Imh0dHBzOi8vYXNzZXRzLmJyYWludHJlZWdhdGV3YXkuY29tIiwiYXV0aFVybCI6Imh0dHBzOi8vYXV0aC52ZW5tby5zYW5kYm94LmJyYWludHJlZWdhdGV3YXkuY29tIiwiYW5hbHl0aWNzIjp7InVybCI6Imh0dHBzOi8vb3JpZ2luLWFuYWx5dGljcy1zYW5kLnNhbmRib3guYnJhaW50cmVlLWFwaS5jb20vazV5bjJ3OXNxNjk2bjdiciJ9LCJ0aHJlZURTZWN1cmVFbmFibGVkIjp0cnVlLCJwYXlwYWxFbmFibGVkIjp0cnVlLCJwYXlwYWwiOnsiZGlzcGxheU5hbWUiOiJSZW50YWwgTWFuYWdlbWVudCBQb3J0YWwiLCJjbGllbnRJZCI6IkFSV3VQWldDNzdIY3BZVlJoN25HTXVaU0tzd0s3eG9wSVVyUDZWWHM1UHAyLWZCY3NLdjNwTGF2eHZNYlBIZEhMc05xclRoVVpvVU1SNWVWIiwicHJpdmFjeVVybCI6Imh0dHA6Ly9leGFtcGxlLmNvbS9wcCIsInVzZXJBZ3JlZW1lbnRVcmwiOiJodHRwOi8vZXhhbXBsZS5jb20vdG9zIiwiYmFzZVVybCI6Imh0dHBzOi8vYXNzZXRzLmJyYWludHJlZWdhdGV3YXkuY29tIiwiYXNzZXRzVXJsIjoiaHR0cHM6Ly9jaGVja291dC5wYXlwYWwuY29tIiwiZGlyZWN0QmFzZVVybCI6bnVsbCwiYWxsb3dIdHRwIjp0cnVlLCJlbnZpcm9ubWVudE5vTmV0d29yayI6ZmFsc2UsImVudmlyb25tZW50Ijoib2ZmbGluZSIsInVudmV0dGVkTWVyY2hhbnQiOmZhbHNlLCJicmFpbnRyZWVDbGllbnRJZCI6Im1hc3RlcmNsaWVudDMiLCJiaWxsaW5nQWdyZWVtZW50c0VuYWJsZWQiOnRydWUsIm1lcmNoYW50QWNjb3VudElkIjoicmVudGFsbWFuYWdlbWVudHBvcnRhbCIsImN1cnJlbmN5SXNvQ29kZSI6IlVTRCJ9LCJtZXJjaGFudElkIjoiazV5bjJ3OXNxNjk2bjdiciIsInZlbm1vIjoib2ZmIn0=' }}
            onInstance={instance => (this.instance = instance)}
          />
						<Button onClick={this.buy.bind(this)}>Buy</Button>
						</Col>
					</Row>
						<p>
								<Link to="/login">Logout</Link>
						</p>
					</div>
        );
    }
}

function mapStateToProps(state) {
    const { notifications, authentication, currentUser, serviceReq, serviceRequestList } = state;
    const { user } = authentication;
    return {
        user,
        notifications,
				serviceReq,
				currentUser,
				serviceRequestList
    };
}

const connectedHomePage = connect(mapStateToProps)(HomePage);
export { connectedHomePage as HomePage };
