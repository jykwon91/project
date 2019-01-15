import React from 'react';
import { Link } from 'react-router-dom';
import { connect } from 'react-redux';
import { FormGroup, FormControl, ControlLabel, Navbar } from 'react-bootstrap';
import { ListGroup, ListGroupItem, ListGroupItemHeading, ListGroupItemText, ButtonGroup, Button } from 'reactstrap';
import { Container, Col, Row } from 'reactstrap';

import { userActions } from '../_actions';

class LandLordPage extends React.Component {

		constructor(props) {
			super(props);

			this.state = {
				submitted: false,
				notification: {
					Message: '',
					AddressList: [{}],
				},
				landLordPropertyList: [{}],
				serviceRequestList: [{}],
				currentUser: {},
				selectedProperty: '',
				addressList: [{}],
				fields: {},
				errors: {}
			}

			this.handleChange = this.handleChange.bind(this);
			this.handleSubmit = this.handleSubmit.bind(this);
		}

		static getDerivedStateFromProps(props, state) {
			if (props.landLordPropertyList !== state.landLordPropertyList) {
				return {
					landLordPropertyList: props.landLordPropertyList
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

		onRadioBtnClick(index, serviceReq, statusString) {
			const SRList = this.state.serviceRequestList;
			for (let i=0; i<this.state.serviceRequestList.items.length; i++) {
				if ( SRList.items[i].RequestID === serviceReq.RequestID) {
					SRList.items[i].Status = statusString;
				}
			}
		
			this.setState(prevState => ({
				serviceRequestList: SRList
			}))

			serviceReq.Status = statusString;
			this.props.dispatch(userActions.updateServiceReq(serviceReq));
		}

		handleChange(field, event) {
			const { name, value } = event.target;
			const { notification } = this.state;
			var fields = this.state.fields;
			fields[field] = event.target.value;

			this.setState({
				notification: {
					...notification,
					[name]: value
				},
				fields: fields
			});
		}

		handleSubmit(event) {
			event.preventDefault();
			event.target.reset();

			this.setState({ submitted: true});
			const { notification, addressList } = this.state;
			const { dispatch } = this.props;

			notification.AddressList = addressList;

			if (notification.Message && addressList ) {
				dispatch(userActions.sendNotification(notification));
				//reset notification
				this.setState({
					notification: {
						Message: '',
						AddressList: [],
					}
				})
				this.setState({ addressList: [] })
			}
		}

		handleDeleteUser(id) {
			return (e) => this.props.dispatch(userActions.delete(id));
		}

    componentDidMount() {
        this.props.dispatch(userActions.getAll());
				this.props.dispatch(userActions.getAllLandLordProperties());
				this.props.dispatch(userActions.getCurrentUser());
				this.props.dispatch(userActions.getServiceRequestList());
				setInterval(() => {
					this.props.dispatch(userActions.getServiceRequestList());
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

		handleAddressList = (event) => {
			const selectedAddressList = Array.from(event.target.options)
				.reduce((addressList, address) => {
					if (address.selected) {
						addressList.push(address.value);
					}

					return addressList;
				}, []);

			this.setState(() => ({ addressList: selectedAddressList }))
		}

    render() {
        const { user, users, submitted, landLordPropertyList, currentUser, addressList, serviceRequestList } = this.props;
				const { notification } = this.state;
        return (
            <div>
							<form onSubmit={this.handleSubmit.bind(this)}>
								<div className="row">
									<div className="col-xs-6 col-sm-4">
										<FormGroup controlId="formMessage">
											<ControlLabel>Please enter the notification message:</ControlLabel>
											<FormControl type="text" name="Message" value={notification.Message} componentClass="textarea" onChange={this.handleChange.bind(this, "Message")} placeholder="Enter Message..." style={{height: "100px", resize: "none"}}/>
										</FormGroup>
										<div>
										{currentUser.loading && <em>Loading properties...</em>}
										{currentUser.error && <span className="text-danger">ERROR: {currentUser.error}</span>}
										{currentUser.items &&
												<ul style={{padding: "0"}}>
													<FormGroup controlId="formAddressID">
														<ControlLabel>Send notification to these properties:</ControlLabel>
														<ControlLabel>(Ctrl+click or Click+Drag to select multiple)</ControlLabel>
														<FormControl 
															componentClass="select" 
															placeholder="Select property address" 
															onChange={this.handleAddressList}
															value={addressList} 
															style={{maxHeight: "75px"}}
														multiple>
															{currentUser.items.OwnedPropertyAddressList.map((property, index) =>
																<option key={index} value={property.AddressID}>{property.Street + ' ' + property.Zipcode + ' ' + property.State + ',' + property.City}</option>
															)}
														</FormControl>
													</FormGroup>
												</ul>
										}
										</div>
										<div style={{float: "right"}}>
											<Button type="submit">Submit</Button>
										</div>
									</div>
									<div className="col-xs-6 col-sm-4">
										<p>All registered users:</p>
										{users.loading && <em>Loading users...</em>}
										{users.error && <span className="text-danger">ERROR: {users.error}</span>}
										{users.items &&
												<ul>
														{users.items.map((user, index) =>
																<li key={index}>
																		{user.FirstName + ' ' + user.LastName}
																		{
																				user.deleting ? <em> - Deleting...</em>
																				: user.deleteError ? <span className="text-danger"> - ERROR: {user.deleteError}</span>
																				: <span> - <a onClick={this.handleDeleteUser(user.id)}>Delete</a></span>
																		}
																</li>
														)}
												</ul>
										}
									</div>
									<div className="col-xs-6 col-sm-4">
										<p> Service Requests: </p>
										{serviceRequestList.error && <span className="text-danger">ERROR: {serviceRequestList.error}</span>}
										{serviceRequestList.items &&
											<ListGroup style={{overflow: "scroll", height: "400px", width: "100%"}}>
												{serviceRequestList.items.map((serviceReq, index) =>
													<ListGroupItem key={index} color={(serviceReq.Status === 'processing') ? 'warning' : (serviceReq.Status === 'completed') ? 'success' : (serviceReq.Status === 'open') ? 'info' : 'danger'}>
														<ListGroupItemHeading>{this.convertEpochTime(serviceReq.RequestTime)}</ListGroupItemHeading>
														<ListGroupItemText style={{overflowWrap: "break-word"}}>{serviceReq.Message}</ListGroupItemText>
														<ButtonGroup size="sm">
															<Button color="info" style={{color: (serviceReq.Status === "open") ? 'black' : 'white'}} onClick={() => this.onRadioBtnClick(index, serviceReq, "open")} active={serviceReq.Status === "open"}>Open</Button>
															<Button color="warning" style={{color: (serviceReq.Status === "processing") ? 'black' : 'white'}} onClick={() => this.onRadioBtnClick(index, serviceReq, "processing")} active={serviceReq.Status === "processing"}>Processing</Button>
															<Button color="success" style={{color: (serviceReq.Status === "completed") ? 'black' : 'white'}} onClick={() => this.onRadioBtnClick(index, serviceReq, "completed")} active={serviceReq.Status === "completed"}>Completed</Button>
														</ButtonGroup>
													</ListGroupItem>
												)}
											</ListGroup>
										}
									</div>
								</div>
							</form>
								<p><Link to="/home/register">Register Rental Home</Link></p>
                <p><Link to="/login">Logout</Link></p>
            </div>
        );
    }
}

function mapStateToProps(state) {
    const { users, notification, authentication, landLordPropertyList, currentUser, serviceRequestList } = state;
    const { user } = authentication;
    return {
        user,
        users,
				notification,
				landLordPropertyList,
				currentUser,
				serviceRequestList,
    };
}

const connectedLandLordPage = connect(mapStateToProps)(LandLordPage);
export { connectedLandLordPage as LandLordPage };
