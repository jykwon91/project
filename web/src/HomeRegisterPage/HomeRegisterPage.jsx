import React from 'react';
import { Link } from 'react-router-dom';
import { connect } from 'react-redux';
import { ListGroup, ListGroupItem, ListGroupItemHeading, ListGroupItemText } from 'reactstrap';
import { Navbar } from 'react-bootstrap';

import { userActions } from '../_actions';

class HomeRegisterPage extends React.Component {

	constructor(props) {
		super(props);

		this.state = {
			address: {
				Street: '',
				City: '',
				Zipcode: '',
				State: ''
			},
			submitted: false,
			stateList: props.stateList,
			selectedState: 'Select State',
			showStateList: false,
			fields: {},
			errors: {}
		}

		this.handleChange = this.handleChange.bind(this);
		this.handleSubmit = this.handleSubmit.bind(this);
		this.handleClick = this.handleClick.bind(this);
		this.handleOutsideClick = this.handleOutsideClick.bind(this);
	}

	dropDownStates = () => {
		this.setState(prevState => ({
			showStateList: !prevState.showStateList
		}))
	}

	selectState = (_state) => this.setState(prevState => ({
		selectedState: _state,
		showStateList: false,
		address: {
			...prevState.address,
			State: _state,
		}
	}))
	
	static getDerivedStateFromProps(props, state) {
		if (props.stateList !== state.stateList) {
			return {
				stateList: props.stateList
			};
		}
		return null;
	}

	handleClick() {
		if (!this.state.showStateList) {
			document.addEventListener('click', this.handleOutsideClick, false);
		} else {
			document.removeEventListener('click', this.handleOutsideClick, false);
		}

		this.setState({ showStateList: !this.state.showStateList });
	}

	handleOutsideClick(e) {
		if (this.node.contains(e.target)) {
			return;
		}

		this.handleClick();
	}

	handleChange(field, event) {
		const { name, value } = event.target;
		const { address } = this.state;
		var fields = this.state.fields;
		fields[field] = event.target.value;

		this.setState({
			address: {
				...address,
				[name]: value
			},
			fields: fields
		});
	}

	handleSubmit(event) {
		event.preventDefault();

		this.setState({ submitted: true});
		const { address } = this.state;
		const { dispatch } = this.props;

		if (this.handleValidation() && address.Street && address.City && address.Zipcode && address.State) {
			dispatch(userActions.registerLandLordProperty(address));
		}

	}

	handleValidation() {
		const fields = this.state.fields;
		const errors = {};
		var formIsValid = true;

		//Zipcode
		if (!fields["Zipcode"]) {
			formIsValid = false;
			errors["Zipcode"] = "Phone number is required";
		}

		if (typeof fields["Zipcode"] !== "undefined") {
			if (!fields["Zipcode"].match(/^[0-9]+$/)) {
				formIsValid = false;
				errors["Zipcode"] = "Zipcode can only consist of numbers(0-9)";
			}
		}

		this.setState({errors: errors});
		return formIsValid;
	}

    componentDidMount() {
        this.props.dispatch(userActions.getStateList());
				//Need this to be an api to get all rental addresses...not sure why
        //this.props.dispatch(userActions.get());
    }

