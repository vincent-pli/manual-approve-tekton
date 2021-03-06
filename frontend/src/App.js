import React from 'react';
import { browserHistory } from 'react-router';
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";
import Dashboard from './components/Dashboard';
import Datatable from './components/Tables';

class App extends React.Component {
    constructor(props) {
      super(props);
    }

    render(){
      return (
        <Router history={browserHistory}>
          <div>       
            <Switch>
              <Route path="/#/requests" component={Datatable} >
              </Route>
              <Route exact path="/" >
                <Dashboard/>
              </Route>
              
            </Switch>
          </div>
        </Router>
      )
    }
  }

function Home() {
  return (
    <div>
      <h2>Dashboard</h2>
    </div>
  );
}

export default App;
