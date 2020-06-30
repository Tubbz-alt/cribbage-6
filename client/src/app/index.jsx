import React from 'react';
import { ConnectedRouter } from 'connected-react-router';
import { Switch, Route } from 'react-router-dom';
import { history } from '../store/reducers';
import Home from '../components/Home';
import Landing from '../components/landing/Landing';
import PrivateRoute from '../components/layout/PrivateRoute';
import Login from '../components/landing/Login';

export const App = () => {
  return (
    <ConnectedRouter history={history}>
      <div className='relative bg-gray-200 h-screen'>
        <Navbar />
        <Alert />
        <Switch>
          <Route exact path='/' component={Landing} />
          <Route exact path='/login' component={Login} />
          <PrivateRoute exact path='/home' component={Home} />
        </Switch>
      </div>
    </ConnectedRouter>
  );
};
