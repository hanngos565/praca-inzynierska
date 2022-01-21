import React from "react";
import './App.css';
import { BrowserRouter as Router, Switch, Route } from 'react-router-dom';
import {MainView} from "./views/MainView/MainView";
import {ResultsView} from "./views/ResultsView/ResultsView";

function App() {
  return (
        <Router>
            <Switch>
                <Route exact path="/" component={MainView} />
                <Route exact path="/simulation-results/demo" component={ResultsView} />
            </Switch>
        </Router>
  );
}

export default App;
