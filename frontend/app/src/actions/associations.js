import _ from 'lodash';
import 'whatwg-fetch';
import * as actionTypes from '../constants/actionTypes';
import { routeToMicroservice } from '../constants/paths';
import {
  emptyPromise,
  timestampExpired,
  checkStatus,
  parseJSON,
  checkCode,
} from '../utility';

/*
  Exported functions:
  * getAssociations
  * invalidateAssociations
*/

// action creators

function requestAssociations() {
  return {
    type: actionTypes.REQUEST_ASSOCIATIONS,
  };
}

function receiveAssociations(data) {
  return {
    type: actionTypes.RECEIVE_ASSOCIATIONS,
    ...data,
  };
}

function fetchAssociations(companyUuid) {
  return (dispatch,getState) => {
    // dispatch action to start the fetch
    dispatch(requestAssociations());
    const associationPath = '/json/company/GetAssociations';

    return fetch(
      routeToMicroservice('company', associationPath),
      { 
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
            company_uuid: companyUuid,
            limit: 100,
            offset: 0
          }
        }),
      })
      .then(checkStatus)
      .then(parseJSON)
      .then(checkCode)
      .then((data) => {
        const result = {};
        _.forEach(_.get(data.rsp, 'accounts', []), (account) => {
          result[account.account.user_uuid] = {
            teams: _.get(account, 'teams', []),
          };
        });

        return dispatch(receiveAssociations({
          data: result,
          lastUpdate: Date.now(),
        }));
      });
  };
}

function shouldFetchAssociations(state) {
  const associationsState = state.associations;
  const associationsData = associationsState.data;

  // it has never been fetched
  if (_.isEmpty(associationsData)) {
    return true;

  // it's currently being fetched
  } else if (associationsState.isFetching) {
    return false;

  // it's been in the UI for more than the allowed threshold
  } else if (!associationsState.lastUpdate ||
    (timestampExpired(associationsState.lastUpdate, 'ASSOCIATIONS'))
  ) {
    return true;
  }

  // otherwise, fetch if it's been invalidated
  return associationsState.didInvalidate;
}

export function getAssociations(companyUuid) {
  return (dispatch, getState) => {
    if (shouldFetchAssociations(getState())) {
      return dispatch(fetchAssociations(companyUuid));
    }
    return emptyPromise();
  };
}

export function invalidateAssociations() {
  return {
    type: actionTypes.INVALIDATE_ASSOCIATIONS,
  };
}
