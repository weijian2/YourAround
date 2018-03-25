import React, { Component } from 'react';
import '../styles/App.css';
import { Header } from './Header'
import { Header_index } from './Header_index';
import { Main } from './Main';
import { TOKEN_KEY } from "../constants"

class App extends Component {
    state = {
        isLoggedIn : !!localStorage.getItem(TOKEN_KEY),
    }

    handleLogin = (token) => {
        localStorage.setItem(TOKEN_KEY, token);
        this.setState({
            isLoggedIn : true,
        });
    }

    handleLogout = () => {
        localStorage.removeItem(TOKEN_KEY);
        this.setState({
            isLoggedIn : false,
        });
    }

    render() {
        return (
            <div className="App">
                {this.state.isLoggedIn ?
                    <Header_index isLoggedIn={this.state.isLoggedIn} handleLogout={this.handleLogout}/>
                    : <Header />}
                <Main isLoggedIn={this.state.isLoggedIn} handleLogin={this.handleLogin}/>
            </div>
        );
    }
}

export default App;