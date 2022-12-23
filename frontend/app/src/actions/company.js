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


function requestCompany() {
  return {
    type: actionTypes.REQUEST_COMPANY,
  };
}

function receiveCompany(data) {
  return {
    type: actionTypes.RECEIVE_COMPANY,
    ...data,
  };
}

function fetchCompany(companyUuid) {
  return (dispatch, getState) => {
    // dispatch an action when the request is initiated
    dispatch(requestCompany());
    const companyPath = '/json/company/GetCompany';

    return fetch(routeToMicroservice('company', companyPath), {
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
          uuid: companyUuid,
        }
      }),
    })
      .then(checkStatus)
      .then(parseJSON)
      .then(checkCode)
      .then(data =>
        dispatch(receiveCompany({
          data: data.rsp,
        }))
      );
  };
}

function shouldFetchCompany(companyUuid, state) {
  const companyState = state.company;
  const companyData = companyState.data;

  // it has never been fetched
  if (_.isEmpty(companyData)) {
    return true;

  // it's a different company
  } else if (companyData.uuid !== companyUuid) {
    return true;

  // it's currently being fetched
  } else if (companyState.isFetching) {
    return false;

  // it's been in the UI for more than the allowed threshold
  } else if (!companyState.lastUpdate ||
    (timestampExpired(companyState.lastUpdate, 'COMPANY'))
  ) {
    return true;
  }

  // otherwise, fetch if it's been invalidated, I suppose
  return companyState.didInvalidate;
}

export function getCompany(companyUuid) {
  return (dispatch, getState) => {
    if (shouldFetchCompany(companyUuid, getState())) {
      return dispatch(fetchCompany(companyUuid));
    }
    return emptyPromise();
  };
}

// TODO - need 'changeCompany' action - invalidate all company related data to force a refetch
