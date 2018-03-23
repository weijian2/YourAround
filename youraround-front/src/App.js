import React, { Component } from 'react';
import './App.css';
import { Header } from './Header'
import { Header_index } from './Header_index';
import { Main } from './Main';

class App extends Component {
  render() {
    return (
      <div className="App">
        <Header_index/>
        <Main/>
      </div>
    );
  }
}

export default App;