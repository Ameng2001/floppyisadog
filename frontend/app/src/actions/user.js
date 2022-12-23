import _ from 'lodash';
import 'whatwg-fetch';
import * as actionTypes from '../constants/actionTypes';
import { routeToMicroservice } from '../constants/paths';
import {
  emptyPromise,
  timestampExpired,
  checkStatus,
  checkCode,
  parseJSON,
} from '../utility';

import { getWhoAmI } from './whoami';

function shouldFetchUser(state) {
  const userState = state.user;
  const userData = userState.data;

  // it has never been fetched
  if (_.isEmpty(userData)) {
    return true;

  // it's currently being fetched
  } else if (userState.isFetching) {
    return false;

  // it's been in the UI for more than the allowed threshold
  } else if (!userState.lastUpdate ||
    (timestampExpired(userState.lastUpdate, 'USER'))
  ) {
    return true;
  }

  // otherwise, fetch if it's been invalidated
  return userState.didInvalidate;
}

function requestUser() {
  return {
    type: actionTypes.REQUEST_USER,
  };
}

function receiveUser(data) {
  return {
    type: actionTypes.RECEIVE_USER,
    lastUpdate: Date.now(),
    ...data,
  };
}

function fetchUser(userUuid) {
  return (dispatch, getState) => {
    dispatch(requestUser());

    return fetch(routeToMicroservice('account', '/json/account/Get'), {
      //credentials: 'include',
      method: 'POST',
      mode:'cors',
      headers: {
        'Content-Type': 'application/json',
        //携带X-Token用于网关鉴权
        //'X-Token': getState().whoami.data.token,
        //whoami相当于鉴权，透传uuid和身份信息authz(设置到request header中，在网关侧会将head透传到tars调用的context中)
        'X-Verify-UID': getState().whoami.data.user_uuid,
        'X-Verify-Data': getState().whoami.data.authz,
      },
      body: JSON.stringify({
        req: {
          uuid: userUuid,
        }
      }),
    })
      .then(checkStatus)
      .then(parseJSON)
      .then(checkCode)
      .then(data =>
        dispatch(receiveUser({
          data: data.rsp,
          lastUpdate: Date.now(),
        }))
      );
  };
}

export function getUser() {
  return (dispatch, getState) => {
    // whoami data is required
    dispatch(getWhoAmI()).then(() => {
      const state = getState();
      const userUuid = state.whoami.data.user_uuid;

      if (shouldFetchUser(state)) {
        return dispatch(fetchUser(userUuid));
      }

      return emptyPromise();
    });
  };
}
