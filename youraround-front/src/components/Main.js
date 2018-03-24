import React from 'react';
import { Register } from './Register'
import { Login } from './Login'

export class Main extends React.Component {
    render() {
        return (
            <div className="main">
                <Login/>
            </div>
        );
    }
}