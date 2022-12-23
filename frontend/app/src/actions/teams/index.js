import _ from 'lodash';
import 'whatwg-fetch';
import { normalize, schema } from 'normalizr';
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
  * getTeams
  * getTeam
*/

// schemas!
const teamSchema = new schema.Entity('teams', {}, { idAttribute: 'uuid' });
const arrayOfTeams = new schema.Array(teamSchema);

// teams

function requestTeams() {
  return {
    type: actionTypes.REQUEST_TEAMS,
  };
}

function receiveTeams(data) {
  return {
    type: actionTypes.RECEIVE_TEAMS,
    ...data,
  };
}

function fetchTeams(companyUuid) {
  return (dispatch,getState) => {
    // dispatch action to start the fetch
    dispatch(requestTeams());

    return fetch(
      routeToMicroservice('company', '/json/company/ListTeams'),
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
          }
        }),
      })
      .then(checkStatus)
      .then(parseJSON)
      .then(checkCode)
      .then((data) => {
        const normalized = normalize(_.get(data.rsp, 'teams', []), arrayOfTeams);

        return dispatch(receiveTeams({
          data: normalized.entities.teams,
          order: normalized.result,
          lastUpdate: Date.now(),
        }));
      });
  };
}

function shouldFetchTeams(state) {
  const teamsState = state.teams;
  const teamsData = teamsState.data;

  // it has never been fetched
  if (_.isEmpty(teamsData)) {
    return true;

  // it's currently being fetched
  } else if (teamsState.isFetching) {
    return false;

  // it's been in the UI for more than the allowed threshold
  } else if (!teamsState.lastUpdate ||
    (timestampExpired(teamsState.lastUpdate, 'TEAMS'))
  ) {
    return true;

  // make sure we have a complete collection too
  } else if (!teamsState.completeSet) {
    return true;
  }

  // otherwise, fetch if it's been invalidated
  return teamsState.didInvalidate;
}

// determines if should fetch teams or extract from current state
export function getTeams(companyUuid) {
  return (dispatch, getState) => {
    if (shouldFetchTeams(getState())) {
      return dispatch(fetchTeams(companyUuid));
    }
    return emptyPromise();
  };
}

// team

function requestTeam() {
  return {
    type: actionTypes.REQUEST_TEAM,
  };
}

function receiveTeam(data) {
  return {
    type: actionTypes.RECEIVE_TEAM,
    ...data,
  };
}

function fetchTeam(companyUuid, teamUuid) {
  return (dispatch, getState) => {
    // dispatch action to start the fetch
    dispatch(requestTeam());
    const teamPath = '/json/company/GetTeam';

    return fetch(
      routeToMicroservice('company', teamPath),
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
            uuid: teamUuid
          }
        }),
      })
      .then(checkStatus)
      .then(parseJSON)
      .then(checkCode)
      .then((data) => {
        const normalized = normalize(data.rsp, teamSchema);

        return dispatch(receiveTeam({
          data: normalized.entities.teams,
        }));
      });
  };
}

function shouldFetchTeam(state, teamUuid) {
  const teamsState = state.teams;
  const teamsData = teamsState.data;

  // no team has ever been fetched
  if (_.isEmpty(teamsData)) {
    return true;

  // the needed teamUuid is not available
  } else if (!_.has(teamsData, teamUuid)) {
    return true;

  // it's been in the UI for more than the allowed threshold
  } else if (!teamsState.lastUpdate ||
    (timestampExpired(teamsState.lastUpdate, 'TEAMS'))
  ) {
    return true;
  }

  // otherwise, fetch if it's been invalidated
  return teamsState.didInvalidate;
}

export function getTeam(companyUuid, teamUuid) {
  return (dispatch, getState) => {
    if (shouldFetchTeam(getState(), teamUuid)) {
      return dispatch(fetchTeam(companyUuid, teamUuid));
    }
    return emptyPromise();
  };
}

export {
  getTeamJobs,
  updateTeamJob,
  updateTeamJobField,
  setTeamJob,
  createTeamJob,
} from './jobs';
export { getTeamEmployees, createTeamEmployee } from './employees';
export {
  getTeamShifts,
  updateTeamShift,
  bulkUpdateTeamShifts,
  deleteTeamShift,
  createTeamShift,
} from './shifts';
