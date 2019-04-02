import React, { Component } from 'react';
import { BrowserRouter as Router, Route} from 'react-router-dom';
import Register from './components/Register'
import Verify from './components/Verify'

class App extends Component {
  render() {
    return (
      <Router>
        <div className="App">
          <div className = "container">
            <Route exact path ="/" component = {Register} />
            <Route path="/verify/:id" component={Verify}/>
          </div>
        </div>
      </Router>
    );
  }
}

export default App;
