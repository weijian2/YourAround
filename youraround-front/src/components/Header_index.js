import React from 'react';
import logo from '../assets/images/logo.svg';

export class Header_index extends React.Component {
    render() {
        return (
            <header className="App-header-index">
                <img src={logo} className="App-logo-index" alt="logo" />
                <h1 className="App-title-index">YourAround</h1>
            </header>
        );
    }
}