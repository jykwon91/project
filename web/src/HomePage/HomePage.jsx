import React from 'react';
import { Link } from 'react-router-dom';
import { connect } from 'react-redux';
import { ListGroup, ListGroupItem, ListGroupItemHeading, ListGroupItemText } from 'reactstrap';

import { userActions } from '../_actions';

class HomePage extends React.Component {

		constructor(props) {
			super(props);

			this.state = {
				notifications: [{}]
			}
		}

/*
    componentDidMount() {
        this.props.dispatch(userActions.getAll());
    }
    handleDeleteUser(id) {
        return (e) => this.props.dispatch(userActions.delete(id));
    }
*/
    componentDidMount() {
        this.props.dispatch(userActions.getAllNotifications());
    }
/*
											<li className="list-group-item" key={index}>
												<label>{notification.From}</label>
												{'\n' + notification.Time}
												{notification.Message}
											</li>
*/
    render() {
        const { notifications } = this.props;
        return (
					<div className="col-md-10 col-md-offset-1">
						<div className="row">
							<div className="col-sm-8">
								<p>Notications:</p>
								{notifications.loading && <em> Loading notifications...</em>}
								{notifications.error && <span className="text-danger">ERROR: {notifications.error}</span>}
								{notifications.items &&
									<ListGroup style={{overflow: "scroll",maxHeight: "400px", width: "100%"}}>
										{notifications.items.map((notification, index) =>
											<ListGroupItem key={index}>
												<ListGroupItemHeading>{notification.Time}</ListGroupItemHeading>
												<ListGroupItemHeading>From: {notification.From}</ListGroupItemHeading>
												<ListGroupItemText>{notification.Message}</ListGroupItemText>
											</ListGroupItem>
										)}
									</ListGroup>
								}
							</div>
						</div>
						<p>
								<Link to="/login">Logout</Link>
						</p>
					</div>
        );
    }
}

function mapStateToProps(state) {
    const { notifications, authentication } = state;
    const { user } = authentication;
    return {
        user,
        notifications
    };
}

const connectedHomePage = connect(mapStateToProps)(HomePage);
export { connectedHomePage as HomePage };
