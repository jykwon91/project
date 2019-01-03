import React from 'react';
import './styles.css'
import { Link } from 'react-router-dom';
import { connect } from 'react-redux';

import { userActions } from '../_actions';

class RegisterPage extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            user: {
                firstName: '',
                lastName: '',
                password: '',
								email: '',
								phoneNumber: '',
								rentalAddress: '',
								billingStreet: '',
								billingCity: '',
								billingZipcode: '',
								billingState: '',
            },
            submitted: false,
						rentalAddressList: [{value:'6738 Peerless St, Houston, TX, 77021', id:1}, {value:'Hello World', id:2}],
						showRentalAddressList: false,
						selectedRentalAddress: {value:'Select Home/Apt Address', id:99999},
						stateList: this.props.stateList || [
							{value:'TX', id:1}, 
							{value:'NO', id:2}
						],
						showStateList: false,
						selectedState: {value:'Select State', id:99999},
						fields: {},
						errors: {}
        };

        this.handleChange = this.handleChange.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
    }

		dropDownRentalAddress = () => {
			this.setState(prevState => ({
				showRentalAddressList: !prevState.showRentalAddressList
			}))
		}
		
		selectRentalAddress = (address) => this.setState(prevState => ({
			selectedRentalAddress: address,
			showRentalAddressList: false,
			user: {
				...prevState.user,
				rentalAddress: address
			}
		}))

		dropDownStates = () => {
			this.setState(prevState => ({
				showStateList: !prevState.showStateList
			}))
		}

		selectState = (_state) => this.setState(prevState => ({
			selectedState: _state,
			showStateList: false,
			user: {
				...prevState.user,
				billingState: _state
			}
		}))

		componentWillMount() {
			document.addEventListener('mousedown', this.handleClick, false);
		}

		componentWillUnmount() {
			document.removeEventListener('mousedown', this.handleClick, false);
		}

		handleClick = (e) => {
			if (this.node.contains(e.target)) {
				this.setState({
					showStateList: false
				})
				return;
			} else if (this.tempnode.contains(e.target)) {
				this.setState({
					showRentalAddressList: false
				})
				return;
			}
			this.handleClickOutside();
		}

		handleClickOutside() {
			this.setState({
				showRentalAddressList: false,
				showStateList: false
			})
		}

    handleChange(field, event) {
        const { name, value } = event.target;
        const { user } = this.state;
				var fields = this.state.fields;
				fields[field] = event.target.value;


        this.setState({
            user: {
                ...user,
                [name]: value
            },
						fields: fields
        });
    }

    handleSubmit(event) {
        event.preventDefault();

        this.setState({ submitted: true });
        const { user } = this.state;
        const { dispatch } = this.props;

        if (this.handleValidation() && user.firstName && user.lastName && user.password && user.email && user.rentalAddress && user.billingStreet && user.billingCity && user.billingZipcode && user.billingState) {
            dispatch(userActions.register(user));
        }
    }

		handleValidation() {
			const fields = this.state.fields;
			const errors = {};
			var formIsValid = true;

			//phone number
			if (!fields["phoneNumber"]) {
				formIsValid = false;
				errors["phoneNumber"] = "Phone number is required";
			}

			if (typeof fields["phoneNumber"] !== "undefined") {
				if (!fields["phoneNumber"].match(/^[0-9]+$/)) {
					formIsValid = false;
					errors["phoneNumber"] = "Phone number can only consist of numbers(0-9)";
				}
			}

			//first name
			if (!fields["firstName"]) {
				formIsValid = false;
				errors["firstName"] = "First name is required";
			}

			if (typeof fields["firstName"] !== "undefined") {
				if (!fields["firstName"].match(/^[a-zA-Z]+$/)) {
					formIsValid = false;
					errors["firstName"] = "Name can only consist of letters(a-z)";
				}
			}

			//last name
			if (!fields["lastName"]) {
				formIsValid = false;
				errors["lastName"] = "Last name is required";
			}

			if (typeof fields["lastName"] !== "undefined") {
				if (!fields["lastName"].match(/^[a-zA-Z]+$/)) {
					formIsValid = false;
					errors["lastName"] = "Name can only consist of letters(a-z)";
				}
			}

			//email
			if (!fields["email"]) {
				formIsValid = false;
				errors["email"] = "Email is required";
			}

			if (typeof fields["email"] !== "undefined") {
				var lastAtPos = fields["email"].lastIndexOf('@');
				var lastDotPos = fields["email"].lastIndexOf('.');

				if (!(lastAtPos < lastDotPos && lastAtPos > 0 && fields["email"].indexOf('@@') == -1 && lastDotPos > 2 && (fields["email"].length - lastDotPos) > 2)) {
					formIsValid = false;
					errors["email"] = "Email is not valid";
				}
			}

			this.setState({errors: errors});
			return formIsValid;
		}

    render() {
        const { registering  } = this.props;
        const { user, submitted } = this.state;
        return (
            <div className="col-md-10 col-md-offset-1">
                <h2>Register</h2>
                <form name="form" onSubmit={this.handleSubmit}>
                    <div className={'form-group' + (submitted && this.state.errors["firstName"] ? ' has-error' : '')}>
                        <label htmlFor="firstName">First Name</label>
                        <input type="text" className="form-control" name="firstName" value={user.firstName} onChange={this.handleChange.bind(this, "firstName")} />
                        <div className="help-block">{this.state.errors["firstName"]}</div>
                    </div>
                    <div className={'form-group' + (submitted && this.state.errors["lastName"] ? ' has-error' : '')}>
                        <label htmlFor="lastName">Last Name</label>
                        <input type="text" className="form-control" name="lastName" value={user.lastName} onChange={this.handleChange.bind(this, "lastName")} />
                        <div className="help-block">{this.state.errors["lastName"]}</div>
                    </div>
                    <div className={'form-group' + (submitted && !user.password ? ' has-error' : '')}>
                        <label htmlFor="password">Password</label>
                        <input type="password" className="form-control" name="password" value={user.password} onChange={this.handleChange.bind(this, "password")} />
                        {submitted && !user.password &&
                            <div className="help-block">Password is required</div>
                        }
                    </div>
                    <div className={'form-group' + (submitted && this.state.errors["email"] ? ' has-error' : '')}>
                        <label htmlFor="email">Email</label>
                        <input type="text" className="form-control" name="email" value={user.email} onChange={this.handleChange.bind(this, "email")} />
                        <div className="help-block">{this.state.errors["email"]}</div>
                    </div>
                    <div className={'form-group' + (submitted && this.state.errors["phoneNumber"] ? ' has-error' : '')}>
                        <label htmlFor="phoneNumber">Phone Number</label>
                        <input type="tel" className="form-control" name="phoneNumber" value={user.phoneNumber} onChange={this.handleChange.bind(this, "phoneNumber")} maxLength="10" />
                        <div className="help-block">{this.state.errors["phoneNumber"]}</div>
                    </div>



                		<h3>Billing Address</h3>
										<div className="row">
											<div className="col-xs-6">
												<div className={'form-group' + (submitted && !user.billingStreet ? ' has-error' : '')}>
														<label htmlFor="billingStreet">Street</label>
														<input type="text" className="form-control" name="billingStreet" value={user.billingStreet} onChange={this.handleChange.bind(this, "billingStreet")} />
														{submitted && !user.billingStreet &&
																<div className="help-block">Street is required</div>
														}
												</div>
											</div>
											<div className="col-xs-6">
												<div className={'form-group' + (submitted && !user.billingCity ? ' has-error' : '')}>
														<label htmlFor="billingCity">City</label>
														<input type="text" className="form-control" name="billingCity" value={user.billingCity} onChange={this.handleChange.bind(this, "billingCity")} />
														{submitted && !user.billingCity &&
																<div className="help-block">City is required</div>
														}
												</div>
											</div>
											<div className="col-xs-6">
												<div className={'form-group' + (submitted && !user.billingZipcode ? ' has-error' : '')}>
														<label htmlFor="billingZipcode">Zipcode</label>
														<input type="text" className="form-control" name="billingZipcode" value={user.billingZipcode} onChange={this.handleChange.bind(this, "billingZipcode")} />
														{submitted && !user.billingZipcode &&
																<div className="help-block">Zipcode is required</div>
														}
												</div>
											</div>
											<div className="col-xs-6">
                        <label htmlFor="billingState">State</label>
                    		<div className={'form-group' + (submitted && !user.billingState ? ' has-error' : '')} ref={tempnode => this.tempnode = tempnode}>
													<div className="select-box--box" style={{width: this.state.width || "100%"}}>
														<div className="select-box--container">
															<div
																className="select-box--selected-item"
																onClick={this.dropDownStates}
															> { this.state.selectedState.value }
															</div>
															<div 
																className="select-box--arrow"
																onClick={this.dropDownStates}>
																	<span className={`${this.state.showStateList ? 'select-box--arrow-up' : 'select-box--arrow-down'}`}/>
															</div>
															<div 
																style={{display: this.state.showStateList ? 'block' : 'none'}} 
																className="select-box--items"
															>
															{
																this.state.stateList.map(_state => <div 
																	key={ _state.id } 
																	onClick={() => this.selectState(_state)} 
																	className={this.state.selectedState === _state ? 'selected' : ''}
																>
																	{ _state.value }
																</div>)
															}
															</div>
														</div>
														<input
															type="hidden"
															value={user.billingState}
															name="billingState"
															onChange={this.handleChange.bind(this, "billingState")}
														/>
														{submitted && !user.billingState &&
															<div className="help-block"> State is required</div>
														}
													</div>
												</div>
											</div>
										</div>

                    <div className={'form-group' + (submitted && !user.rentalAddress ? ' has-error' : '')} ref={node => this.node = node}>
                      <label htmlFor="rentalAddress">Rental Address</label>
											<div className="select-box--box" style={{width: this.state.width || "100%"}}>
												<div className="select-box--container">
													<div
														className="select-box--selected-item"
														onClick={this.dropDownRentalAddress}
													> { this.state.selectedRentalAddress.value }
													</div>
													<div 
														className="select-box--arrow"
														onClick={this.dropDownRentalAddress}>
															<span className={`${this.state.showRentalAddressList ? 'select-box--arrow-up' : 'select-box--arrow-down'}`}/>
													</div>
													<div 
														style={{display: this.state.showRentalAddressList ? 'block' : 'none'}} 
														className="select-box--items"
													>
													{
														this.state.rentalAddressList.map(address => <div 
															key={ address.id } 
															onClick={() => this.selectRentalAddress(address)} 
															className={this.state.selectedRentalAddress === address ? 'selected' : ''}
														>
															{ address.value }
														</div>)
													}
													</div>
												</div>
												<input
													type="hidden"
													value={user.rentalAddress}
													name="rentalAddress"
													onChange={this.handleChange.bind(this, "rentalAddress")}
												/>
												{submitted && !user.rentalAddress &&
													<div className="help-block"> Rental address is required</div>
												}
											</div>
										</div>

                    <div className="form-group">
                        <button className="btn btn-primary">Register</button>
                        {registering && 
                            <img src="data:image/gif;base64,R0lGODlhEAAQAPIAAP///wAAAMLCwkJCQgAAAGJiYoKCgpKSkiH/C05FVFNDQVBFMi4wAwEAAAAh/hpDcmVhdGVkIHdpdGggYWpheGxvYWQuaW5mbwAh+QQJCgAAACwAAAAAEAAQAAADMwi63P4wyklrE2MIOggZnAdOmGYJRbExwroUmcG2LmDEwnHQLVsYOd2mBzkYDAdKa+dIAAAh+QQJCgAAACwAAAAAEAAQAAADNAi63P5OjCEgG4QMu7DmikRxQlFUYDEZIGBMRVsaqHwctXXf7WEYB4Ag1xjihkMZsiUkKhIAIfkECQoAAAAsAAAAABAAEAAAAzYIujIjK8pByJDMlFYvBoVjHA70GU7xSUJhmKtwHPAKzLO9HMaoKwJZ7Rf8AYPDDzKpZBqfvwQAIfkECQoAAAAsAAAAABAAEAAAAzMIumIlK8oyhpHsnFZfhYumCYUhDAQxRIdhHBGqRoKw0R8DYlJd8z0fMDgsGo/IpHI5TAAAIfkECQoAAAAsAAAAABAAEAAAAzIIunInK0rnZBTwGPNMgQwmdsNgXGJUlIWEuR5oWUIpz8pAEAMe6TwfwyYsGo/IpFKSAAAh+QQJCgAAACwAAAAAEAAQAAADMwi6IMKQORfjdOe82p4wGccc4CEuQradylesojEMBgsUc2G7sDX3lQGBMLAJibufbSlKAAAh+QQJCgAAACwAAAAAEAAQAAADMgi63P7wCRHZnFVdmgHu2nFwlWCI3WGc3TSWhUFGxTAUkGCbtgENBMJAEJsxgMLWzpEAACH5BAkKAAAALAAAAAAQABAAAAMyCLrc/jDKSatlQtScKdceCAjDII7HcQ4EMTCpyrCuUBjCYRgHVtqlAiB1YhiCnlsRkAAAOwAAAAAAAAAAAA==" />
                        }
                        <Link to="/login" className="btn btn-link">Cancel</Link>
                    </div>
                </form>
            </div>
        );
    }
}

function mapStateToProps(state) {
    const { registering } = state.registration;
    return {
        registering
    };
}

const connectedRegisterPage = connect(mapStateToProps)(RegisterPage);
export { connectedRegisterPage as RegisterPage };
