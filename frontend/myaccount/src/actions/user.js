import _ from 'lodash';
import 'whatwg-fetch';

// TODO when we have an API import request from 'request';
// import PhoneNumber from 'awesome-phonenumber';

import * as actionTypes from '../constants/actionTypes';
import { setFormResponse } from './forms';
import { getWhoAmI } from './whoami';
import {
  emptyPromise,
  timestampExpired,
  checkStatus,
  parseJSON,
  routeToMicroservice,
  checkCode,
} from '../utility';

// TODO delete this once we start fetching
const delay = ms => new Promise(resolve =>
  setTimeout(resolve, ms)
);

function updatingPassword() {
  return {
    type: actionTypes.UPDATING_PASSWORD,
  };
}

function updatedPassword() {
  return {
    type: actionTypes.UPDATED_PASSWORD,
  };
}

function updatePassword(userId, password) {
  return (dispatch, getState) => {
    dispatch(updatingPassword());

    //tars json调用参数，与tars接口定义一致
    const req = {
      req: {
        uuid: userId,
        password: password,
      }
    }

    return fetch(routeToMicroservice(
      'account',
      //json固定；account为accountserver别名；UpdatePassword同rpc方法名，统一用POST方法
      '/json/account/UpdatePassword'
    ), {
      credentials: 'include',
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        //携带X-Token用于网关鉴权
        //'X-Token': getState().whoami.data.token,
        //whoami相当于鉴权，透传uuid和身份信息authz
        'X-Verify-UID': getState.whoami.data.user_uuid,
        'X-Verify-Data': getState.whoami.data.authz,
      },
      body: JSON.stringify(req),
    })
      .then(checkStatus)
      .then(parseJSON)
      .then(checkCode)
      .then(() => {
        dispatch(setFormResponse('passwordUpdate', {
          type: 'success',
          message: 'Password updated!',
        }));
        return dispatch(updatedPassword());
      })
      .catch(() =>
        dispatch(setFormResponse('passwordUpdate', {
          type: 'danger',
          message: 'Passwords must be at least 6 characters long',
        }))
      );
  };
}

export function changePassword(newPassword) {
  return (dispatch, getState) =>
    dispatch(updatePassword(
      getState().whoami.data.user_uuid,
      newPassword
    ));
}

function updatingPhoto() {
  return {
    type: actionTypes.UPDATING_PHOTO,
  };
}

function updatedPhoto(photoUrl) {
  return {
    type: actionTypes.UPDATED_PHOTO,
    data: {
      photoUrl,
    },
  };
}

function updatePhoto(userId, photoReference) {
  return (dispatch) => {
    dispatch(updatingPhoto());

    return delay(500).then(() => {
      const response = {
        data: {
          photoUrl: photoReference,
        },
      };

      return dispatch(updatedPhoto(response.data.photoUrl));
    });
  };
}

export function changePhoto(event) {
  const photoLocalLocation = event.target.value;

  return (dispatch, getState) =>
    dispatch(updatePhoto(getState().user.userId, photoLocalLocation));
}

// state will be update before the patch is made
function updatingUser(data) {
  return {
    type: actionTypes.UPDATING_USER,
    ...data,
  };
}

function updatedUser(data) {
  return {
    type: actionTypes.UPDATED_USER,
    ...data,
  };
}

function updateUser(userId, data) {
  return (dispatch, getState) => {
    const userData = getState().user.data;
    const originalEmail = userData.email;
    let successMessage = 'Success!';
    //tars json调用参数，与tars接口定义一致
    const updateReq = {
      req: _.extend({}, userData, data)
    }

    if (data.email !== originalEmail) {
      successMessage += ' Check your email for a confirmation link.';
    }

    dispatch(updatingUser({ data }));

    return fetch(routeToMicroservice('account', '/json/account/Update'), {
      credentials: 'include',
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        //携带X-Token用于网关鉴权
        //'X-Token': getState().whoami.data.token,
        //whoami相当于鉴权，透传uuid和身份信息authz
        'X-Verify-UID': getState.whoami.data.user_uuid,
        'X-Verify-Data': getState.whoami.data.authz,
      },
      body: JSON.stringify(updateReq),
    })
      .then(checkStatus)
      .then(parseJSON)
      .then(checkCode)
      .then((response) => {
        dispatch(setFormResponse('accountUpdate', {
          type: 'success',
          message: successMessage,
        }));
        return dispatch(updatedUser({
          data: response.account,
          lastUpdate: Date.now(),
        }));
      })
      .catch(() =>
        dispatch(setFormResponse('accountUpdate', {
          type: 'danger',
          message: 'Unable to save changes',
        }))
      );
  };
}

function requestUser() {
  return {
    type: actionTypes.REQUEST_USER,
  };
}

function receiveUser(data) {
  return {
    type: actionTypes.RECEIVE_USER,
    ...data,
  };
}

function fetchUser(userId) {
  return (dispatch, getState) => {
    // dispatch action to start the fetch
    dispatch(requestUser());

    //tars json调用参数，与tars接口定义一致
    const getReq = {
      req: {
        uuid: userId,
      }
    }

    // eslint-disable-next-line max-len
    return fetch(routeToMicroservice('account', `/json/account/Get`), {
      credentials: 'include',
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        //携带X-Token用于网关鉴权
        //'X-Token': getState().whoami.data.token,
        //whoami相当于鉴权，透传uuid和身份信息authz
        'X-Verify-UID': getState.whoami.data.user_uuid,
        'X-Verify-Data': getState.whoami.data.authz,
      },
      body: JSON.stringify(getReq),
    })
      .then(checkStatus)
      .then(parseJSON)
      .then(checkCode)
      .then(data =>
        dispatch(receiveUser({
          //tars接口的返回格式，rsp为AccountInfo结构
          data: data.rsp,
          lastUpdate: Date.now(),
        }))
      );
  };
}

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

/*
  Exported funcitons:
  * initialize  // gets the userId and then calls getUser
  * getUser     // data for the user (needs a user id)
  * changeAccountData
  * modifyUserAttribute
*/

export function getUser(userId) {
  return (dispatch, getState) => {
    if (shouldFetchUser(getState())) {
      return dispatch(fetchUser(userId));
    }
    return emptyPromise();
  };
}

export function initialize() {
  //利用redux-thunk中间件，dispatch action creater，将dispatch和getState作为参数传入
  return (dispatch, getState) => {
    dispatch(getWhoAmI()).then(() => {
      const userId = getState().whoami.data.user_uuid;

      return dispatch(getUser(userId));
    });
  };
}

export function changeAccountData(email, name, phoneNumber) {
  // make API call to save the submitted changes
  return (dispatch, getState) =>
    dispatch(updateUser(getState().whoami.data.user_uuid, {
      email,
      name,
      phoneNumber,
    }));
}

export function modifyUserAttribute(event) {
  const target = event.target;
  const inputType = target.getAttribute('type');
  const attribute = target.getAttribute('data-model-attribute');
  const payload = {};

  if (inputType === 'checkbox') {
    payload[attribute] = target.checked;
  } else {
    payload[attribute] = target.value;
  }

  return (dispatch, getState) => {
    dispatch(updateUser(getState().whoami.data.user_uuid, payload));
  };
}
