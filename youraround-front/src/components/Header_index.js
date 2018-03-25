import React from 'react';
import logo from '../assets/images/logo.svg';
import { Icon } from 'antd'
export class Header_index extends React.Component {
    render() {
        return (
            <header className="App-header-index">
                <img src={logo} className="App-logo-index" alt="logo" />
                <h1 className="App-title-index">YourAround</h1>
                {this.props.isLoggedIn ?
                    <a className="logout"
                       onClick={this.props.handleLogout}>
                        <Icon type="logout" />{' '}logout
                    </a>
                    : null}
            </header>
        );
    }
}