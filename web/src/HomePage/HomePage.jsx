import React from 'react';
import axios from 'axios';
import { Link } from 'react-router-dom';
import { authHeader } from '../_helpers';
import { connect } from 'react-redux';
import { ListGroup, ListGroupItem, ListGroupItemHeading, ListGroupItemText, ButtonGroup, Button } from 'reactstrap';
import { FormGroup, FormControl, ControlLabel, Modal } from 'react-bootstrap';
import { Container, Col, Row } from 'reactstrap';
import CurrencyInput from 'react-currency-input';
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
				selectedPayment: {},
				paymentAmount: 0,
				serviceRequestList: [{}],
				paymentList: {loading: true},
				currentUser: '',
				instance: '',
				clientToken: '',
				modal: false,
				nestedModal: false,
				closeAll: false,
				paymentOverview: {},
				submitted: false,
				fields: {},
				errors: {}
			}

			this.toggle = this.toggle.bind(this);
			this.toggleNested = this.toggleNested.bind(this);
			this.toggleAll = this.toggleAll.bind(this);
			this.fetchCurrentUser = this.fetchCurrentUser.bind(this);

			this.handleChange = this.handleChange.bind(this);
			this.handleSubmit = this.handleSubmit.bind(this);
		}

		openModal(payment) {
			this.setState({
				selectedPayment: payment,
				modal: !this.state.modal,
			});
		}

		toggle() {
			this.setState({
				modal: !this.state.modal
			});
		}

		toggleNested() {
			this.setState({
				nestedModal: !this.state.nestedModal,
				closeAll: false
			});
		}

		toggleAll() {
			this.setState({
				nestedModal: !this.state.nestedModal,
				closeAll: true,
			});
		}

		static getDerivedStateFromProps(props, state) {
			if (props.notifications !== state.notifications) {
				return {
					notifications: props.notifications
				};
			}

			if ((props.currentUser !== state.currentUser) && !props.currentUser) {
				return {
					currentUser: props.currentUser
				};
			}

			if (props.serviceRequestList !== state.serviceRequestList) {
				return {
					serviceRequestList: props.serviceRequestList
				};
			}

			if (props.paymentList !== state.paymentList) {
				return {
					paymentList: props.paymentList
				};
			}

			if (props.paymentOverview !== state.paymentOverview) {
				return {
					paymentOverview: props.paymentOverview
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

		handleValidation(payment) {
			const fields = this.state.fields;
			const errors = {};
			var formIsValid = true;

			if (!fields["paymentAmount"]) {
				formIsValid = false;
				errors["paymentAmount"] = "Payment amount must be greater than 0";
			}

			if (typeof fields["paymentAmount"] !== "undefined") {
				if (fields["paymentAmount"] >= 0) {
					formIsValid = false;
					errors["paymentAmount"] = "Payment amount must be greater than 0";
				}

				if (this.convertDollarToInt(fields["paymentAmount"]) > payment.Amount) {
					formIsValid = false;
					errors["paymentAmount"] = "Payment amount cannot be greater than the amount due(" + this.convertDollarAmount(payment.Amount) + ")";
				}
			}

			this.setState({errors: errors});
			return formIsValid;
		}

		handleChange(field, event) {
			const { name, value } = event.target;
			const { serviceReq } = this.state;
			const { paymentAmount } = this.state;
			var fields = this.state.fields;
			fields[field] = event.target.value;

			this.setState({
				serviceReq: {
					...serviceReq,
					[name]: value
				},
				paymentAmount: value,
				fields: fields
			})
		}

		fetchCurrentUser = () => {
			const requestOptions = {
				method: 'GET',
				mode: 'cors',
				headers: authHeader()
			};
			fetch (`http://rentalmgmt.co:8000/users/currentUser`, requestOptions)
				.then( response => {
					response.text().then( text => {
						const data = text && JSON.parse(text);
						this.setState({currentUser: data});
						this.props.dispatch(userActions.getPaymentOverview(data.LandLordID));
					})
				})
		}

    async componentDidMount() {

			await this.fetchCurrentUser();
			this.props.dispatch(userActions.getAllNotifications());
			this.props.dispatch(userActions.getServiceRequestList());
			this.props.dispatch(userActions.getPaymentList());
			setInterval(() => {
				this.props.dispatch(userActions.getCurrentUser());
				this.props.dispatch(userActions.getPaymentOverview(this.state.currentUser.LandLordID));
				this.props.dispatch(userActions.getServiceRequestList());
				this.props.dispatch(userActions.getPaymentList());
			}, 5000);
    }

		convertDollarAmount(amount) {
			return parseFloat(Math.round(amount) / 100).toFixed(2);
		}

		convertDollarToInt(amount) {
			return Number(amount.substr(1).replace(',','')) * 100;
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
			const { dispatch } = this.props;
			const { paymentAmount, currentUser, selectedPayment } = this.state;
			this.setState({submitted: true});
			if (this.handleValidation(this.state.selectedPayment) && paymentAmount) {
				const { nonce } = await this.instance.requestPaymentMethod();

				var requestBody = {
					PaymentID: selectedPayment.PaymentID,
					LandLordID: currentUser.LandLordID,
					TenantID: currentUser.UserID,
					Amount: this.convertDollarToInt(paymentAmount),
				};

				const requestOptions = {
					method: 'POST',
					mode: 'cors',
					body: JSON.stringify(requestBody),
				}
				await fetch(`http://rentalmgmt.co:8000/tenant/pay/${nonce}`, requestOptions)
								.then(response => {
									this.setState({
										modal: !this.state.modal
									});
									response.text().then( text => {
										const data = text && JSON.parse(text);
										console.log(data);
										console.log(text);
										console.log(response);
										if (!data) {
											this.props.dispatch(userActions.payment(null));
										} else if (data.type === "error") {
											this.props.dispatch(userActions.payment(data));
										}
									})
								})
			}
		}

    render() {
        const { notifications, currentUser, serviceRequestList, paymentList, paymentOverview } = this.props;
				const { submitted, serviceReq, selectedPayment, paymentAmount } = this.state;
        return (
					<div>
					<h1>{currentUser.UserID}</h1>
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
							<Row>
								<div>
									{paymentOverview.items && 
										<div>
											<h6>Current Pay Period: {paymentOverview.items.CurrentPayPeriod}</h6>
											<h6>Current Amount Due: ${this.convertDollarAmount(paymentOverview.items.CurrentAmountDue)}</h6>
											<h6>Late Amount: ${this.convertDollarAmount(paymentOverview.items.TotalLateAmount)}</h6>
											<h6>Late Fees: ${this.convertDollarAmount(paymentOverview.items.TotalLateFees)}</h6>
											<h6>Total Due: ${this.convertDollarAmount(paymentOverview.items.TotalDue)}</h6>
										</div>
									}
								</div>
							</Row>
							<Row>
									{paymentList.loading && <em>Loading Payment List...</em>}
									{paymentList.error && <span className="text-danger">ERROR: {paymentList.error}</span>}
									{paymentList.items &&
										<ListGroup style={{overflow: "scroll", height: "550px", width: "100%"}}>
											{paymentList.items.map((payment, index) =>
												<ListGroupItem key={index} color={(payment.Status === "processing") ? 'warning' : (payment.Status === "paid") ? 'success' : (payment.Status === "late") ? 'danger' : 'info'}>
													<div>
														<div style={{float: "left"}}></div>
														<div style={{float: "right"}}><Button color="info" onClick={() => this.openModal(payment)}>Pay</Button></div>
													</div>
													<ListGroupItemHeading>{this.convertEpochTime(payment.PaidDate)}</ListGroupItemHeading>
													<ListGroupItemHeading>Amount: ${this.convertDollarAmount(payment.Amount)}</ListGroupItemHeading>
													<ListGroupItemHeading>Category: {payment.Category}</ListGroupItemHeading>
													<ListGroupItemHeading>Status: {payment.Status}</ListGroupItemHeading>
												</ListGroupItem>
											)}
										</ListGroup>
									}
							</Row>
						</Col>
					</Row>
						<p>
								<Link to="/login">Logout</Link>
						</p>
								<Modal show={this.state.modal} style={{opacity: "1"}} onHide={this.toggle}>
									<Modal.Header closeButton><Modal.Title>Payment</Modal.Title></Modal.Header>
										<Modal.Body>
													<p>Category: {selectedPayment.Category}</p>
													<p>Description: {selectedPayment.Description}</p>
													<p>Date Posted: {this.convertEpochTime(selectedPayment.DueDate)}</p>
													<p>Amount Due: ${this.convertDollarAmount(selectedPayment.Amount)}</p>
													<div className={'form-group' + (submitted && this.state.errors["paymentAmount"] ? ' has-error' : '')}>
														<CurrencyInput inputType="text" prefix="$" name="paymentAmount" value={paymentAmount} onChangeEvent={this.handleChange.bind(this, "paymentAmount")}/>
														<div className="help-block">{this.state.errors["paymentAmount"]}</div>
													</div>
													<DropIn
														options={{ authorization: 'eyJ2ZXJzaW9uIjoyLCJhdXRob3JpemF0aW9uRmluZ2VycHJpbnQiOiJkOGY5MGU4ZDQ4NzExMzI3M2MzNDY5MmU5YjZiZDhhMjczZjBhZDM4NjdiMWM1YWRiOTg4Y2NiOWM0YTg3MDdlfGNyZWF0ZWRfYXQ9MjAxOS0wMS0xN1QwMjoyMjo0Ny4xMjcwMzY4NjgrMDAwMFx1MDAyNm1lcmNoYW50X2lkPWs1eW4ydzlzcTY5Nm43YnJcdTAwMjZwdWJsaWNfa2V5PXg4OHhicmt5enE0OWg0N2IiLCJjb25maWdVcmwiOiJodHRwczovL2FwaS5zYW5kYm94LmJyYWludHJlZWdhdGV3YXkuY29tOjQ0My9tZXJjaGFudHMvazV5bjJ3OXNxNjk2bjdici9jbGllbnRfYXBpL3YxL2NvbmZpZ3VyYXRpb24iLCJncmFwaFFMIjp7InVybCI6Imh0dHBzOi8vcGF5bWVudHMuc2FuZGJveC5icmFpbnRyZWUtYXBpLmNvbS9ncmFwaHFsIiwiZGF0ZSI6IjIwMTgtMDUtMDgifSwiY2hhbGxlbmdlcyI6WyJjdnYiXSwiZW52aXJvbm1lbnQiOiJzYW5kYm94IiwiY2xpZW50QXBpVXJsIjoiaHR0cHM6Ly9hcGkuc2FuZGJveC5icmFpbnRyZWVnYXRld2F5LmNvbTo0NDMvbWVyY2hhbnRzL2s1eW4ydzlzcTY5Nm43YnIvY2xpZW50X2FwaSIsImFzc2V0c1VybCI6Imh0dHBzOi8vYXNzZXRzLmJyYWludHJlZWdhdGV3YXkuY29tIiwiYXV0aFVybCI6Imh0dHBzOi8vYXV0aC52ZW5tby5zYW5kYm94LmJyYWludHJlZWdhdGV3YXkuY29tIiwiYW5hbHl0aWNzIjp7InVybCI6Imh0dHBzOi8vb3JpZ2luLWFuYWx5dGljcy1zYW5kLnNhbmRib3guYnJhaW50cmVlLWFwaS5jb20vazV5bjJ3OXNxNjk2bjdiciJ9LCJ0aHJlZURTZWN1cmVFbmFibGVkIjp0cnVlLCJwYXlwYWxFbmFibGVkIjp0cnVlLCJwYXlwYWwiOnsiZGlzcGxheU5hbWUiOiJSZW50YWwgTWFuYWdlbWVudCBQb3J0YWwiLCJjbGllbnRJZCI6IkFSV3VQWldDNzdIY3BZVlJoN25HTXVaU0tzd0s3eG9wSVVyUDZWWHM1UHAyLWZCY3NLdjNwTGF2eHZNYlBIZEhMc05xclRoVVpvVU1SNWVWIiwicHJpdmFjeVVybCI6Imh0dHA6Ly9leGFtcGxlLmNvbS9wcCIsInVzZXJBZ3JlZW1lbnRVcmwiOiJodHRwOi8vZXhhbXBsZS5jb20vdG9zIiwiYmFzZVVybCI6Imh0dHBzOi8vYXNzZXRzLmJyYWludHJlZWdhdGV3YXkuY29tIiwiYXNzZXRzVXJsIjoiaHR0cHM6Ly9jaGVja291dC5wYXlwYWwuY29tIiwiZGlyZWN0QmFzZVVybCI6bnVsbCwiYWxsb3dIdHRwIjp0cnVlLCJlbnZpcm9ubWVudE5vTmV0d29yayI6ZmFsc2UsImVudmlyb25tZW50Ijoib2ZmbGluZSIsInVudmV0dGVkTWVyY2hhbnQiOmZhbHNlLCJicmFpbnRyZWVDbGllbnRJZCI6Im1hc3RlcmNsaWVudDMiLCJiaWxsaW5nQWdyZWVtZW50c0VuYWJsZWQiOnRydWUsIm1lcmNoYW50QWNjb3VudElkIjoicmVudGFsbWFuYWdlbWVudHBvcnRhbCIsImN1cnJlbmN5SXNvQ29kZSI6IlVTRCJ9LCJtZXJjaGFudElkIjoiazV5bjJ3OXNxNjk2bjdiciIsInZlbm1vIjoib2ZmIn0=' }}
														onInstance={instance => (this.instance = instance)}
													/>
											<Button color="success" onClick={this.toggleNested}>Show Nested Modal</Button>
											<Modal show={this.state.nestedModal} style={{opacity: "1"}} onHide={this.toggleNested}>
												<Modal.Header>Nested Modal title</Modal.Header>
												<Modal.Body>
												</Modal.Body>
												<Modal.Footer>
													<Button color="primary" onClick={this.toggleNested}>Done</Button>{' '}
													<Button color="secondary" onClick={this.toggleAll}>All Done</Button>
												</Modal.Footer>
											</Modal>
										</Modal.Body>
									<Modal.Footer>
										<Button color="primary" onClick={this.buy.bind(this)}>Pay Now</Button>
										<Button color="secondary" onClick={this.toggle}>Cancel</Button>
									</Modal.Footer>
								</Modal>
					</div>
        );
    }
}

function mapStateToProps(state) {
    const { notifications, authentication, currentUser, serviceReq, serviceRequestList, paymentList, paymentOverview } = state;
    const { user } = authentication;
    return {
        user,
        notifications,
				serviceReq,
				currentUser,
				serviceRequestList,
				paymentList,
				paymentOverview,
    };
}

const connectedHomePage = connect(mapStateToProps)(HomePage);
export { connectedHomePage as HomePage };

/*
            <Modal show={this.state.nestedModal} style={{opacity: "1"}} onHide={this.toggleNested} onClosed={this.state.closeAll ? this.toggle : undefined}>
*/
