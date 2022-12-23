import _ from 'lodash';
import 'whatwg-fetch';
import { normalize, schema } from 'normalizr';
import { invalidateAssociations } from '../associations';
import * as actionTypes from '../../constants/actionTypes';
import { routeToMicroservice } from '../../constants/paths';
import {
  emptyPromise,
  timestampExpired,
  checkStatus,
  parseJSON,
  checkCode,
} from '../../utility';

/*
  Exported functions:
  * getTeamEmployees
  * createTeamEmployee
*/

// schemas!
const teamEmployeesSchema = new schema.Entity('employees', {}, { idAttribute: 'user_uuid' });
const arrayOfTeamEmployees = new schema.Array(teamEmployeesSchema);

// team employees
function requestTeamEmployees(teamUuid) {
  return {
    type: actionTypes.REQUEST_TEAM_EMPLOYEES,
    teamUuid,
  };
}

function receiveTeamEmployees(teamUuid, data) {
  return {
    type: actionTypes.RECEIVE_TEAM_EMPLOYEES,
    teamUuid,
    ...data,
  };
}

function creatingTeamEmployee(teamUuid) {
  return {
    type: actionTypes.CREATING_TEAM_EMPLOYEE,
    teamUuid,
  };
}

function createdTeamEmployee(teamUuid, userUuid, data) {
  return {
    type: actionTypes.CREATED_TEAM_EMPLOYEE,
    teamUuid,
    userUuid,
    ...data,
  };
}

function fetchTeamEmployees(companyUuid, teamUuid) {
  return (dispatch, getState) => {
    // dispatch action to start the fetch
    dispatch(requestTeamEmployees(teamUuid));
    const teamEmployeePath =
      '/json/company/ListWorkers';

    return fetch(
      routeToMicroservice('company', teamEmployeePath),
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
            team_uuid: teamUuid
          }
        }),
      })
      .then(checkStatus)
      .then(parseJSON)
      .then(checkCode)
      .then((data) => {
        const normalized = normalize(
          _.get(data.rsp, 'workers', []),
          arrayOfTeamEmployees
        );

        return dispatch(receiveTeamEmployees(teamUuid, {
          data: normalized.entities.employees,
          order: normalized.result,
          lastUpdate: Date.now(),
        }));
      });
  };
}

function shouldFetchTeamEmployees(state, teamUuid) {
  const employeesData = state.teams.employees;
  const teamEmployees = _.get(employeesData, teamUuid, {});

  // no team employees have ever been fetched
  if (_.isEmpty(employeesData)) {
    return true;

  // the needed teamUuid is empty
  } else if (_.isEmpty(teamEmployees)) {
    return true;

  // teamEmployees is at least partially populated with a trusted object at this point
  // the order of these is related to how the 1st fetch might play out

  // this data set is currently being fetched
  } else if (teamEmployees.isFetching) {
    return false;

  // this data set is not complete
  } else if (!teamEmployees.completeSet) {
    return true;

  // this data set is stale
  } else if (!teamEmployees.lastUpdate ||
    (timestampExpired(teamEmployees.lastUpdate, 'TEAM_EMPLOYEES'))
  ) {
    return true;
  }

  // check if invalidated
  return teamEmployees.didInvalidate;
}

export function getTeamEmployees(companyUuid, teamUuid) {
  return (dispatch, getState) => {
    if (shouldFetchTeamEmployees(getState(), teamUuid)) {
      return dispatch(fetchTeamEmployees(companyUuid, teamUuid));
    }
    return emptyPromise();
  };
}

export function createTeamEmployee(companyUuid, teamUuid, userUuid) {
  return (dispatch, getState) => {
    dispatch(creatingTeamEmployee(teamUuid));
    const workerPath = '/json/company/CreateWorker';

    return fetch(
      routeToMicroservice('company', workerPath), {
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
            team_uuid: teamUuid,
            user_uuid: userUuid
          }
        }),
      })
      .then(checkStatus)
      .then(parseJSON)
      .then(checkCode)
      .then((data) => {
        dispatch(invalidateAssociations());
        return dispatch(createdTeamEmployee(teamUuid, data.rsp.user_uuid, data.rsp));
      });
  };
}
