import React from 'react';
import './styles.css'
import { Link } from 'react-router-dom';
import { connect } from 'react-redux';

import CurrencyInput from 'react-currency-input';
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
								landLord: '',
								billingStreet: '',
								billingCity: '',
								billingZipcode: '',
								billingState: '',
								rentalPaymentAmt: '',
            },
            submitted: false,
						selectedLandLord: {Name: 'Select your land lord', LandLordID:99999},
						stateList: props.stateList,
						landLordList: props.landLordList,
						showStateList: false,
						showLandLordList: false,
						selectedState: {value:'Select State', id:99999},
						fields: {},
						errors: {}
        };

        this.handleChange = this.handleChange.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
				this.handleClickSL = this.handleClickSL.bind(this);
				this.handleOutsideClickSL = this.handleOutsideClickSL.bind(this);
				this.handleClickLLL = this.handleClickLLL.bind(this);
				this.handleOutsideClickLLL = this.handleOutsideClickLLL.bind(this);
    }

		dropDownLandLordList = () => {
			this.setState(prevState => ({
				showLandLordList: !prevState.showLandLordList
			}))
		}
		
		selectLandLord = (landLord) => this.setState(prevState => ({
			selectedLandLord: landLord,
			showLandLordList: false,
			user: {
				...prevState.user,
				landLord: landLord
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

		static getDerivedStateFromProps(props, state) {
			if (props.stateList !== state.stateList) {
				return {
					stateList: props.stateList
				};
			}

			if (props.landLordList !== state.landLordList) {
				return {
					landLordList: props.landLordList
				};
			};

			return null;
		}

		componentDidMount() {
			this.props.dispatch(userActions.getStateList());
			this.props.dispatch(userActions.getLandLordList());
		}

		convertDollarToInt(amount) {
			return Number(amount.substr(1).replace(',','')) * 100;
		}

		handleClickSL() {
			if (!this.state.showStateList) {
				document.addEventListener('click', this.handleOutsideClickSL, false);
			} else {
				document.removeEventListener('click', this.handleOutsideClickSL, false);
			}

			this.setState({ showStateList: !this.state.showStateList});
		}

		handleOutsideClickSL(e) {
			if (this.SLNode.contains(e.target)) {
				return;
			}

			this.handleClickSL();
		}

		handleClickLLL() {
			if (!this.state.showLandLordList) {
				document.addEventListener('click', this.handleOutsideClickLLL, false);
			} else {
				document.removeEventListener('click', this.handleOutsideClickLLL, false);
			}

			this.setState({ showLandLordList: !this.state.showLandLordList});
		}

		handleOutsideClickLLL(e) {
			if (this.LLLNode.contains(e.target)) {
				return;
			}

			this.handleClickLLL();
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

				console.log(this.state.user);
        this.setState({ submitted: true });
        const { user } = this.state;
        const { dispatch } = this.props;

        if (this.handleValidation() && user.firstName && user.lastName && user.password && user.email && user.landLord && user.billingStreet && user.billingCity && user.billingZipcode && user.billingState) {
					user.rentalPaymentAmt = this.convertDollarToInt(user.rentalPaymentAmt);
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

			//billingZipcode
			if (!fields["billingZipcode"]) {
				formIsValid = false;
				errors["billingZipcode"] = "Zipcode is required";
			}

			if (typeof fields["billingZipcode"] !== "undefined") {
				if (!fields["billingZipcode"].match(/^[0-9]+$/)) {
					formIsValid = false;
					errors["billingZipcode"] = "Zipcode can only consist of numbers(0-9)";
				}
			}

			//rentalPaymentAmt
			if (!fields["rentalPaymentAmt"]) {
				formIsValid = false;
				errors["rentalPaymentAmt"] = "Monthly rental payment is required";
			}

			if (typeof fields["rentalPaymentAmt"] !== "undefined") {
				if (fields["rentalPaymentAmt"] >= 0) {
					formIsValid = false;
					errors["rentalPaymentAmt"] = "Monthly rental payment is required";
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
        const { registering, stateList, landLordList } = this.props;
        const { user, submitted } = this.state;
        return (
            <div className="col-md-10 col-md-offset-1">
                <h2>Register</h2>
                <form name="form" onSubmit={this.handleSubmit}>
									<div className="row">
										<div className="col-xs-6">
											<div className={'form-group' + (submitted && this.state.errors["firstName"] ? ' has-error' : '')}>
													<label htmlFor="firstName">First Name</label>
													<input type="text" className="form-control" name="firstName" value={user.firstName} onChange={this.handleChange.bind(this, "firstName")} />
													<div className="help-block">{this.state.errors["firstName"]}</div>
											</div>
										</div>
										<div className="col-xs-6">
											<div className={'form-group' + (submitted && this.state.errors["lastName"] ? ' has-error' : '')}>
													<label htmlFor="lastName">Last Name</label>
													<input type="text" className="form-control" name="lastName" value={user.lastName} onChange={this.handleChange.bind(this, "lastName")} />
													<div className="help-block">{this.state.errors["lastName"]}</div>
											</div>
										</div>
										<div className="col-xs-6">
											<div className={'form-group' + (submitted && !user.password ? ' has-error' : '')}>
													<label htmlFor="password">Password</label>
													<input type="password" className="form-control" name="password" value={user.password} onChange={this.handleChange.bind(this, "password")} />
													{submitted && !user.password &&
															<div className="help-block">Password is required</div>
													}
											</div>
										</div>
										<div className="col-xs-6">
											<div className={'form-group' + (submitted && this.state.errors["email"] ? ' has-error' : '')}>
													<label htmlFor="email">Email</label>
													<input type="text" className="form-control" name="email" value={user.email} onChange={this.handleChange.bind(this, "email")} />
													<div className="help-block">{this.state.errors["email"]}</div>
											</div>
										</div>
										<div className="col-xs-6">
											<div className={'form-group' + (submitted && this.state.errors["phoneNumber"] ? ' has-error' : '')}>
													<label htmlFor="phoneNumber">Phone Number</label>
													<input type="tel" className="form-control" name="phoneNumber" value={user.phoneNumber} onChange={this.handleChange.bind(this, "phoneNumber")} maxLength="10" />
													<div className="help-block">{this.state.errors["phoneNumber"]}</div>
											</div>
										</div>
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
												<div className={'form-group' + (submitted && this.state.errors["billingZipcode"] ? ' has-error' : '')} style={{width: "100px"}}>
														<label htmlFor="billingZipcode">Zipcode</label>
														<input type="text" className="form-control" name="billingZipcode" value={user.billingZipcode} onChange={this.handleChange.bind(this, "billingZipcode")} maxLength="5" />
														<div className="help-block">{this.state.errors["billingZipcode"]}</div>
												</div>
											</div>
											<div className="col-xs-6">
                        <label htmlFor="billingState">State</label>
                    		<div className={'form-group' + (submitted && !user.billingState ? ' has-error' : '')} ref={SLNode => this.SLNode = SLNode}>
													<div className="select-box--box" style={{width: this.state.width || "150px"}}>
														<div className="select-box--container" onClick={this.handleClickSL}>
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
																style={{display: this.state.showStateList ? 'block' : 'none', maxHeight: '400px', overflow: 'scroll'}} 
																className="select-box--items"
															>
															{ this.state.stateList.items &&
																this.state.stateList.items.map(_state => <div 
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
											<div className="col-xs-6">

											</div>
										</div>

										<div className="row">
											<div className="col-xs-6">
												<div className={'form-group' + (submitted && !user.landLord ? ' has-error' : '')}  ref={LLLNode => this.LLLNode = LLLNode}>
													<label htmlFor="landLord">Land Lord List:</label>
													<div className="select-box--box" style={{width: this.state.width || "100%"}}>
														<div className="select-box--container" onClick={this.handleClickLLL}>
															<div
																className="select-box--selected-item"
																onClick={this.dropDownLandLordList}
															> { this.state.selectedLandLord.Name}
															</div>
															<div 
																className="select-box--arrow"
																onClick={this.dropDownLandLordList}>
																	<span className={`${this.state.showLandLordList ? 'select-box--arrow-up' : 'select-box--arrow-down'}`}/>
															</div>
															<div 
																style={{display: this.state.showLandLordList ? 'block' : 'none'}} 
																className="select-box--items"
															>
															{ this.state.landLordList.items &&
																this.state.landLordList.items.map(landLord => <div 
																	key={landLord.LandLordID} 
																	onClick={() => this.selectLandLord(landLord)} 
																	className={this.state.selectedLandLord === landLord ? 'selected' : ''}
																>
																	{ landLord.Name }
																</div>)
															}
															</div>
														</div>
														<input
															type="hidden"
															value={user.landLord}
															name="landLord"
															onChange={this.handleChange.bind(this, "landLord")}
														/>
														{submitted && !user.landLord &&
															<div className="help-block"> Need to select a land lord</div>
														}
													</div>
												</div>
											</div>
											<div className="col-xs-6">
												<div className={'form-group' + (submitted && this.state.errors["rentalPaymentAmt"] ? ' has-error' : '')} style={{width: "200px"}}>
													<label htmlFor="rentalPaymentAmt">Rental Payment Amount</label>
													<CurrencyInput inputType="text" prefix="$" name="rentalPaymentAmt" value={user.rentalPaymentAmt} onChangeEvent={this.handleChange.bind(this, "rentalPaymentAmt")} style={{hidden: true}}/>
													<div className="help-block">{this.state.errors["rentalPaymentAmt"]}</div>
												</div>
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
		const { stateList, landLordList } = state;
    return {
        registering,
				stateList,
				landLordList,
    };
}

const connectedRegisterPage = connect(mapStateToProps)(RegisterPage);
export { connectedRegisterPage as RegisterPage };
