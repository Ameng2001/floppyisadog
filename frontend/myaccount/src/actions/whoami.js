import _ from 'lodash';
import 'whatwg-fetch';
import * as actionTypes from '../constants/actionTypes';
import {
  emptyPromise,
  timestampExpired,
  checkStatus,
  parseJSON,
  routeToMicroservice,
  checkCode,
} from '../utility';

function shouldFetchWhoAmI(state) {
  const whoAmIState = state.whoami;
  const whoAmIData = whoAmIState.data;

  // it has never been fetched
  if (_.isEmpty(whoAmIData)) {
    return true;

  // it's currently being fetched
  } else if (whoAmIState.isFetching) {
    return false;

  // it's been in the UI for more than the allowed threshold
  } else if (!whoAmIState.lastUpdate ||
    (timestampExpired(whoAmIState.lastUpdate, 'WHOAMI'))
  ) {
    return true;
  }

  // otherwise, fetch if it's been invalidated
  return whoAmIState.didInvalidate;
}

function requestWhoAmI() {
  return {
    type: actionTypes.REQUEST_WHO_AM_I,
  };
}

function receiveWhoAmI(data) {
  return {
    type: actionTypes.RECEIVE_WHO_AM_I,
    lastUpdate: Date.now(),
    data,
  };
}

function fetchWhoAmI() {
  return (dispatch) => {
    dispatch(requestWhoAmI());

    return fetch(routeToMicroservice('whoami', '/whoami/'), {
      credentials: 'include', //携带cookie用作鉴权
    })
      .then(checkStatus)
      .then(parseJSON)
      //.then(checkCode) 非tars服务断点，不用checkCode
      .then(data =>
        dispatch(receiveWhoAmI(data))
      );
  };
}

function requestIntercomSettings() {
  return {
    type: actionTypes.REQUEST_INTERCOM_SETTINGS,
  };
}

function receiveIntercomSettings(data) {
  return {
    type: actionTypes.RECEIVE_INTERCOM_SETTINGS,
    data,
  };
}

export function getWhoAmI() {
  return (dispatch, getState) => {
    if (shouldFetchWhoAmI(getState())) {
      return dispatch(fetchWhoAmI());
    }
    return emptyPromise();
  };
}

export function fetchIntercomSettings() {
  return (dispatch) => {
    dispatch(requestIntercomSettings());

    return fetch(routeToMicroservice('whoami', '/intercom/'), {
      credentials: 'include', //携带cookie用作鉴权
    })
      .then(checkStatus)
      .then(parseJSON)
      //.then(checkCode)
      .then(data => dispatch(receiveIntercomSettings(data)));
  };
}
