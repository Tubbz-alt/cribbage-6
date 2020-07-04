import React from 'react';
import { useSelector } from 'react-redux';
import { selectCurrentUser } from '../../../auth/selectors';
import ActiveGames from './active_game_list';

const Home = () => {
  const currentUser = useSelector(selectCurrentUser);

  return <div>
    Welcome, {currentUser.name}!
    <ActiveGames/>
  </div>;
};

export default Home;
