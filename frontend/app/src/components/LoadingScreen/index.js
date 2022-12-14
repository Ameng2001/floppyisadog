import React from 'react';
import PropTypes from 'prop-types';
import { CircularProgress } from '@rmwc/circular-progress';

require('./loading-screen.scss');
const imgUrl = require(
  '../../../../resources/images/staffjoy.png'
);

function LoadingScreen({ containerProps = {} }) {
  return (
    <div className="loading-container" {...containerProps}>
      <img role="presentation" alt="Staffjoy logo" src={imgUrl} />
      <CircularProgress size="small" />
    </div>
  );
}

LoadingScreen.propTypes = {
  containerProps: PropTypes.object,
};

export default LoadingScreen;
