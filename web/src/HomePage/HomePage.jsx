import React from 'react';
import { Link } from 'react-router-dom';
import { connect } from 'react-redux';
import { ListGroup, ListGroupItem, ListGroupItemHeading, ListGroupItemText, ButtonGroup, Button } from 'reactstrap';
import { FormGroup, FormControl, ControlLabel } from 'react-bootstrap';
import { Container, Col, Row } from 'reactstrap';

import { userActions } from '../_actions';

class HomePage extends React.Component {

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

    render() {
        const { notifications, currentUser, serviceRequestList } = this.props;
				const { serviceReq } = this.state;
        return (
					<div>
					<Row>
						<Col xs="6" sm="4">
							<p>Notications:</p>
							{currentUser.error && <span className="text-danger">ERROR: {currentUser.error}</span>}
							{currentUser.items &&
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