    render() {
        const { registering, stateList, user } = this.props;
		const { address, submitted } = this.state
        return (
            <div className="col-md-10 col-md-offset-1">
				<form name="form" onSubmit={this.handleSubmit}>
					<h3>Rental Home Address</h3>
					<div className="row">
						<div className="col-xs-6">
							<div className={'form-group' + (submitted && !address.Street ? ' has-error' : '')}>
								<label htmlFor="Street">Street</label>
								<input type="text" className="form-control" name="Street" value={address.Street} onChange={this.handleChange.bind(this, "Street")} />
								{submitted && !address.Street &&
										<div className="help-block">Street is required</div>
								}
							</div>
						</div>
						<div className="col-xs-6">
							<div className={'form-group' + (submitted && !address.City ? ' has-error' : '')}>
								<label htmlFor="City">City</label>
								<input type="text" className="form-control" name="City" value={address.City} onChange={this.handleChange.bind(this, "City")} />
								{submitted && !address.City &&
										<div className="help-block">City is required</div>
								}
							</div>
						</div>
						<div className="col-xs-6">
							<div className={'form-group' + (submitted && this.state.errors["Zipcode"] ? ' has-error' : '')} style={{width: "100px"}}>
								<label htmlFor="Zipcode">Zipcode</label>
								<input type="text" className="form-control" name="Zipcode" value={address.Zipcode} onChange={this.handleChange.bind(this, "Zipcode")} maxLength="5" />
								<div className="help-block">{this.state.errors["Zipcode"]}</div>
							</div>
						</div>
						<div className="col-xs-6">
							<label htmlFor="State">State</label>
							<div className={'form-group' + (submitted && !address.State ? ' has-error' : '')} ref={node => this.node = node}>
								<div className="select-box--box" style={{width: this.state.width || "150px"}}>
									<div className="select-box--container" onClick={this.handleClick}>
										<div
											className="select-box--selected-item"
											onClick={this.dropDownStates}
										> { this.state.selectedState }
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
												onClick={() => this.selectState(_state.name)}
												className={this.state.selectedState === _state.name ? 'selected' : ''}
											>
												{ _state.value }
											</div>)
										}
										</div>
									</div>
									<input
										type="hidden"
										value={address.State.name}
										name="State"
										onChange={this.handleChange.bind(this, "State")}
									/>
									{submitted && !address.State &&
										<div className="help-block"> State is required</div>
									}
								</div>
							</div>
						</div>
					</div>
                    <div className="form-group">
                        <button className="btn btn-primary">Register</button>
                        {registering &&
                            <img src="data:image/gif;base64,R0lGODlhEAAQAPIAAP///wAAAMLCwkJCQgAAAGJiYoKCgpKSkiH/C05FVFNDQVBFMi4wAwEAAAAh/hpDcmVhdGVkIHdpdGggYWpheGxvYWQuaW5mbwAh+QQJCgAAACwAAAAAEAAQAAADMwi63P4wyklrE2MIOggZnAdOmGYJRbExwroUmcG2LmDEwnHQLVsYOd2mBzkYDAdKa+dIAAAh+QQJCgAAACwAAAAAEAAQAAADNAi63P5OjCEgG4QMu7DmikRxQlFUYDEZIGBMRVsaqHwctXXf7WEYB4Ag1xjihkMZsiUkKhIAIfkECQoAAAAsAAAAABAAEAAAAzYIujIjK8pByJDMlFYvBoVjHA70GU7xSUJhmKtwHPAKzLO9HMaoKwJZ7Rf8AYPDDzKpZBqfvwQAIfkECQoAAAAsAAAAABAAEAAAAzMIumIlK8oyhpHsnFZfhYumCYUhDAQxRIdhHBGqRoKw0R8DYlJd8z0fMDgsGo/IpHI5TAAAIfkECQoAAAAsAAAAABAAEAAAAzIIunInK0rnZBTwGPNMgQwmdsNgXGJUlIWEuR5oWUIpz8pAEAMe6TwfwyYsGo/IpFKSAAAh+QQJCgAAACwAAAAAEAAQAAADMwi6IMKQORfjdOe82p4wGccc4CEuQradylesojEMBgsUc2G7sDX3lQGBMLAJibufbSlKAAAh+QQJCgAAACwAAAAAEAAQAAADMgi63P7wCRHZnFVdmgHu2nFwlWCI3WGc3TSWhUFGxTAUkGCbtgENBMJAEJsxgMLWzpEAACH5BAkKAAAALAAAAAAQABAAAAMyCLrc/jDKSatlQtScKdceCAjDII7HcQ4EMTCpyrCuUBjCYRgHVtqlAiB1YhiCnlsRkAAAOwAAAAAAAAAAAA==" />
                        }
						<p><Link to="/landlord">Cancel</Link></p>
                    </div>
				</form>
            </div>
        );
    }
}

function mapStateToProps(state) {
	const { registering } = state.registration;
    const { address, stateList, authentication } = state;
	const { user } = authentication;
    return {
		user,
		registering,
        address,
		stateList
    };
}

const connectedHomeRegisterPage = connect(mapStateToProps)(HomeRegisterPage);
export { connectedHomeRegisterPage as HomeRegisterPage };
