import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';

const title = 'Manual Approve';

window._env_ = {
  APPROVE_URL: window.location.href,
}

ReactDOM.render(
  <App title={title} />,
  document.getElementById('app')
);

module.hot.accept();
