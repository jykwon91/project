import React from 'react';
import { Link } from 'react-router-dom';
import { connect } from 'react-redux';

import { userActions } from '../_actions';

class LandLordPage extends React.Component {

		constructor(props) {
			super(props);

			this.state = {
				submitted: false,
				notification: {
					Message: '',
					From: '',
					Time: '', //in epoch time
				},
				fields: {},
				errors: {}
			}

			this.handleChange = this.handleChange.bind(this);
			this.handleSubmit = this.handleSubmit.bind(this);
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

			this.setState({ submitted: true});
			//const { notification } = this.state;
			const { dispatch } = this.props;

			if (notification.message) {
				dispatch(userActions.sendNotification(notification));
			}
		}

		handleDeleteUser(id) {
			return (e) => this.props.dispatch(userActions.delete(id));
		}

    componentDidMount() {
        this.props.dispatch(userActions.getAll());
    }

    render() {
        const { user, users, notification, submitted } = this.props;
        return (
            <div className="col-md-10 col-md-offset-1">
								<h2 style={{width: "100%"}}>Send notification:</h2>
								<div className={'form-group' + (submitted && this.state.errors["notification"] ? ' has-error' : '')}>

								</div>
                <h3>All registered users:</h3>
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

								<p><Link to="/home/register">Register Rental Home</Link></p>
                <p><Link to="/login">Logout</Link></p>
            </div>
        );
    }
}

function mapStateToProps(state) {
    const { users, authentication } = state;
    const { user } = authentication;
    return {
        user,
        users
    };
}

const connectedLandLordPage = connect(mapStateToProps)(LandLordPage);
export { connectedLandLordPage as LandLordPage };
